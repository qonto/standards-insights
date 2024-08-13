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

	ownerMap := map[string]string{
		"@team-cft-account-management":          "account-management",
		"@team-cft-bookkeeping":                 "bookkeeping",
		"@team-cft-cards":                       "cards",
		"@team-cft-cash-flow-management":        "cash-flow-management",
		"@team-cft-connect":                     "connect",
		"@team-cft-data":                        "data",
		"@team-cft-data-products":               "data-products",
		"@team-cft-design-system":               "design-system",
		"@team-cft-financing":                   "financing",
		"@team-cft-invoices":                    "invoices",
		"@team-cft-ledger":                      "ledger",
		"@team-cft-local-payments":              "local-payments",
		"@team-cft-onboarding":                  "onboarding",
		"@team-cft-onboarding-company-creation": "onboarding-company-creation",
		"@team-cft-onboarding-kyb":              "onboarding-kyb",
		"@team-cft-onboarding-kyc":              "onboarding-kyc",
		"@team-cft-onboarding-registration":     "onboarding-registration",
		"@team-cft-ops-tooling":                 "ops-tooling",
		"@team-cft-pricing":                     "pricing",
		"@team-cft-security":                    "security",
		"@team-cft-sepa":                        "sepa",
		"@team-cft-spend-management":            "spend-management",
		"@team-cft-sre":                         "sre",
		"@team-cft-swift":                       "swift",
	}

	owner := owners[0]
	if mappedOwner, found := ownerMap[owner]; found {
		return mappedOwner, true
	}

	return owner, true
}
