package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"core-gitlab.corp.zulily.com/core/stevedore/image"
	"core-gitlab.corp.zulily.com/core/stevedore/repo"
	"core-gitlab.corp.zulily.com/core/stevedore/slack"
	"core-gitlab.corp.zulily.com/core/stevedore/ui"
)

var (
	sleepDuration = 1 * time.Minute
	cfg           config
	notifications *slack.Slack
)

type config struct {
	sync.Mutex
	PublishCommand []string `json:"publishCommand"`
	RegistryURL    string   `json:"registryUrl"`
	Notifications  []string `json:"notifications"`
	Slack          struct {
		Channel  string `json:"channel"`
		Username string `json:"username"`
		Webhook  string `json:"webhook"`
	}
}

func loadConfig() error {
	cfg.Lock()
	defer cfg.Unlock()

	jsonFile := filepath.Clean("./config.json")
	file, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, &cfg)
}

// ImagePublishCommand returns the command line strings to use to publish an image.
func (c *config) ImagePublishCommand(image string) []string {
	if c.PublishCommand == nil {
		return []string{"docker", "push", image}
	}

	return append(c.PublishCommand, image)
}

func main() {
	var err error
	err = loadConfig()
	if err != nil {
		ui.Err(err.Error())
		os.Exit(1)
	}

	ui.Info("loaded config")

	if contains(cfg.Notifications, "slack") && cfg.Slack.Webhook != "" {
		notifications, err = slack.New(
			slack.WithWebhook(cfg.Slack.Webhook),
			slack.WithChannelAndUsername(cfg.Slack.Channel, cfg.Slack.Username))
	}
	if err != nil {
		ui.Err(err.Error())
		return
	}

	for {
		updated := check()
		ui.Task("%s repo images updated, built and published. Sleeping for %s...", strconv.Itoa(updated), sleepDuration.String())
		time.Sleep(sleepDuration)
	}
}

func check() (updated int) {
	ui.Task("Checking repos.")
	repos, err := repo.All()
	if err != nil {
		ui.Err(err.Error())
		return 0
	}

	for _, repo := range repos {
		if checkRepo(repo, cfg.RegistryURL) {
			updated++
		}
	}

	return updated
}

func checkRepo(r *repo.Repo, registry string) (updated bool) {
	if strings.Index(r.URL, "https://") != 0 {
		ui.Warn(fmt.Sprintf("Skipping %s, only https is supported", r.URL))
		return false
	}

	head, err := r.Checkout()
	if err != nil {
		ui.Err(fmt.Sprintf("Error checking %s: %v\n", r.URL, err))
		return false
	}

	if r.SHA == head {
		return false
	}

	ui.Info("%s has been updated from %s to %s. Starting a new build.", r.URL, r.SHA, head)
	if err := image.Make(r); err != nil {
		ui.Err(fmt.Sprintf("Error making %s: %v", r.URL, err))
		return false
	}

	imgs, err := image.Build(r, head, registry)
	if err != nil {
		ui.Err(fmt.Sprintf("Error building %s: %v", r.URL, err))
		return false
	}

	for _, img := range imgs {
		ui.Info("%s version %s has been built", r.URL, head)

		if err := image.Publish(cfg.ImagePublishCommand(img)); err != nil {
			ui.Err(fmt.Sprintf("Error publishing %s: %v", r.URL, err))
			return false
		}

		msg := fmt.Sprintf("A new image for %s has been published to %s", r.URL, img)
		ui.Info(msg)

		if notifications != nil {
			if err = notifications.Notify(msg); err != nil {
				ui.Err(err.Error())
			}
		}
	}

	r.SHA = head
	r.Images = imgs
	if err := r.Save(); err != nil {
		ui.Err(fmt.Sprintf("Error updating %s: %v", r.URL, err))
	}

	return true
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
