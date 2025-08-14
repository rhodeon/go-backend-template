package email

import (
	"embed"
	"html/template"
	"strings"

	"github.com/go-errors/errors"
)

//go:embed "templates"
var templatesFS embed.FS

const (
	otpTemplateFilePlainText = "otp.gotmpl"
	otpTemplateFileHtml      = "otp.gohtml"
)

// EmailBody represents the structure of an email body, containing both plain text and HTML content.
type EmailBody struct {
	PlainText string
	Html      string
}

// NewTemplate parses and initialises all relevant email templates.
// This avoids having to reparse the templates for every email sending attempt.
func NewTemplate() (*template.Template, error) {
	tmpl, err := template.New("root").ParseFS(
		templatesFS,
		"templates/"+otpTemplateFilePlainText,
		"templates/"+otpTemplateFileHtml,
	)
	if err != nil {
		return nil, errors.Errorf("parsing template files: %w", err)
	}
	return tmpl, nil
}

type OtpTemplateData struct {
	Subject           string
	Otp               string
	DurationInSeconds int
}

func GenerateOtpTemplate(tmpl *template.Template, data OtpTemplateData) (EmailBody, error) {
	var plainTextBody strings.Builder
	var htmlBody strings.Builder

	if err := tmpl.ExecuteTemplate(&plainTextBody, otpTemplateFilePlainText, data); err != nil {
		return EmailBody{}, errors.Errorf("executing otp plain text template file: %w", err)
	}

	if err := tmpl.ExecuteTemplate(&htmlBody, otpTemplateFileHtml, data); err != nil {
		return EmailBody{}, errors.Errorf("executing otp html template file: %w", err)
	}

	return EmailBody{
		PlainText: plainTextBody.String(),
		Html:      htmlBody.String(),
	}, nil
}
