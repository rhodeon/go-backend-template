package smtp

import (
	"context"
	"html/template"
	"time"

	"github.com/rhodeon/go-backend-template/repositories/email"

	"github.com/go-errors/errors"
	"github.com/wneessen/go-mail"
)

type Email struct {
	client *mail.Client
	config *Config
	tmpl   *template.Template
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

	return &Email{
		client: client,
		config: cfg,
		tmpl:   tmpl,
	}, nil
}

// pingServer mimics pinging by checking the SMTP connection to confirm the credentials are valid.
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

	if err := e.client.DialAndSendWithContext(ctx, message); err != nil {
		return errors.Errorf("dialing and sending email: %w", err)
	}

	return nil
}
