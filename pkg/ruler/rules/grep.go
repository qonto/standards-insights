package rules

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/project"
)

type GrepRule struct {
	Path      string
	Recursive bool
	Pattern   string
	Match     bool
}

func NewGrepRule(config config.GrepRule) *GrepRule {
	return &GrepRule{
		Path:      config.Path,
		Recursive: config.Recursive,
		Pattern:   config.Pattern,
		Match:     config.Match,
	}
}

func (rule *GrepRule) Do(ctx context.Context, project project.Project) error {
	arguments := []string{}
	if rule.Recursive {
		arguments = append(arguments, "-r")
	}
	arguments = append(arguments, rule.Pattern, rule.Path)

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
				return fmt.Errorf("no match for pattern %s on path %s", rule.Pattern, rule.Path)
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
