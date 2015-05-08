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

	"github.com/zulily/stevedore/image"
	"github.com/zulily/stevedore/repo"
	"github.com/zulily/stevedore/slack"
	"github.com/zulily/stevedore/ui"
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

	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return err
	}

	if len(cfg.PublishCommand) == 0 {
		cfg.PublishCommand = []string{"docker", "push"}
	}
	return nil
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

// notify reports a msg to the specified ui output func, and to any configured
// notifications.  Additionally, if any errors ocurr during notificiation, that
// error is sent to the UI's Err func.
func notify(msg, output string, uiFunc func(string, ...string)) {
	uiFunc(msg)
	if notifications != nil {
		msg = fmt.Sprintf("%s\n%s", msg, output)
		if err := notifications.Notify(msg); err != nil {
			ui.Err(err.Error())
		}
	}
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

	// Update and persist the new SHA now, so that if a build/publish fails, it
	// won't repeate endlessly
	r.SHA = head
	if err := r.Save(); err != nil {
		ui.Err(fmt.Sprintf("Error updating %s: %v", r.URL, err))
	}

	if err := r.PrepareMake(); err != nil {
		ui.Err(fmt.Sprintf("Error preparing %s: %v", r.URL, err))
		return false
	}

	if output, err := image.Make(r); err != nil {
		msg := fmt.Sprintf("Error making %s: %v", r.URL, err)
		notify(msg, output, ui.Err)
		return false
	}

	imgs, output, err := image.Build(r, head, registry)
	if err != nil {
		msg := fmt.Sprintf("Error building %s: %v", r.URL, err)
		notify(msg, output, ui.Err)
		return false
	}

	cmd := cfg.PublishCommand
	for _, img := range imgs {
		ui.Info("%s version %s has been built", r.URL, head)

		if output, err := image.Publish(img, cmd); err != nil {
			msg := fmt.Sprintf("Error publishing %s: %v", r.URL, err)
			notify(msg, output, ui.Err)
			return false
		}

		msg := fmt.Sprintf("A new image for %s has been published to %s", r.URL, img)
		notify(msg, "", ui.Info)
	}

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
