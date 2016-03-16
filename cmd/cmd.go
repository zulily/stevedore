package cmd

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

var (
	Registry string
	Output   io.Writer  = ioutil.Discard
	Filter   FilterFunc = matchAll
	Tag      string
)

type FilterFunc func(dockerfile string) bool

func matchAll(dockerfile string) bool {
	return true
}

func matchAny(dockerfiles ...string) FilterFunc {
	return func(dockerfile string) bool {
		for _, v := range dockerfiles {
			if dockerfile == v {
				return true
			}
		}

		return false
	}
}

func matchRegexp(expr string) FilterFunc {
	rexpr := regexp.MustCompile(expr)
	return func(dockerfile string) bool {
		return rexpr.MatchString(dockerfile)
	}
}

func init() {
	var expr string
	verbose := false

	flag.StringVar(&expr, "i", "", "include only dockerfiles that match this regular expression")
	flag.StringVar(&Registry, "registry-base", "docker.io", "the registry name to prepend to each Docker image")
	flag.BoolVar(&verbose, "verbose", false, "enables verbose output")
	flag.BoolVar(&Tag, "tag", "", "manually specify a tag")

	flag.Parse()

	if verbose {
		Output = os.Stdout
	}

	args := flag.Args()

	switch {
	case expr != "" && len(args) == 0:
		Filter = matchRegexp(expr)
	case expr == "" && len(args) > 0:
		Filter = matchAny(args[0:]...)
	case expr != "" && len(args) > 0:
		log.Fatal("Cannot mix -i and dockerfile args")
	}
}
