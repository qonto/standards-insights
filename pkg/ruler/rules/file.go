package rules

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/qonto/standards-insights/config"
	"github.com/qonto/standards-insights/pkg/project"
	"github.com/qonto/standards-insights/types"
)

type FileRule struct {
	Path        string
	Contains    *types.Regexp
	NotContains *types.Regexp
	Exists      *bool
}

func NewFileRule(config config.FileRule) *FileRule {
	return &FileRule{
		Path:        config.Path,
		Contains:    config.Contains,
		NotContains: config.NotContains,
		Exists:      config.Exists,
	}
}

func (rule *FileRule) Do(ctx context.Context, project project.Project) error {
	path := filepath.Join(project.Path, rule.Path)
	if rule.Contains != nil {
		file, err := os.Open(path) //nolint
		if err != nil {
			return fmt.Errorf("fail to read file %s: %w", rule.Path, err)
		}
		defer file.Close() //nolint

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Bytes()
			if rule.Contains.Regexp.Match(line) {
				return nil
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error while reading file %s: %w", rule.Path, err)
		}

		return fmt.Errorf("pattern %s not found in file %s", rule.Contains, rule.Path)
	}
	if rule.NotContains != nil {
		file, err := os.Open(path) //nolint
		if err != nil {
			return fmt.Errorf("fail to read file %s: %w", rule.Path, err)
		}
		defer file.Close() //nolint

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Bytes()
			if rule.NotContains.Regexp.Match(line) {
				return fmt.Errorf("pattern %s found in file %s", rule.NotContains, rule.Path)
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error while reading file %s: %w", rule.Path, err)
		}

		return nil
	}
	if rule.Exists != nil {
		shouldExists := *rule.Exists
		_, err := os.Stat(path)
		isNotExist := os.IsNotExist(err)
		if err != nil && !isNotExist {
			return fmt.Errorf("unknown error while checking file %s: %w", rule.Path, err)
		}
		if shouldExists && isNotExist {
			return fmt.Errorf("file %s does not exist", rule.Path)
		}
		if !shouldExists && !isNotExist {
			return fmt.Errorf("file %s exists", rule.Path)
		}
		return nil
	}
	return nil
}
