package email

import (
	"bytes"
	"context"
	"io"

	"github.com/oneee-playground/r2d2-api-server/internal/global/email"
	"gopkg.in/gomail.v2"
)

type GomailSender struct {
	fromAddr string
	dialer   *gomail.Dialer
}

var _ email.Sender = (*GomailSender)(nil)

type GomailOptions struct {
	Host     string
	Port     int
	Username string
	Password string

	FromAddr string
}

func NewGomailSender(opt GomailOptions) *GomailSender {
	s := GomailSender{}

	s.dialer = gomail.NewDialer(opt.Host, opt.Port, opt.Username, opt.Password)
	s.fromAddr = opt.FromAddr

	return &s
}

func (s *GomailSender) Send(ctx context.Context, address string, subject string, content *bytes.Buffer) error {
	msg := gomail.NewMessage()

	msg.SetHeader("From", s.fromAddr)
	msg.SetHeader("To", address)
	msg.SetHeader("Subject", subject)
	msg.AddAlternativeWriter("test/plain", func(w io.Writer) error {
		_, err := content.WriteTo(w)
		return err
	})

	return s.dialer.DialAndSend(msg)
}
