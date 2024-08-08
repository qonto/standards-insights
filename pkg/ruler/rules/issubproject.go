package rules

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/qonto/standards-insights/pkg/project"
)

type IsSubProjectRule struct {
}

func NewIsSubProjectRule() *IsSubProjectRule {
	return &IsSubProjectRule{}
}

func (rule *IsSubProjectRule) Do(ctx context.Context, project project.Project) error {
	if project.SubProject == "" {
		return errors.New("this project is not a subproject")
	}

	// Check if the SubProject path is a file
	fileInfo, err := os.Stat(filepath.Join(project.Path, project.SubProject))
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return errors.New("the SubProject path points to a directory, not a file")
	}

	return nil
}