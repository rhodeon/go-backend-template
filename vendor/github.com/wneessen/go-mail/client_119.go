// SPDX-FileCopyrightText: The go-mail Authors
//
// SPDX-License-Identifier: MIT

//go:build !go1.20
// +build !go1.20

package mail

import (
	"errors"

	"github.com/wneessen/go-mail/smtp"
)

// SendWithSMTPClient attempts to send one or more Msg using a provided smtp.Client with an
// established connection to the SMTP server. If the smtp.Client has no active connection to
// the server, SendWithSMTPClient will fail with an error. For each of the provided Msg, it
// will associate a SendError with the Msg in case of a transmission or delivery error.
//
// This method first checks for an active connection to the SMTP server. If the connection is
// not valid, it returns a SendError. It then iterates over the provided messages, attempting
// to send each one. If an error occurs during sending, the method records the error and
// associates it with the corresponding Msg. If multiple errors are encountered, it aggregates
// them into a single SendError to be returned.
//
// Parameters:
//   - client: A pointer to the smtp.Client that holds the connection to the SMTP server
//   - messages: A variadic list of pointers to Msg objects to be sent.
//
// Returns:
//   - An error that represents the sending result, which may include multiple SendErrors if
//     any occurred; otherwise, returns nil.
func (c *Client) SendWithSMTPClient(client *smtp.Client, messages ...*Msg) error {
	escSupport := false
	if client != nil {
		escSupport, _ = client.Extension("ENHANCEDSTATUSCODES")
	}
	if err := c.checkConn(client); err != nil {
		return &SendError{
			Reason: ErrConnCheck, errlist: []error{err}, isTemp: isTempError(err),
			errcode: errorCode(err), enhancedStatusCode: enhancedStatusCode(err, escSupport),
		}
	}
	var errs []*SendError
	for id, message := range messages {
		if message == nil {
			continue
		}
		if sendErr := c.sendSingleMsg(client, message); sendErr != nil {
			messages[id].sendError = sendErr

			var msgSendErr *SendError
			if errors.As(sendErr, &msgSendErr) {
				errs = append(errs, msgSendErr)
			}
		}
	}

	if len(errs) > 0 {
		if len(errs) > 1 {
			returnErr := &SendError{Reason: ErrAmbiguous}
			for i := range errs {
				returnErr.errlist = append(returnErr.errlist, errs[i].errlist...)
				returnErr.rcpt = append(returnErr.rcpt, errs[i].rcpt...)
			}

			// We assume that the error codes and flags from the last error we received should be the
			// indicator for the returned isTemp flag as well
			returnErr.isTemp = errs[len(errs)-1].isTemp
			returnErr.errcode = errs[len(errs)-1].errcode
			returnErr.enhancedStatusCode = errs[len(errs)-1].enhancedStatusCode

			return returnErr
		}
		return errs[0]
	}
	return nil
}
