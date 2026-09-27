package version

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Fatal("expected non-empty Version")
	}
	full := Full()
	if !strings.Contains(full, "spp-shadowsocks-plugin version") {
		t.Fatalf("unexpected full version string: %s", full)
	}
	short := Short()
	if !strings.Contains(short, "spp-shadowsocks-plugin") {
		t.Fatalf("unexpected short version string: %s", short)
	}
}
