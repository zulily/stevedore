package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var (
	// mu guards access to repos (and other data structures eventually)
	mu  sync.Mutex
	cfg Config
)

// Config wraps the entire contents of the "repos.json" file.
type Config struct {
	Repos []*Repo `json:"repos"`
}

type BuildStatus int

const (
	StatusPassing BuildStatus = iota
	StatusInProgress
	StatusFailing
)

// Repo represents a git source code repository.
type Repo struct {
	URL             string            `json:"url"`
	SHA             string            `json:"sha"`
	Images          []string          `json:"images"`
	LastPublishDate int64             `json:"lastPublishDate"`
	Log             string            `json:"log"`
	Status          BuildStatus       `json:"buildStatus"`
	ExternalFiles   map[string]string `json:"external,omitempty"`
}

// Validate ensures that the Repo instance has a valid configuration
func (r *Repo) Validate() error {
	if strings.Index(r.URL, "https://") != 0 {
		return fmt.Errorf("Invalid repo URL: '%s', only https is supported", r.URL)
	}
	return nil
}

// LocalPath returns the location on the local file-system where this repo will
// be synced.
func (r *Repo) LocalPath() string {
	id := strings.NewReplacer("https://", "", "http://", "", "/", "_", ":", "_").Replace(r.URL)
	wd, err := os.Getwd()
	if err != nil {
		wd = os.TempDir()
	}
	local := filepath.Join(wd, "builds", id)
	return local
}

// All returns all repositories that Stevedore needs to sync and build.
func All() ([]*Repo, error) {
	mu.Lock()
	defer mu.Unlock()

	jsonFile := filepath.Clean("./repos.json")
	file, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}

	cfg = Config{}
	err = json.Unmarshal(file, &cfg)
	return cfg.Repos, err
}

// Add adds a new repo to the global repo configuration
func Add(r *Repo) ([]*Repo, error) {
	mu.Lock()
	defer mu.Unlock()

	if r == nil {
		return cfg.Repos, nil
	}

	cfg.Repos = append(cfg.Repos, r)

	if err := doSave(); err != nil {
		return nil, err
	}

	return cfg.Repos, nil
}

// Save updates the Stevedore configuration
func (r *Repo) Save() error {
	mu.Lock()
	defer mu.Unlock()

	return doSave()
}

func doSave() error {
	bytes, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		fmt.Printf("ERROR saving: %#v\n", err)
		return err
	}

	jsonFile := filepath.Clean("./repos.json")
	return ioutil.WriteFile(jsonFile, bytes, 0644)
}

// Checkout clones or fetches latest code from a remote repo, ensures the local
// copy is pointed at origin/master, and returns the SHA of the head revision.
func (r *Repo) Checkout() (head string, err error) {
	if r.URL == "" {
		return "", fmt.Errorf("Repo has empty URL")
	}

	local := r.LocalPath()
	if err := os.MkdirAll(local, 0755); err != nil {
		return "", err
	}

	if _, err := os.Stat(filepath.Join(local, ".git")); os.IsNotExist(err) {
		return clone(r.URL, local)
	}

	cmd := prepareGitCommand(local, "git", "clean", "-d", "-f", "-x")
	if err := cmd.Run(); err != nil {
		return "", err
	}

	cmd = prepareGitCommand(local, "git", "fetch", "--all")
	if err := cmd.Run(); err != nil {
		return "", err
	}

	cmd = prepareGitCommand(local, "git", "merge", "origin/master")
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return getHead(local)
}

func (r *Repo) PrepareMake() error {
	if len(r.ExternalFiles) == 0 {
		return nil
	}

	for src, dest := range r.ExternalFiles {
		dest = filepath.Join(r.LocalPath(), dest)
		fmt.Println("Copying", src, "to", dest)
		in, err := os.Open(src)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.Create(dest)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, in)
		cerr := out.Close()
		if err != nil {
			return err
		}
		if cerr != nil {
			return cerr
		}
	}

	return nil
}

func clone(url, dest string) (head string, err error) {
	cmd := prepareGitCommand(path.Dir(dest), "git", "clone", url, dest)
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return getHead(dest)
}

func getHead(local string) (head string, err error) {
	cmd := prepareGitCommand(local, "git", "rev-parse", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err == nil {
		return strings.Trim(out.String(), "\r\n\t "), nil
	}

	return "", err
}

func prepareGitCommand(dir, cmd string, args ...string) *exec.Cmd {
	c := exec.Command(cmd, args...)
	c.Env = []string{"GIT_SSL_NO_VERIFY=true"}
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = c.Stdout
	return c
}
