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

func imageName(r *repo.Repo, registry string, dockerfile string) string {
	urlTokens := strings.Split(strings.TrimSuffix(r.URL, ".git"), "/")
	imgTokens := []string{registry}
	imgTokens = append(imgTokens, urlTokens[3:]...)
	img := strings.Join(imgTokens, "/")
	if strings.Count(img, "/") > 2 {
		imgTokens = strings.SplitN(img, "/", 3)
		imgTokens[2] = strings.Replace(imgTokens[2], "/", "-", -1)
		img = strings.Join(imgTokens, "/")
	}

	fname := filepath.Base(dockerfile)
	if strings.HasPrefix(fname, "Dockerfile.") {
		suffix := fname[len("Dockerfile."):]
		img = strings.Join([]string{img, suffix}, "-")
	}

	return img
}

func versionToTag(version string) string {
	return version[0:7]
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

	makeCmd := prepareCommand(r.LocalPath(), "make")
	if err := makeCmd.Run(); err != nil {
		return err
	}

	return nil
}

// Build creates one or more docker images, as specified by the Dockerfile(s) in the repo root
// path.  Valid Dockerfiles are either named 'Dockerfile' or use the naming convention 'Dockerfile.<SUFFIX>'
func Build(r *repo.Repo, version, registry string) (name []string, err error) {

	var names []string
	dockerfiles, err := filepath.Glob(filepath.Join(r.LocalPath(), "Dockerfile*"))
	if err == filepath.ErrBadPattern {
		return names, err
	}

	if dockerfiles == nil {
		return names, fmt.Errorf("Cannot build %s, no Dockerfile(s) found in root of repository", r.URL)
	}

	for _, dockerfile := range dockerfiles {
		nameAndTag := imageName(r, registry, dockerfile) + ":" + versionToTag(version)

		buildCmd := prepareCommand(r.LocalPath(), "docker", "build", "-f", dockerfile, "-t", nameAndTag, ".")

		if err := buildCmd.Run(); err != nil {
			return names, err
		}

		names = append(names, nameAndTag)
	}

	return names, nil
}

// Publish pushes a local docker image to its registry via `gcloud preview docker push`.
func Publish(publishCmd []string) error {
	cmd := publishCmd[0]
	args := publishCmd[1:]
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	return prepareCommand(wd, cmd, args...).Run()
}

func prepareCommand(path, cmd string, args ...string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.Dir = path
	c.Stdout = ui.Wrap(os.Stdout)
	c.Stderr = c.Stdout
	return c
}
