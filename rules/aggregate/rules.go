package aggregate

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
)

var ErrRuleFailed = errors.New("rule verification failed")

type RuleType interface {
	Verify(ctx context.Context) error
	Summary() string
}

type Rule struct {
	Name  string
	Files []FileRule
}

type FileRule struct {
	Path     string
	Contains *Regexp
	Exists   *bool
}

func (r *FileRule) Verify(ctx context.Context) error {

	if r.Contains != nil {
		file, err := os.Open(r.Path)
		if err != nil {
			return fmt.Errorf("fail to read file %s: %w", r.Path, err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Bytes()
			if r.Contains.Regexp.Match(line) {
				return nil
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error while reading file %s: %w", r.Path, err)
		}

		return ErrRuleFailed
	}
	if r.Exists != nil {
		_, err := os.Stat(r.Path)
		if err != nil {
			return ErrRuleFailed
		}
		return nil
	}
	return nil
}

func (r *FileRule) Summary() string {
	if r.Contains != nil {
		return fmt.Sprintf("file %s should contain %s", r.Path, r.Contains.Regexp.String())
	} else {
		return fmt.Sprintf("file %s should exist", r.Path)
	}
}

type RuleError struct {
	RuleName string
	Messages []string
}

func (e *RuleError) Error() string {
	return fmt.Sprintf("rule %s failed: %s", e.RuleName, strings.Join(e.Messages, "-"))
}
