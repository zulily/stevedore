package image

import (
	"strings"
	"testing"

	"github.com/zulily/stevedore/repo"
)

func TestNonHttpsImageName(t *testing.T) {
	registry := "gcr.io/eternal_empire_754"
	repo := &repo.Repo{
		URL:    "git@github.com:dcarney/actually-test.git",
		SHA:    "ecaf0d06834ec132fedd74a61a3e3871367c5833",
		Images: []string{},
	}
	img := imageName(repo, registry, "Dockerfile")

	if strings.Count(img, "/") > 2 {
		t.Error("Docker image name contains too many '/' separators")
	}

	if img != "gcr.io/eternal_empire_754/dcarney-actually-test" {
		t.Errorf("Expected gcr.io/eternal_empire_754/dcarney-actually-test, got %q", img)
	}
}

func TestImageName(t *testing.T) {
	registry := "gcr.io/eternal_empire_754"
	repo := &repo.Repo{
		URL:    "https://github.com/dcarney/actually-test.git",
		SHA:    "ecaf0d06834ec132fedd74a61a3e3871367c5833",
		Images: []string{},
	}
	img := imageName(repo, registry, "Dockerfile")

	if strings.Count(img, "/") > 2 {
		t.Error("Docker image name contains too many '/' separators")
	}

	if img != "gcr.io/eternal_empire_754/dcarney-actually-test" {
		t.Error("Docker image name is incorrect")
	}
}

func TestImageNameWithNonDefaultDockerfile(t *testing.T) {
	registry := "gcr.io/eternal_empire_754"
	repo := &repo.Repo{
		URL:    "https://github.com/dcarney/actually-test.git",
		SHA:    "ecaf0d06834ec132fedd74a61a3e3871367c5833",
		Images: []string{},
	}
	dockerfile := "Dockerfile.foobar"

	img := imageName(repo, registry, dockerfile)

	if img != "gcr.io/eternal_empire_754/dcarney-actually-test-foobar" {
		t.Error("Docker image name is incorrect when using non-default Dockerfile name")
	}

}
