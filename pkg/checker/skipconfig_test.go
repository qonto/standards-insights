package checker

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeSkipFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, skipConfigFileName), []byte(content), 0o600)
	require.NoError(t, err)
	return dir
}

func TestLoadSkipConfigMissingFile(t *testing.T) {
	groups, checks, err := loadSkipConfig(t.TempDir())
	assert.NoError(t, err)
	assert.Empty(t, groups)
	assert.Empty(t, checks)
}

func TestLoadSkipConfigGroupsAndChecks(t *testing.T) {
	dir := writeSkipFile(t, `
skip:
  groups:
    - golang
  checks:
    - go-version-latest
    - go-main
`)
	groups, checks, err := loadSkipConfig(dir)
	assert.NoError(t, err)
	assert.Equal(t, map[string]bool{"golang": true}, groups)
	assert.Equal(t, map[string]bool{"go-version-latest": true, "go-main": true}, checks)
}

func TestLoadSkipConfigOnlyChecks(t *testing.T) {
	dir := writeSkipFile(t, `
skip:
  checks:
    - go-version-latest
`)
	groups, checks, err := loadSkipConfig(dir)
	assert.NoError(t, err)
	assert.Empty(t, groups)
	assert.Equal(t, map[string]bool{"go-version-latest": true}, checks)
}

func TestLoadSkipConfigOnlyGroups(t *testing.T) {
	dir := writeSkipFile(t, `
skip:
  groups:
    - golang
`)
	groups, checks, err := loadSkipConfig(dir)
	assert.NoError(t, err)
	assert.Equal(t, map[string]bool{"golang": true}, groups)
	assert.Empty(t, checks)
}

func TestLoadSkipConfigMalformed(t *testing.T) {
	dir := writeSkipFile(t, "skip: [this is not valid: yaml")
	groups, checks, err := loadSkipConfig(dir)
	assert.Error(t, err)
	// Fails open: nothing is skipped when the file cannot be parsed.
	assert.Empty(t, groups)
	assert.Empty(t, checks)
}
