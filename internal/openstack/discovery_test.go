package openstack

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscoverCloudsFilesExplicitPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "clouds.yaml")

	if err := os.WriteFile(
		path,
		[]byte("clouds: {}\n"),
		0600,
	); err != nil {
		t.Fatal(err)
	}

	got, err := DiscoverCloudsFiles(path)
	if err != nil {
		t.Fatal(err)
	}

	want, err := filepath.Abs(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 {
		t.Fatalf("DiscoverCloudsFiles() returned %d paths, want 1", len(got))
	}

	if got[0] != want {
		t.Fatalf("DiscoverCloudsFiles()[0] = %q, want %q", got[0], want)
	}
}
