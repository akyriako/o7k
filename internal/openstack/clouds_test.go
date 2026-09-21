package openstack

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadCloudsMultipleFiles(t *testing.T) {
	dir := t.TempDir()

	firstPath := filepath.Join(dir, "first.yaml")
	secondPath := filepath.Join(dir, "second.yaml")

	first := `
clouds:
  prod:
    region_name: region-one
    auth:
      auth_url: https://prod.example.com
      project_name: prod-project
      domain_name: prod-domain
  shared:
    region_name: first-region
    auth:
      auth_url: https://first.example.com
      project_name: first-project
      domain_name: first-domain
`

	second := `
clouds:
  dev:
    region_name: region-two
    auth:
      auth_url: https://dev.example.com
      project_name: dev-project
      domain_name: dev-domain
  shared:
    region_name: second-region
    auth:
      auth_url: https://second.example.com
      project_name: second-project
      domain_name: second-domain
`

	if err := os.WriteFile(firstPath, []byte(first), 0600); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(secondPath, []byte(second), 0600); err != nil {
		t.Fatal(err)
	}

	clouds, err := LoadClouds([]string{firstPath, secondPath})
	if err != nil {
		t.Fatal(err)
	}

	if len(clouds.Items) != 3 {
		t.Fatalf("LoadClouds() returned %d clouds, want 3", len(clouds.Items))
	}

	wantNames := []string{"prod", "shared", "dev"}

	for i, want := range wantNames {
		if clouds.Items[i].Name != want {
			t.Fatalf(
				"clouds.Items[%d].Name = %q, want %q",
				i,
				clouds.Items[i].Name,
				want,
			)
		}
	}

	shared, ok := clouds.Get("shared")
	if !ok {
		t.Fatal(`cloud "shared" not found`)
	}

	if shared.Region != "first-region" {
		t.Fatalf(
			`shared.Region = %q, want "first-region"`,
			shared.Region,
		)
	}

	if shared.Path != firstPath {
		t.Fatalf(
			"shared.Path = %q, want %q",
			shared.Path,
			firstPath,
		)
	}

	dev, ok := clouds.Get("dev")
	if !ok {
		t.Fatal(`cloud "dev" not found`)
	}

	if dev.Path != secondPath {
		t.Fatalf(
			"dev.Path = %q, want %q",
			dev.Path,
			secondPath,
		)
	}
}
