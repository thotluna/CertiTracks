package mailer_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"certitrack/internal/mailer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSMTPIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	mailhog := NewMailHogHelper()
	cfg := TestConfig(t)

	mailer, err := mailer.NewSMTPMailer(cfg)
	require.NoError(t, err, "Should not fail when creating mailer")

	testID := fmt.Sprintf("test-%d", time.Now().UnixNano())
	testSubject := fmt.Sprintf("[CertiTracks] Test %s", testID)

	testData := map[string]string{
		"Timestamp": time.Now().Format(time.RFC3339),
		"TestData":  "Test data for integration " + testID,
	}

	err = mailhog.ClearMessages()
	require.NoError(t, err, "Should not fail when clearing previous messages")

	t.Run("Email sending integration", func(t *testing.T) {
		err := mailer.SendEmail(
			"test@example.com",
			testSubject,
			"test_email.html",
			testData,
		)

		assert.NoError(t, err, "Should not fail when sending email")

		time.Sleep(500 * time.Millisecond)

		found, err := mailhog.FindMessage("test@example.com", testSubject)
		require.NoError(t, err, "Should not fail when searching for email in MailHog")
		assert.True(t, found, "Email should be in MailHog")
	})
}

func TestMain(m *testing.M) {
	result := m.Run()

	if result == 0 && testing.CoverMode() != "" {
		if tc := testing.Coverage(); tc < 0.8 {
			fmt.Printf("Tests passed but coverage failed: got %.2f, want 0.80\n", tc)
			os.Exit(1)
		}
	}

	os.Exit(result)
}
