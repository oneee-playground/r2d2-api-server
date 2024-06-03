package email

import (
	"bytes"
	"context"
	"io"

	"github.com/oneee-playground/r2d2-api-server/internal/global/email"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type GomailSender struct {
	dialer *gomail.Dialer
	logger *zap.Logger

	fromAddr string
}

var _ email.Sender = (*GomailSender)(nil)

type GomailOptions struct {
	Host     string
	Port     int
	Username string
	Password string

	FromAddr string
}

func NewGomailSender(logger *zap.Logger, opt GomailOptions) *GomailSender {
	s := GomailSender{logger: logger}

	s.dialer = gomail.NewDialer(opt.Host, opt.Port, opt.Username, opt.Password)
	s.fromAddr = opt.FromAddr

	return &s
}

func (s *GomailSender) Send(ctx context.Context, address string, subject string, content *bytes.Buffer) error {
	msg := gomail.NewMessage()

	msg.SetHeader("From", s.fromAddr)
	msg.SetHeader("To", address)
	msg.SetHeader("Subject", subject)
	msg.AddAlternativeWriter("text/plain", func(w io.Writer) error {
		_, err := content.WriteTo(w)
		return err
	})

	err := s.dialer.DialAndSend(msg)
	if err != nil {
		return errors.Wrap(err, "sending email")
	}

	s.logger.Info("sent email", zap.String("subject", subject), zap.String("email", address))

	return nil
}
