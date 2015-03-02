package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type message struct {
	Text     string `json:"text"`
	Channel  string `json:"channel"`
	Username string `json:"username"`
}

// Slack represents a Notifier for the Slack team messaging system (https://slack.com/)
type Slack struct {
	cfg *Config
}

// Config contains all the necessary config for sending slack notifications to a channel
type Config struct {
	Username string
	Channel  string
	Webhook  string
}

// WithWebHook returns a configuration function that sets the slack webhook address.
func WithWebhook(webhook string) func(*Config) error {
	return func(c *Config) error {
		c.Webhook = webhook
		return nil
	}
}

// WithChannelAndUsername returns a configuration function that sets the slack channel and username.
func WithChannelAndUsername(channel, username string) func(*Config) error {
	return func(c *Config) error {
		c.Username = username
		c.Channel = channel
		return nil
	}
}

// NewSlack creates a Slack instance. Configuration functions may be
// supplied to override the default settings of the Config.
// The default settings are as follows:
//    Webhook:       "",
//    Username:   "stevedore",
//    Password:   "#stevedore",
func New(fns ...func(*Config) error) (*Slack, error) {
	var err error
	cfg := &Config{
		Username: "stevedore",
		Channel:  "#stevedore",
		Webhook:  "",
	}

	for _, fn := range fns {
		if err == nil {
			err = fn(cfg)
		}
	}

	return &Slack{
		cfg: cfg,
	}, nil
}

// Notify implements Notifier for Slack
func (s *Slack) Notify(msg string) error {
	body, err := json.Marshal(&message{
		Channel:  s.cfg.Channel,
		Username: s.cfg.Username,
		Text:     msg,
	})

	if err != nil {
		return err
	}

	client := http.Client{}
	resp, err := client.PostForm(s.cfg.Webhook, url.Values{"payload": {string(body)}})

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
