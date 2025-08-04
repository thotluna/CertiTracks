package mailer_test

import (
	"certitrack/internal/mailer"
	"certitrack/pkg/util"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordResetEmailIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	mailhog := NewMailHogHelper()
	cfg := TestConfig(t)

	mailer, err := mailer.NewSMTPMailer(cfg)
	require.NoError(t, err, "Should not fail when creating mailer")

	testID := fmt.Sprintf("test-%d", time.Now().UnixNano())
	testSubject := fmt.Sprintf("[CertiTracks] Test %s", testID)

	resetToken := "test-reset-token-123"
	appURL := util.GetEnv("APP_URL", "http://localhost:3000")
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", appURL, resetToken)

	testData := map[string]interface{}{
		"Data": map[string]string{
			"Name":     "Test User",
			"ResetURL": resetURL,
		},
	}

	err = mailhog.ClearMessages()
	require.NoError(t, err, "Should not fail when clearing previous messages")

	t.Run("Email sending integration", func(t *testing.T) {
		err := mailer.SendEmail(
			"test@example.com",
			testSubject,
			"password_reset.html",
			testData,
		)

		assert.NoError(t, err, "Should not fail when sending email")

		time.Sleep(500 * time.Millisecond)

		found, err := mailhog.FindMessage("test@example.com", testSubject)
		require.NoError(t, err, "Should not fail when searching for email in MailHog")
		assert.True(t, found, "Email should be in MailHog")

		// Verificar que el correo se envió al destinatario correcto
		messages, err := mailhog.GetMessages()
		require.NoError(t, err, "Should not fail when getting messages from MailHog")
		require.Greater(t, messages.Total, 0, "Should have at least one message")

		// Verificar que el correo se envió al destinatario correcto
		msg := messages.Items[0]
		assert.Contains(t, msg.Content.Headers.To[0], "test@example.com", "Should send to the correct email")
	})
}
