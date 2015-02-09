package main

import (
	"fmt"
	"strings"
	"time"

	"core-gitlab.corp.zulily.com/core/stevedore/Godeps/_workspace/src/github.com/mgutz/ansi"
	"core-gitlab.corp.zulily.com/core/stevedore/image"
	"core-gitlab.corp.zulily.com/core/stevedore/repo"
)

var (
	taskColor       = ansi.ColorCode("blue+h")
	brightTaskColor = ansi.ColorCode("white+h")
	errColor        = ansi.ColorCode("red")
	brightErrColor  = ansi.ColorCode("red+h")
	warnColor       = ansi.ColorCode("yellow")
	brightWarnColor = ansi.ColorCode("yellow+h")
	infoColor       = ansi.ColorCode("white")
	brightInfoColor = ansi.ColorCode("cyan+h")
	reset           = ansi.ColorCode("reset")
)

func printTask(msg string, args ...string) {
	if len(args) == 0 {
		fmt.Println(taskColor, msg, reset)
		return
	}
	colored := colorArgs(args, brightTaskColor, taskColor)
	fmt.Printf(msg+"\n", colored...)
}

func printErr(msg string, args ...string) {
	fmt.Println(errColor, msg, reset)
	if len(args) == 0 {
		fmt.Println(errColor, msg, reset)
		return
	}
	colored := colorArgs(args, brightErrColor, errColor)
	fmt.Printf(msg+"\n", colored...)
}

func printWarn(msg string, args ...string) {
	fmt.Println(warnColor, msg, reset)
	if len(args) == 0 {
		fmt.Println(warnColor, msg, reset)
		return
	}
	colored := colorArgs(args, brightWarnColor, warnColor)
	fmt.Printf(msg+"\n", colored...)
}

func printInfo(msg string, args ...string) {
	if len(args) == 0 {
		fmt.Println(infoColor, msg, reset)
	} else {
		colored := colorArgs(args, brightInfoColor, infoColor)
		fmt.Printf(msg+"\n", colored...)
	}
}

func colorArgs(args []string, color, reset string) []interface{} {
	var colored []interface{}
	for _, v := range args {
		colored = append(colored, color+v+reset)
	}
	return colored
}

var (
	sleepDuration = 1 * time.Minute
)

func main() {
	for {
		check()
		printTask(fmt.Sprintf("Sleeping for %s...", sleepDuration))
		time.Sleep(sleepDuration)
	}
}

func check() {
	printTask("Checking repos.")
	repos, registry, err := repo.All()
	if err != nil {
		printErr(err.Error())
		return
	}

	for _, repo := range repos {
		checkRepo(repo, registry)
	}
}

func checkRepo(r *repo.Repo, registry string) {
	if strings.Index(r.URL, "http") != 0 {
		printWarn(fmt.Sprintf("Skipping %s, only http[s] is supported", r.URL))
		return
	}

	head, err := r.Checkout()
	if err != nil {
		printErr(fmt.Sprintf("Error checking %s: %v\n", r.URL, err))
		return
	}

	if r.SHA == head {
		return
	}

	printInfo("%s has been updated from %s to %s. Starting a new build.", r.URL, r.SHA, head)
	if err := image.Make(r); err != nil {
		printErr(fmt.Sprintf("Error making %s: %v", r.URL, err))
		return
	}

	img, err := image.Build(r, head, registry)
	if err != nil {
		printErr(fmt.Sprintf("Error building %s: %v", r.URL, err))
		return
	}

	printInfo("%s version %s has been built", r.URL, head)
	if err := image.Publish(img); err != nil {
		printErr(fmt.Sprintf("Error publishing %s: %v", r.URL, err))
		return
	}
	printInfo("%s has been published to %s", r.URL, img)
	r.SHA = head
	r.Image = img
	if err := r.Save(); err != nil {
		printErr(fmt.Sprintf("Error updating %s: %v", r.URL, err))
	}
}
