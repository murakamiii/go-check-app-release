package slack

import (
	"bytes"
	"fmt"
	"net/http"
)

// Slack ...
type Slack struct {
	Client *http.Client
}

// PostMessage to Slack
func (slack *Slack) PostMessage(path string, text string) error {
	res, err := slack.Client.Post(
		fmt.Sprintf("https://hooks.slack.com/services/%s", path),
		"application/json",
		bytes.NewBuffer([]byte(`{"text":"`+text+`"}`)),
	)

	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Slack POST Error: HTTP Status %d", res.StatusCode)
	}
	return nil
}
