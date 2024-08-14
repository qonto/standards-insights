package codeowners

import (
	"os"

	"github.com/hairyhenderson/go-codeowners"
)

type Codeowners struct {
	codeowners *codeowners.Codeowners
}

func NewCodeowners(path string) (*Codeowners, error) {
	// open filesystem rooted at current working directory
	fsys := os.DirFS(path)

	c, err := codeowners.FromFileWithFS(fsys, "CODEOWNERS")
	if err != nil {
		return nil, err
	}

	return &Codeowners{
		codeowners: c,
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

	// TODO: load ownermap from standard-insights local depoyment file
	ownerMap := map[string]string{}

	owner := owners[0]
	if mappedOwner, found := ownerMap[owner]; found {
		return mappedOwner, true
	}

	return owner, true
}
