package git

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

type PrivateKey struct {
	Path     string
	Password string
}

type BasicAuth struct {
	Username string
	Password string
}

type Config struct {
	PrivateKey PrivateKey `yaml:"private-key" validate:"omitempty,file"`
	BasicAuth  BasicAuth  `yaml:"basic-auth"`
}

type Git struct {
	logger     *slog.Logger
	config     Config
	publicKeys *ssh.PublicKeys
}

func New(logger *slog.Logger, config Config) (*Git, error) {
	keyPassword := os.Getenv("GIT_PRIVATE_KEY_PASSWORD")
	if keyPassword != "" {
		config.PrivateKey.Password = keyPassword
		err := os.Unsetenv("GIT_PRIVATE_KEY_PASSWORD")
		if err != nil {
			return nil, err
		}
	}

	basicAuthPassword := os.Getenv("GIT_BASIC_AUTH_PASSWORD")
	if basicAuthPassword != "" {
		config.BasicAuth.Password = basicAuthPassword
		err := os.Unsetenv("GIT_BASIC_AUTH_PASSWORD")
		if err != nil {
			return nil, err
		}
	}

	component := &Git{
		logger: logger,
		config: config,
	}

	if config.PrivateKey.Path != "" {
		key, err := ssh.NewPublicKeysFromFile("git", config.PrivateKey.Path, config.PrivateKey.Password)
		if err != nil {
			return nil, err
		}
		component.publicKeys = key
	}
	return component, nil
}

func (g *Git) Clone(url string, ref string, path string) error {
	g.logger.Debug(fmt.Sprintf("cloning repository %s", url))
	options := &git.CloneOptions{
		URL: url,
		// Depth:         1, shallow clone
		SingleBranch:  true,
		ReferenceName: plumbing.ReferenceName(ref),
	}
	if g.publicKeys != nil {
		options.Auth = g.publicKeys
	}
	if g.config.BasicAuth.Username != "" {
		options.Auth = &http.BasicAuth{
			Username: g.config.BasicAuth.Username,
			Password: g.config.BasicAuth.Password,
		}
	}
	_, err := git.PlainClone(path, false, options)
	return err
}

func (g *Git) Pull(path string, ref string) error {
	g.logger.Debug(fmt.Sprintf("pulling repository %s ref %s", path, ref))
	repository, err := git.PlainOpen(path)
	if err != nil {
		return err
	}
	worktree, err := repository.Worktree()
	if err != nil {
		return err
	}
	options := &git.PullOptions{
		SingleBranch: true,
	}
	if g.publicKeys != nil {
		options.Auth = g.publicKeys
	}
	if g.config.BasicAuth.Username != "" {
		options.Auth = &http.BasicAuth{
			Username: g.config.BasicAuth.Username,
			Password: g.config.BasicAuth.Password,
		}
	}
	err = worktree.Pull(options)
	if err != nil && errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}
	return nil
}
