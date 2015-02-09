package image

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"core-gitlab.corp.zulily.com/core/stevedore/repo"
	"core-gitlab.corp.zulily.com/core/stevedore/ui"
)

func imageName(r *repo.Repo, registry string) string {
	urlTokens := strings.Split(strings.TrimSuffix(r.URL, ".git"), "/")
	imgTokens := []string{registry}
	imgTokens = append(imgTokens, urlTokens[3:]...)
	img := strings.Join(imgTokens, "/")
	if strings.Count(img, "/") > 2 {
		imgTokens = strings.SplitN(img, "/", 3)
		imgTokens[2] = strings.Replace(imgTokens[2], "/", "-", -1)
		img = strings.Join(imgTokens, "/")
	}
	return img
}

func versionToTag(version string) string {
	return version[0:8]
}

// Make will run the `make` command in the repository root directory if there
// is a Makefile there.
func Make(r *repo.Repo) error {
	makefile := filepath.Join(r.LocalPath(), "Makefile")
	if _, err := os.Stat(makefile); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	makeCmd := prepareMakeCommand(r.LocalPath(), "make")
	if err := makeCmd.Run(); err != nil {
		return err
	}

	return nil
}

// Build creates a docker image as specified by the Dockerfile in the repo root
// path.
func Build(r *repo.Repo, version, registry string) (name string, err error) {
	dockerfile := filepath.Join(r.LocalPath(), "Dockerfile")
	if _, err := os.Stat(dockerfile); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("Cannot build %s, no Dockerfile found in root of repository", r.URL)
		}
		return "", err
	}

	nameAndTag := imageName(r, registry) + ":" + versionToTag(version)

	buildCmd := prepareDockerCommand(r.LocalPath(), "docker", "build", "-t", nameAndTag, ".")
	if err := buildCmd.Run(); err != nil {
		return "", err
	}

	return nameAndTag, nil
}

// Publish pushes a local docker image to its registry via `gcloud preview docker push`.
func Publish(image string) error {
	cmdAndArgs := repo.PublishCommand(image)
	cmd := cmdAndArgs[0]
	args := cmdAndArgs[1:]
	publishCmd := prepareGcloudCommand(cmd, args...)
	return publishCmd.Run()
}

func prepareMakeCommand(path, cmd string, args ...string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.Dir = path
	c.Stdout = ui.Wrap(os.Stdout)
	c.Stderr = c.Stdout
	return c
}

func prepareDockerCommand(path, cmd string, args ...string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.Dir = path
	c.Stdout = ui.Wrap(os.Stdout)
	c.Stderr = c.Stdout
	return c
}

func prepareGcloudCommand(cmd string, args ...string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.Stdout = ui.Wrap(os.Stdout)
	c.Stderr = c.Stdout
	return c
}
