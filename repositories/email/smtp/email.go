package smtp

import (
	"context"
	"html/template"
	"strings"
	"time"

	"github.com/rhodeon/go-backend-template/repositories/email"

	"github.com/go-errors/errors"
	"github.com/wneessen/go-mail"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

type Email struct {
	client  *mail.Client
	config  *Config
	tmpl    *template.Template
	tracer  trace.Tracer
	meter   metric.Meter
	counter metric.Int64Counter
}

func New(ctx context.Context, cfg *Config) (email.Email, error) {
	client, err := mail.NewClient(
		cfg.Host,
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.User),
		mail.WithPassword(cfg.Password),
		mail.WithPort(cfg.Port),
	)
	if err != nil {
		return nil, errors.Errorf("creating client: %w", err)
	}

	if err := pingServer(ctx, client); err != nil {
		return nil, errors.Errorf("pinging server: %w", err)
	}

	tmpl, err := email.NewTemplate()
	if err != nil {
		return nil, errors.Errorf("creating email template: %w", err)
	}

	meter := otel.GetMeterProvider().Meter(cfg.OtelServiceName)
	counter, err := meter.Int64Counter(
		"email.total_sent",
		metric.WithDescription("Total number of emails sent."),
	)
	if err != nil {
		return nil, errors.Errorf("creating sent emails counter: %w", err)
	}

	return &Email{
		client:  client,
		config:  cfg,
		tmpl:    tmpl,
		tracer:  otel.GetTracerProvider().Tracer(cfg.OtelServiceName),
		meter:   meter,
		counter: counter,
	}, nil
}

// pingServer simulates pinging by checking the SMTP connection to confirm the credentials are valid.
// A temporary SMTP client is established for this.
func pingServer(ctx context.Context, client *mail.Client) error {
	smtpClient, err := client.DialToSMTPClientWithContext(ctx)
	if err != nil {
		return errors.Errorf("establishing smtp connecting: %w", err)
	}
	if err := smtpClient.Close(); err != nil {
		return errors.Errorf("closing smtp ping client: %w", err)
	}

	return nil
}

// send is a wrapper around the sending operation that traces the outgoing email.
func (e *Email) send(ctx context.Context, message *mail.Msg, spanName string) error {
	sender, _ := message.GetSender(true)
	recipients, _ := message.GetRecipients()
	recipientsStr := strings.Join(recipients, ", ")

	newCtx, span := e.tracer.Start(ctx, spanName,
		trace.WithAttributes(
			semconv.NetworkProtocolName("smtp"),
			semconv.ServerAddress(e.config.Host),
			semconv.ServerPort(e.config.Port),
			attribute.String("email.from", sender),
			attribute.String("email.to", recipientsStr),
		),
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	e.counter.Add(
		newCtx,
		1,
		metric.WithAttributes(attribute.String("email.operation_id", spanName)),
	)

	if err := e.client.DialAndSendWithContext(newCtx, message); err != nil {
		err = errors.Errorf("dialing and sending: %w", err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	return nil
}

// SendVerificationEmail sends an email with the OTP to an unverified user.
func (e *Email) SendVerificationEmail(ctx context.Context, recipient string, otp string) error {
	subject := "Account Verification"
	emailTemplate, err := email.GenerateOtpTemplate(e.tmpl, email.OtpTemplateData{
		Subject:           subject,
		Otp:               otp,
		DurationInSeconds: int(e.config.OtpDuration / time.Second),
	})
	if err != nil {
		return errors.Errorf("generating otp template: %w", err)
	}

	message := mail.NewMsg()
	if err := message.From(e.config.Sender); err != nil {
		return errors.Errorf("setting sender email: %w", err)
	}

	if err := message.To(recipient); err != nil {
		return errors.Errorf("setting recipient email: %w", err)
	}

	message.Subject(subject)
	message.SetBodyString(mail.TypeTextHTML, emailTemplate.Html)
	message.AddAlternativeString(mail.TypeTextPlain, emailTemplate.PlainText)

	if err := e.send(ctx, message, "send-verification-email"); err != nil {
		return errors.Errorf("sending email: %w", err)
	}

	return nil
}
