package mailer_test

import (
	"certitrack/pkg/util"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type mailhogMessage struct {
	Content struct {
		Headers struct {
			To      []string `json:"To"`
			Subject []string `json:"Subject"`
			From    []string `json:"From"`
		} `json:"Headers"`
		Body string `json:"Body"`
	} `json:"Content"`
	Raw *struct {
		Data string `json:"Data"`
	} `json:"Raw"`
}

type mailhogResponse struct {
	Total int              `json:"total"`
	Items []mailhogMessage `json:"items"`
}

type MailHogHelper struct {
	baseURL string
}

func NewMailHogHelper() *MailHogHelper {
	return &MailHogHelper{
		baseURL: getMailHogAPIURL(),
	}
}

func (h *MailHogHelper) ClearMessages() error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v1/messages", h.baseURL), nil)
	if err != nil {
		return fmt.Errorf("error creating request to clear MailHog: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error clearing MailHog messages: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code when clearing MailHog: %d", resp.StatusCode)
	}

	return nil
}

func (h *MailHogHelper) GetMessages() (*mailhogResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/v2/messages", h.baseURL))
	if err != nil {
		return nil, fmt.Errorf("error getting messages from MailHog: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading MailHog response: %v", err)
	}

	var messages mailhogResponse
	if err := json.Unmarshal(body, &messages); err != nil {
		return nil, fmt.Errorf("error decoding MailHog response: %v", err)
	}

	return &messages, nil
}

func (h *MailHogHelper) FindMessage(to, subject string) (bool, error) {
	messages, err := h.GetMessages()
	if err != nil {
		return false, err
	}

	for _, msg := range messages.Items {
		if len(msg.Content.Headers.To) > 0 &&
			len(msg.Content.Headers.Subject) > 0 &&
			msg.Content.Headers.To[0] == to &&
			msg.Content.Headers.Subject[0] == subject {
			return true, nil
		}
	}

	return false, nil
}

func getMailHogAPIURL() string {
	return fmt.Sprintf("http://%s:8025/api", util.GetEnv("SMTP_HOST", "localhost"))
}
