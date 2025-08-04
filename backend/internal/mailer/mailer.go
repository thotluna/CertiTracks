// Package mailer provides email sending functionality for the application.
// It supports different email providers and templates.
package mailer

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

type Mailer interface {
	SendEmail(to, subject, templateName string, data interface{}) error
}
