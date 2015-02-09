package main

import (
	"fmt"
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
		check()
		ui.Task(fmt.Sprintf("Sleeping for %s...", sleepDuration))
		time.Sleep(sleepDuration)
	}
}

func check() {
	ui.Task("Checking repos.")
	repos, registry, err := repo.All()
	if err != nil {
		ui.Err(err.Error())
		return
	}

	for _, repo := range repos {
		checkRepo(repo, registry)
	}
}

func checkRepo(r *repo.Repo, registry string) {
	if strings.Index(r.URL, "http") != 0 {
		ui.Warn(fmt.Sprintf("Skipping %s, only http[s] is supported", r.URL))
		return
	}

	head, err := r.Checkout()
	if err != nil {
		ui.Err(fmt.Sprintf("Error checking %s: %v\n", r.URL, err))
		return
	}

	if r.SHA == head {
		return
	}

	ui.Info("%s has been updated from %s to %s. Starting a new build.", r.URL, r.SHA, head)
	if err := image.Make(r); err != nil {
		ui.Err(fmt.Sprintf("Error making %s: %v", r.URL, err))
		return
	}

	img, err := image.Build(r, head, registry)
	if err != nil {
		ui.Err(fmt.Sprintf("Error building %s: %v", r.URL, err))
		return
	}

	ui.Info("%s version %s has been built", r.URL, head)
	if err := image.Publish(img); err != nil {
		ui.Err(fmt.Sprintf("Error publishing %s: %v", r.URL, err))
		return
	}
	ui.Info("%s has been published to %s", r.URL, img)
	r.SHA = head
	r.Image = img
	if err := r.Save(); err != nil {
		ui.Err(fmt.Sprintf("Error updating %s: %v", r.URL, err))
	}
}
