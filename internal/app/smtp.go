package app

import (
	"fmt"
	"net"
	"net/smtp"
)

type SMTP struct {
	Enabled bool
	MailHog MailHog
}

type MailHog struct {
	Host   string
	Port   uint32
	Client *smtp.Client
}

func (s *SMTP) prepare() error {
	return s.MailHog.connect()
}

func (m *MailHog) connect() error {
	addr := net.JoinHostPort(m.Host, fmt.Sprintf("%d", m.Port))

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	if m.Client, err = smtp.NewClient(conn, m.Host); err != nil {
		return err
	}

	return m.Client.Hello("loan-service")
}
