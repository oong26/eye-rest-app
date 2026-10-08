package main

import (
	"strconv"
	"strings"
	"testing"
)

func TestAppVersionIsSemanticVersion(t *testing.T) {
	parts := strings.Split(appVersion, ".")
	if len(parts) != 3 {
		t.Fatalf("appVersion %q must use MAJOR.MINOR.PATCH", appVersion)
	}

	for _, part := range parts {
		if part == "" {
			t.Fatalf("appVersion %q must not contain empty version segments", appVersion)
		}
		if _, err := strconv.Atoi(part); err != nil {
			t.Fatalf("appVersion %q must contain numeric version segments: %v", appVersion, err)
		}
	}
}
