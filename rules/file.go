package rules

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/qonto/standards-insights/rules/aggregates"
)

func executeFileRule(_ context.Context, fileRule aggregates.FileRule) error {
	if fileRule.Contains != nil {
		file, err := os.Open(fileRule.Path)
		if err != nil {
			return fmt.Errorf("fail to read file %s: %w", fileRule.Path, err)
		}
		defer file.Close() //nolint

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Bytes()
			if fileRule.Contains.Regexp.Match(line) {
				return nil
			}
		}

		if err := scanner.Err(); err != nil {
			return fmt.Errorf("error while reading file %s: %w", fileRule.Path, err)
		}

		return fmt.Errorf("pattern %s not found in file %s", fileRule.Contains, fileRule.Path)
	}
	if fileRule.Exists != nil {
		exists := *fileRule.Exists
		_, err := os.Stat(fileRule.Path)
		if exists {
			if err != nil {
				return fmt.Errorf("file %s does not exist", fileRule.Path)
			}
		} else {
			if !os.IsNotExist(err) {
				return fmt.Errorf("unknown error while checking file %s: %w", fileRule.Path, err)
			}
		}
	}
	return nil
}

func executeFileRules(ctx context.Context, rule aggregates.Rule) error {
	errorMessages := []string{}
	if len(rule.Files) != 0 {
		for _, fileRule := range rule.Files {
			err := executeFileRule(ctx, fileRule)
			if err != nil {
				message := err.Error()
				errorMessages = append(errorMessages, message)
			}
		}
	}
	if len(errorMessages) > 0 {
		return errors.New(strings.Join(errorMessages, ", "))
	}
	return nil
}
