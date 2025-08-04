package mailer

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"io"
	"net/smtp"
	"os"
	"path/filepath"
	"runtime"
)

var smtpDial = func(addr string) (smtpClient, error) {
	return smtp.Dial(addr)
}

type smtpClient interface {
	Mail(string) error
	Rcpt(string) error
	Data() (io.WriteCloser, error)
	Close() error
	StartTLS(*tls.Config) error
	Auth(smtp.Auth) error
}

type smtpMailer struct {
	config    Config
	templates *template.Template
	tlsConfig *tls.Config
}

func NewSMTPMailer(cfg Config) (Mailer, error) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	templatePath := filepath.Join(dir, "templates/*.html")

	tmpl, err := template.ParseGlob(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load email templates from %s: %w", templatePath, err)
	}

	return &smtpMailer{
		config:    cfg,
		templates: tmpl,
		tlsConfig: &tls.Config{
			InsecureSkipVerify: os.Getenv("APP_ENV") == "development",
			ServerName:         cfg.Host,
		},
	}, nil
}

func (m *smtpMailer) SendEmail(to, subject, templateName string, data interface{}) error {
	var body bytes.Buffer

	if m.templates.Lookup(templateName) == nil {
		return fmt.Errorf("template %s not found", templateName)
	}

	templateData := struct {
		Data    interface{}
		Subject string
	}{
		Data:    data,
		Subject: subject,
	}

	if err := m.templates.ExecuteTemplate(&body, templateName, templateData); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"+
		"%s\r\n", to, m.config.From, subject, body.String()))

	c, err := smtpDial(fmt.Sprintf("%s:%d", m.config.Host, m.config.Port))
	if err != nil {
		return fmt.Errorf("error connecting to SMTP server: %w", err)
	}
	defer c.Close()

	if m.config.Port == 587 || m.config.Port == 465 {
		if err = c.StartTLS(m.tlsConfig); err != nil {
			return fmt.Errorf("error starting TLS: %w", err)
		}
	}

	if m.config.Username != "" && m.config.Password != "" {
		auth := smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host)
		if err = c.Auth(auth); err != nil {
			return fmt.Errorf("error authenticating: %w", err)
		}
	}

	if err = c.Mail(m.config.From); err != nil {
		return fmt.Errorf("error setting sender: %w", err)
	}
	if err = c.Rcpt(to); err != nil {
		return fmt.Errorf("error setting recipient: %w", err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("error preparing to send email: %w", err)
	}

	if _, err = w.Write(msg); err != nil {
		return fmt.Errorf("error writing email: %w", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("error closing email writer: %w", err)
	}

	return nil
}
