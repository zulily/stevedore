package image

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/zulily/stevedore/repo"
)

func imageName(r *repo.Repo, registry string, dockerfile string) string {
	urlTokens := strings.Split(strings.TrimSuffix(r.URL, ".git"), "/")
	imgTokens := []string{registry}
	// 3 tokens assumes "https:" "" "hostname" ... then the rest
	discard := 3
	if !strings.HasPrefix(r.URL, "http") {
		discard = 1
	}
	imgTokens = append(imgTokens, urlTokens[discard:]...)
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

	return strings.ToLower(img)
}

func versionToTag(version string) string {
	return version[0:7]
}

// Make will run the `make` command in the repository root directory if there
// is a Makefile there.  It returns the combined stdout/stderr from the
// execution of the command, along with any error that may have occurred
func Make(r *repo.Repo) (string, error) {
	makefile := filepath.Join(r.LocalPath(), "Makefile")
	if _, err := os.Stat(makefile); err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	return execAndCapture(r.LocalPath(), "make")
}

// BuildResult encapsulates the status of a build attempt for a repo
type BuildResult struct {

	// ImageName is the full name of the Docker image that was built
	ImageName string

	// Output is the combined stdout/stderr of the build process
	Output string

	// Err is any error that may have occurred during the image building
	Err error
}

// Build creates one or more docker images, as specified by the Dockerfile(s)
// in the repo root path.  Valid Dockerfiles are either named 'Dockerfile' or
// use the naming convention 'Dockerfile.<SUFFIX>'.  The func results the
// results from each build, along with any error that may have occurred.  Note
// that this error is a general error, and is different from the `Err` property
// of each returned build result.
func Build(r *repo.Repo, version, registry string) ([]*BuildResult, error) {
	var results []*BuildResult
	dockerfiles, err := filepath.Glob(filepath.Join(r.LocalPath(), "Dockerfile*"))
	if err == filepath.ErrBadPattern {
		return results, err
	}

	if dockerfiles == nil {
		return results, fmt.Errorf("Cannot build %s, no Dockerfile(s) found in root of repository", r.URL)
	}

	for _, dockerfile := range dockerfiles {
		nameAndTag := imageName(r, registry, dockerfile) + ":" + versionToTag(version)

		output, err := execAndCapture(r.LocalPath(), "docker", "build", "--force-rm", "-f", dockerfile, "-t", nameAndTag, ".")
		result := &BuildResult{ImageName: nameAndTag, Output: output, Err: err}
		results = append(results, result)
	}

	return results, nil
}

// Publish pushes a local docker image to its registry, using the specified publish command.  For
// example: `gcloud preview docker push`, or `docker push`.
func Publish(image string, publishCmd []string) (string, error) {
	cmd := publishCmd[0]
	args := append(publishCmd[1:], image)
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return execAndCapture(wd, cmd, args...)
}

// execAndCapture execs the given command, returning the combined stdout/stderr
// of the command, along with any error that may have occurred.  Additionally,
// the same combined output is streamed to stdout while the command is
// executing. This behavior is analogous to the POSIX `tee` command.
func execAndCapture(path, cmd string, args ...string) (string, error) {
	c := exec.Command(cmd, args...)
	c.Dir = path
	stdout, err := c.StdoutPipe()
	if err != nil {
		return "", err
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	// capture stdout and stderr with the same reader
	r := io.MultiReader(stdout, stderr)
	w := io.MultiWriter(os.Stdout, &buf)

	if err := c.Start(); err != nil {
		return buf.String(), err
	}

	if _, err := io.Copy(w, r); err != nil {
		fmt.Println(err.Error())
	}

	if err := c.Wait(); err != nil {
		return buf.String(), err
	}

	return buf.String(), nil
}
