package codeowners

import (
	"os"
	"testing"
)

func TestCodeowners_FromFile(t *testing.T) {
	// Ensure the CODEOWNERS file exists
	if _, err := os.Stat(".gitlab/CODEOWNERS"); os.IsNotExist(err) {
		t.Fatalf("CODEOWNERS file does not exist")
	}

	co, err := NewCodeowners(".")
	if err != nil {
		t.Fatalf("Error creating codeowners: %v", err)
	}

	tests := []struct {
		path     string
		expected string
	}{
		{"app/announcers/attachment_announcer.rb", "bookkeeping"},
		{"app/controllers/v1/external_transfers_controller.rb", "sepa"},
		{"app/controllers/somepath/v1/external_transfers_controller.rb", "sepa"},
		{"app/controllers/anotherpath/mandates_controller.rb", "sepa"},
		{"app/services/transfers/somefile.rb", "sepa"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			owner, exists := co.GetOwners(tt.path)
			if !exists {
				t.Errorf("GetOwners() for path %s, expected owner %s, but got none", tt.path, tt.expected)
			}
			if owner != tt.expected {
				t.Errorf("GetOwners() for path %s, expected owner %s, but got %s", tt.path, tt.expected, owner)
			}
		})
	}
}