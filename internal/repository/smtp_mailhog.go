package repository

import (
	"context"
	"fmt"

	"github.com/ffauzann/loan-service/internal/model"
	"github.com/ffauzann/loan-service/internal/util"
)

// SendMail sends an email using the SMTP client.
func (r *notificationRepository) SendMail(ctx context.Context, req *model.EmailRequest) error {
	if r.enabled == false {
		util.LogContext(ctx).Info("Email sending is disabled, skipping SendMail")
		return nil
	}

	// 1. MAIL FROM
	if err := r.smtp.Mail(req.From); err != nil {
		util.LogContext(ctx).Error(err.Error())
		return err
	}

	// 2. RCPT TO
	if err := r.smtp.Rcpt(req.To); err != nil {
		util.LogContext(ctx).Error(err.Error())
		return err
	}

	// 3. DATA
	writer, err := r.smtp.Data()
	if err != nil {
		util.LogContext(ctx).Error(err.Error())
		return err
	}

	// 4. Write email headers + body
	body := fmt.Sprintf(
		"To: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n%s\r\n",
		req.To, req.From, req.Subject, req.Body,
	)

	if _, err := writer.Write([]byte(body)); err != nil {
		util.LogContext(ctx).Error(err.Error())
		return err
	}

	if err := writer.Close(); err != nil {
		util.LogContext(ctx).Error(err.Error())
		return err
	}

	// 5. Close the SMTP session
	if err := r.smtp.Quit(); err != nil {
		util.LogContext(ctx).Error(err.Error())
		return err
	}

	util.LogContext(ctx).Info(fmt.Sprintf("Email sent successfully to: %s, subject: %s", req.To, req.Subject))
	return nil
}
