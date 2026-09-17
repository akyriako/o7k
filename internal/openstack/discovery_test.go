package openstack

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverCloudsFileExplicitPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "clouds.yaml")

	if err := os.WriteFile(
		path,
		[]byte("clouds: {}\n"),
		0600,
	); err != nil {
		t.Fatal(err)
	}

	got, err := DiscoverCloudsFile(path)
	if err != nil {
		t.Fatal(err)
	}

	want, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Fatalf(
			"DiscoverCloudsFile() = %q, want %q",
			got,
			want,
		)
	}
}
