package mailer

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSMTPMailer_SendEmail(t *testing.T) {
	testCases := []struct {
		name        string
		shouldFail  bool
		expectError bool
		template    string
		data        interface{}
	}{
		{
			name:        "Successful send",
			shouldFail:  false,
			expectError: false,
			template:    "test_email.html",
			data: map[string]string{
				"Timestamp": time.Now().Format(time.RFC3339),
				"TestData":  "Test data",
			},
		},
		{
			name:        "SMTP error",
			shouldFail:  true,
			expectError: true,
			template:    "test_email.html",
			data:        map[string]string{},
		},
	}

	originalDial := smtpDial
	defer func() { smtpDial = originalDial }()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockClient := &mockSMTP{
				shouldFail: tc.shouldFail,
			}

			smtpDial = func(addr string) (smtpClient, error) {
				return mockClient, nil
			}

			mailer := &smtpMailer{
				config: Config{
					From:     "test@example.com",
					Host:     "smtp.example.com",
					Port:     587,
					Username: "user",
					Password: "pass",
				},
				templates: template.Must(template.ParseGlob("templates/*.html")),
			}

			err := mailer.SendEmail(
				"recipient@example.com",
				"Test Subject",
				tc.template,
				tc.data,
			)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "test@example.com", mockClient.from)
				assert.Contains(t, mockClient.to, "recipient@example.com")
			}
		})
	}
}

func TestTemplateLoading(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	if filepath.Base(cwd) != "mailer" {
		t.Fatalf("Tests must be run from the mailer directory")
	}

	tpl, err := template.ParseGlob("templates/test_email.html")
	require.NoError(t, err)
	require.NotNil(t, tpl)

	var buf bytes.Buffer
	err = tpl.ExecuteTemplate(&buf, "test_email.html", map[string]string{
		"Timestamp": "2023-01-01T00:00:00Z",
		"TestData":  "Test",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, buf.String())
}
