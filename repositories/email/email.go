package email

import "context"

type Email interface {
	SendVerificationEmail(ctx context.Context, recipient string, otp string) error
}
