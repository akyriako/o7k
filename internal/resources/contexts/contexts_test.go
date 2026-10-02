package contexts

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/akyriako/o7k/internal/openstack"
)

func writeCloudsFile(t *testing.T, name string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "clouds.yaml")

	content := "clouds:\n  " + name + ":\n    region_name: region-one\n    auth:\n      auth_url: https://prod.example.com\n      project_name: prod-project\n      domain_name: prod-domain\n"

	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	return path
}

func TestListWithoutEnvironment(t *testing.T) {
	t.Setenv("OS_AUTH_URL", "")

	path := writeCloudsFile(t, "prod")

	rows, err := New([]string{path}).List(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(rows) != 1 || rows[0].ID != "prod" {
		t.Fatalf("List() = %+v, want only prod", rows)
	}

	if rows[0].Fields["cloudsyaml"] != path {
		t.Fatalf("cloudsyaml = %q, want %q", rows[0].Fields["cloudsyaml"], path)
	}
}

func TestListAppendsEnvironmentLast(t *testing.T) {
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com/v3")

	rows, err := New([]string{writeCloudsFile(t, "prod")}).List(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(rows) != 2 || rows[0].ID != "prod" || rows[1].ID != openstack.EnvCloudName {
		t.Fatalf("List() = %+v, want prod then %s", rows, openstack.EnvCloudName)
	}

	if rows[1].Fields["cloudsyaml"] != openstack.EnvCloudSource {
		t.Fatalf("cloudsyaml = %q, want %q", rows[1].Fields["cloudsyaml"], openstack.EnvCloudSource)
	}
}

func TestListSkipsEnvironmentOnNameCollision(t *testing.T) {
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com/v3")

	path := writeCloudsFile(t, openstack.EnvCloudName)

	rows, err := New([]string{path}).List(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(rows) != 1 || rows[0].ID != openstack.EnvCloudName {
		t.Fatalf("List() = %+v, want a single %s row", rows, openstack.EnvCloudName)
	}

	if rows[0].Fields["cloudsyaml"] != path {
		t.Fatalf("cloudsyaml = %q, want the clouds.yaml entry %q", rows[0].Fields["cloudsyaml"], path)
	}
}

func TestListEnvironmentOnly(t *testing.T) {
	t.Setenv("OS_AUTH_URL", "https://keystone.example.com/v3")

	rows, err := New(nil).List(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(rows) != 1 || rows[0].ID != openstack.EnvCloudName {
		t.Fatalf("List() = %+v, want a single %s row", rows, openstack.EnvCloudName)
	}
}
