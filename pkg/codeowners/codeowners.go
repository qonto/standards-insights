package codeowners

import (
	"os"

	"github.com/hairyhenderson/go-codeowners"
	"github.com/qonto/standards-insights/config"
)

type Codeowners struct {
	codeowners *codeowners.Codeowners
	configPath string
}

func NewCodeowners(path string, configPath string) (*Codeowners, error) {
	// open filesystem rooted at current working directory
	fsys := os.DirFS(path)

	c, err := codeowners.FromFileWithFS(fsys, "CODEOWNERS")
	if err != nil {
		return nil, err
	}

	return &Codeowners{
		codeowners: c,
		configPath: configPath,
	}, nil
}

func (c *Codeowners) GetOwners(path string) (string, bool) {
	if c == nil || c.codeowners == nil {
		return "", false
	}

	owners := c.codeowners.Owners(path)
	if len(owners) == 0 {
		return "", false
	}

	// Attempt to load the configuration
	config, _, err := config.New(c.configPath)
	var ownerMap map[string]string
	if err != nil {
		// If config cannot be found, use an empty ownerMap
		ownerMap = map[string]string{}
	} else {
		ownerMap = config.OwnerMap
	}

	owner := owners[0]
	if mappedOwner, found := ownerMap[owner]; found {
		return mappedOwner, true
	}

	return owner, true
}