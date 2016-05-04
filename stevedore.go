package stevedore

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zulily/stevedore/cmd"
)

type Image struct {
	Dockerfile string
	Url        string
}

func FindImagesInCwd(filter cmd.FilterFunc) ([]Image, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return findImages(filter, wd)
}

func findImages(filter cmd.FilterFunc, wd string) (images []Image, err error) {
	repo, path, tag := detectRepoPathAndTag(wd)
	dockerfiles := findDockerfiles()
	for dockerfile, repos := range mapDockerfileToRepos(repo, path, tag, dockerfiles...) {
		if !filter(dockerfile) {
			continue
		}

		for _, repo := range repos {
			img := Image{
				Dockerfile: dockerfile,
				Url:        repo,
			}
			images = append(images, img)
		}
	}
	return images, nil
}

func (i Image) String() string {
	return i.Url
}

func (i Image) Build() (err error) {
	return runCmdAndPipeOutput(cmd.Output, "docker", "build", "-t", i.Url, "-f", i.Dockerfile, ".")
}

func (i Image) Push() (err error) {
	return runCmdAndPipeOutput(cmd.Output, "docker", "push", i.Url)
}

func detectRepoPathAndTag(wd string) (repo, path, tag string) {
	repo, err := runCmdAndGetOutput("git", "config", "--get", "remote.origin.url")
	if err != nil {
		log.Fatal("error detecting git repo", err)
	}

	if index := strings.LastIndex(repo, ":"); index != -1 {
		repo = repo[index+1:]
	}

	if strings.HasSuffix(repo, ".git") {
		repo = repo[:len(repo)-4]
	}

	path, err = runCmdAndGetOutput("git", "rev-parse", "--show-toplevel")
	switch {
	case wd == path:
		path = ""
	case strings.HasPrefix(wd, path):
		path = wd[len(path)+1:]
		path = strings.Replace(path, "/", "-", -1)
	default:
		log.Fatal("Current directory is not child of top level", wd, path)
	}

	if cmd.Tag == "" {

		tag, err = runCmdAndGetOutput("git", "rev-parse", "HEAD")
		if err != nil {
			log.Fatal("error detecting git HEAD revision", err)
		}
	} else {
		tag = cmd.Tag
	}

	if len(tag) > 7 {
		tag = tag[:7]
	}

	return repo, path, tag
}

func runCmdAndPipeOutput(w io.Writer, name string, arg ...string) error {
	fmt.Println(">", name, strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)

	cmd.Stdout, cmd.Stderr = w, w
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func runCmdAndGetOutput(name string, arg ...string) (string, error) {
	log.Println(">", name, strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func findDockerfiles() []string {
	matches, err := filepath.Glob("Dockerfile*")
	if err != nil {
		log.Fatal("error finding Dockerfiles", err)
	}
	return matches
}

func mapDockerfileToRepos(base, path, tag string, dockerfile ...string) map[string][]string {
	m := make(map[string][]string)
	for _, f := range dockerfile {
		m[f] = generateRepoNames(base, path, tag, f)
	}
	return m
}

func generateRepoNames(base, path, tag, dockerfile string) []string {
	if strings.HasSuffix(cmd.Registry, "/") {
		base = cmd.Registry + base
	} else {
		base = cmd.Registry + "/" + base
	}

	if path != "" {
		base = base + "-" + path
	}

	name := base

	// grab the suffix from the Dockerfile, if any (e.g. "Dockerfile.foo" => "foo")
	suffix := dockerfile
	if index := strings.LastIndex(suffix, "."); index != -1 {
		name = name + "-" + suffix[index+1:]
	}

	// Docker image names can't have more than 2 '/' chars in them ¯\_(ツ)_/¯
	// replace any offending '/' chars w/ '-'
	if strings.Count(name, "/") > 2 {
		nameTokens := strings.SplitN(name, "/", 3)
		nameTokens[2] = strings.Replace(nameTokens[2], "/", "-", -1)
		name = strings.Join(nameTokens, "/")
	}

	if cmd.NoLatest {
		return []string{name + ":" + tag}
	}

	return []string{name + ":" + tag, name + ":latest"}
}
