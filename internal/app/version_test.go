package app

import (
	"strconv"
	"strings"
	"testing"
)

func TestDefaultVersionIsSemanticVersion(t *testing.T) {
	parts := strings.Split(DefaultVersion, ".")
	if len(parts) != 3 {
		t.Fatalf("DefaultVersion %q must use MAJOR.MINOR.PATCH", DefaultVersion)
	}

	for _, part := range parts {
		if part == "" {
			t.Fatalf("DefaultVersion %q must not contain empty version segments", DefaultVersion)
		}
		if _, err := strconv.Atoi(part); err != nil {
			t.Fatalf("DefaultVersion %q must contain numeric version segments: %v", DefaultVersion, err)
		}
	}
}
