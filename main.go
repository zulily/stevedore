package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"core-gitlab.corp.zulily.com/core/stevedore/image"
	"core-gitlab.corp.zulily.com/core/stevedore/repo"
	"core-gitlab.corp.zulily.com/core/stevedore/ui"
)

var (
	sleepDuration = 1 * time.Minute
)

func main() {
	for {
		updated := check()
		ui.Task("%s repo images updated, built and published. Sleeping for %s...", strconv.Itoa(updated), sleepDuration.String())
		time.Sleep(sleepDuration)
	}
}

func check() (updated int) {
	ui.Task("Checking repos.")
	repos, registry, err := repo.All()
	if err != nil {
		ui.Err(err.Error())
		return 0
	}

	for _, repo := range repos {
		if checkRepo(repo, registry) {
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

	img, err := image.Build(r, head, registry)
	if err != nil {
		ui.Err(fmt.Sprintf("Error building %s: %v", r.URL, err))
		return false
	}

	ui.Info("%s version %s has been built", r.URL, head)
	if err := image.Publish(img); err != nil {
		ui.Err(fmt.Sprintf("Error publishing %s: %v", r.URL, err))
		return false
	}
	ui.Info("%s has been published to %s", r.URL, img)
	r.SHA = head
	r.Image = img
	if err := r.Save(); err != nil {
		ui.Err(fmt.Sprintf("Error updating %s: %v", r.URL, err))
	}

	return true
}
