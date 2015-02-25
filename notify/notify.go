package notify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
)

// A Notificer represents a component that sends stevedore messages to some external system
type Notifier interface {
	Notify(string) error
}

// Init instantiates zero or more Notifier instances, based on the configuration JSON
func Init() ([]Notifier, error) {
	var notifiers []Notifier

	jsonFile := filepath.Clean("./config.json")
	file, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return notifiers, err
	}

	// unmarshall into an anonymous struct to get the list of notifications
	config := struct {
		Notifications []string `json:"notifications"`
	}{}
	err = json.Unmarshal(file, &config)

	if contains(config.Notifications, "slack") {
		notifier, err := newSlack(file)
		if err != nil {
			return notifiers, fmt.Errorf("Error configuration slack notifications: %s", err.Error())
		}
		notifiers = append(notifiers, notifier)
	}
	// TODO: other notifiers

	return notifiers, err
}

type slackOptions struct {
	Channel  string `json:"channel"`
	Username string `json:"username"`
	Webhook  string `json:"webhook"`
}

type slackMessage struct {
	Text     string `json:"text"`
	Channel  string `json:"channel"`
	Username string `json:"username"`
}

// Slack represents a Notifier for the Slack team messaging system (https://slack.com/)
type Slack struct {
	opts *slackOptions
}

// newSlack creates a Slack Notifier, configured according to the supplied configuration
// JSON bytes
func newSlack(configFile []byte) (Notifier, error) {

	config := struct {
		Slack slackOptions `json:"slack"`
	}{}

	err := json.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	return &Slack{
		opts: &config.Slack,
	}, nil
}

// Notify implements Notifier for Slack
func (s *Slack) Notify(msg string) error {
	body, err := json.Marshal(&slackMessage{
		Channel:  s.opts.Channel,
		Username: s.opts.Username,
		Text:     msg,
	})

	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.PostForm(s.opts.Webhook, url.Values{"payload": {string(body)}})

	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		var errMsg string
		if byts, err := ioutil.ReadAll(resp.Body); err == nil {
			errMsg = string(byts)
		}
		return fmt.Errorf("Got a %d response from the slack API: %s", resp.StatusCode, errMsg)
	}
	return nil
}

// contains returns a boolean indicating whether or not `e` is contained in `s`
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
