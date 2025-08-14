package mock

import (
	"context"
	"log/slog"

	"github.com/rhodeon/go-backend-template/internal/log"
	"github.com/rhodeon/go-backend-template/repositories/email"
	"github.com/rhodeon/go-backend-template/utils/contextutils"
)

type Email struct{}

func New() email.Email {
	return &Email{}
}

func (p *Email) SendVerificationEmail(ctx context.Context, recipient string, otp string) error {
	contextutils.GetLogger(ctx).Info(
		"Sending verification email.",
		slog.String(log.AttrRecipient, recipient),
		slog.String(log.AttrOtp, otp),
	)
	return nil
}
