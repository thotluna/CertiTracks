package mailer

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"net/smtp"
)

type mockSMTP struct {
	from       string
	to         []string
	msg        []byte
	shouldFail bool
}

type mockSMTPWriter struct {
	buffer *bytes.Buffer
}

func (w *mockSMTPWriter) Write(p []byte) (int, error) {
	return w.buffer.Write(p)
}

func (w *mockSMTPWriter) Close() error {
	return nil
}

func (m *mockSMTP) Mail(from string) error {
	m.from = from
	if m.shouldFail {
		return errors.New("mock SMTP error: authentication required")
	}
	return nil
}

func (m *mockSMTP) Rcpt(to string) error {
	m.to = append(m.to, to)
	if m.shouldFail {
		return errors.New("mock SMTP error: recipient rejected")
	}
	return nil
}

func (m *mockSMTP) Data() (io.WriteCloser, error) {
	if m.shouldFail {
		return nil, errors.New("mock SMTP error: data command failed")
	}
	return &mockSMTPWriter{buffer: &bytes.Buffer{}}, nil
}

func (m *mockSMTP) Close() error {
	return nil
}

func (m *mockSMTP) StartTLS(config *tls.Config) error {
	if m.shouldFail {
		return errors.New("mock SMTP error: could not start TLS")
	}
	return nil
}

func (m *mockSMTP) Auth(auth smtp.Auth) error {
	if m.shouldFail {
		return errors.New("mock SMTP error: authentication failed")
	}
	return nil
}
