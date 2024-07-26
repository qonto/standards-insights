package rules

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/project"
)

type GrepRule struct {
	Path           string
	Include        string
	Pattern        string
	ExtendedRegexp bool
	Recursive      bool
	Match          bool
	SkipNotFound   bool
	NullData       bool
}

func NewGrepRule(config config.GrepRule) *GrepRule {
	return &GrepRule{
		Path:           config.Path,
		Recursive:      config.Recursive,
		Pattern:        config.Pattern,
		ExtendedRegexp: config.ExtendedRegexp,
		Include:        config.Include,
		Match:          config.Match,
		SkipNotFound:   config.SkipNotFound,
		NullData:       config.NullData,
	}
}

func (rule *GrepRule) Do(ctx context.Context, project project.Project) error {
	path := filepath.Join(project.Path, project.SubProject, rule.Path)
	_, err := os.Stat(path)
	isNotExist := os.IsNotExist(err)
	if isNotExist && rule.SkipNotFound {
		return nil
	}
	arguments := []string{}
	if rule.Recursive {
		arguments = append(arguments, "-r")
	}
	if rule.Include != "" {
		arguments = append(arguments, "--include", rule.Include)
	}
	if rule.ExtendedRegexp {
		arguments = append(arguments, "--extended-regexp")
	}
	if rule.NullData {
		arguments = append(arguments, "-z")
	}
	arguments = append(arguments, rule.Pattern, path)

	cmd := exec.CommandContext(ctx, "grep", arguments...) //nolint

	var stdErrBuffer bytes.Buffer
	var stdOutBuffer bytes.Buffer
	cmd.Stderr = &stdErrBuffer
	cmd.Stdout = &stdOutBuffer

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode := exitErr.ExitCode()
			// grep returns 1 if no match
			if exitCode == 1 && !rule.Match {
				return nil
			}
			if exitCode == 1 && rule.Match {
				return fmt.Errorf("no match for pattern %s on path %s", rule.Pattern, filepath.Join(project.SubProject, rule.Path))
			}
			return fmt.Errorf("failed to execute grep command (error code %d), stderr=%s, error=%w", exitErr.ExitCode(), stdErrBuffer.String(), err)
		}
		return fmt.Errorf("the grep command failed, stderr=%s, error=%w", stdErrBuffer.String(), err)
	}
	// exit code is zero so lines matching the pattern were detected
	if !rule.Match {
		return fmt.Errorf("match found for pattern %s on path %s", rule.Pattern, rule.Path)
	}
	return nil
}
