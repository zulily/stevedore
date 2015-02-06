package main

import (
	"fmt"
	"strings"
	"time"

	"core-gitlab.corp.zulily.com/core/stevedore/image"
	"core-gitlab.corp.zulily.com/core/stevedore/repo"
	"github.com/mgutz/ansi"
)

var (
	taskColor = ansi.ColorCode("blue+h:black")
	errColor  = ansi.ColorCode("red+hb:black")
	warnColor = ansi.ColorCode("yellow:black")
	infoColor = ansi.ColorCode("white:black")
	reset     = ansi.ColorCode("reset")
)

func printTask(msg string) {
	fmt.Println(taskColor, msg, reset)
}

func printErr(msg string) {
	fmt.Println(errColor, msg, reset)
}

func printWarn(msg string) {
	fmt.Println(warnColor, msg, reset)
}

func printInfo(msg string) {
	fmt.Println(infoColor, msg, reset)
}

var (
	sleepDuration = 1 * time.Minute
)

func main() {
	shutdown := make(chan bool)
	startBuilder(shutdown)
	<-shutdown
}

func startBuilder(shutdown chan bool) {
	go func() {
		for {
			check()
			printTask(fmt.Sprintf("Sleeping for %s...", sleepDuration))
			time.Sleep(sleepDuration)
		}
		shutdown <- true
	}()
}

func check() {
	printTask("Checking repos.")
	repos, err := repo.All()
	if err != nil {
		printErr(err.Error())
		return
	}

	for _, repo := range repos {
		checkRepo(repo)
	}
}

func checkRepo(r *repo.Repo) {
	if strings.Index(r.URL, "http") != 0 {
		printWarn(fmt.Sprintf("Skipping %s, only http[s] is supported", r.URL))
		return
	}

	head, err := r.Checkout()
	if err != nil {
		printErr(fmt.Sprintf("Error checking %s: %v\n", r.URL, err))
		return
	}

	if r.LastCommit != head {
		printInfo(fmt.Sprintf("%s has been updated from %s to %s. Starting a new build.", r.URL, r.LastCommit, head))
		if img, err := image.Build(r, head); err == nil {
			printInfo(fmt.Sprintf("%s version %s has been built", r.URL, head))
			if err := image.Publish(img); err == nil {
				printInfo(fmt.Sprintf("%s has been published to %s", r.URL, img))
				r.LastCommit = head
				if err := r.Save(); err != nil {
					printErr(fmt.Sprintf("Error updating %s: %v", r.URL, err))
				}
			} else {
				printErr(fmt.Sprintf("Error publishing %s: %v", r.URL, err))
			}
		} else {
			printErr(fmt.Sprintf("Error building %s: %v", r.URL, err))
		}
	}
}
