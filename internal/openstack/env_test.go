package openstack

import (
	"context"
	"strings"
	"testing"
)

// envVariables lists every variable read by env.go so that tests can clear
// them all and not depend on the developer's shell.
var envVariables = []string{
	"OS_AUTH_URL",
	"OS_USERNAME",
	"OS_USER_ID",
	"OS_USERID",
	"OS_PASSWORD",
	"OS_PASSCODE",
	"OS_USER_DOMAIN_ID",
	"OS_USER_DOMAIN_NAME",
	"OS_PROJECT_DOMAIN_ID",
	"OS_PROJECT_DOMAIN_NAME",
	"OS_DOMAIN_ID",
	"OS_DOMAIN_NAME",
	"OS_PROJECT_ID",
	"OS_PROJECT_NAME",
	"OS_TENANT_ID",
	"OS_TENANT_NAME",
	"OS_TOKEN",
	"OS_APPLICATION_CREDENTIAL_ID",
	"OS_APPLICATION_CREDENTIAL_NAME",
	"OS_APPLICATION_CREDENTIAL_SECRET",
	"OS_REGION_NAME",
}

func setEnv(t *testing.T, vars map[string]string) {
	t.Helper()

	for _, name := range envVariables {
		t.Setenv(name, "")
	}

	for name, value := range vars {
		t.Setenv(name, value)
	}
}

func TestEnvCloudNotConfigured(t *testing.T) {
	setEnv(t, nil)

	if _, ok := EnvCloud(); ok {
		t.Fatal("EnvCloud() reported a cloud without OS_AUTH_URL")
	}
}

func TestEnvCloudOpenRC(t *testing.T) {
	setEnv(t, map[string]string{
		"OS_AUTH_URL":            "https://keystone.example.com/v3",
		"OS_USERNAME":            "alice",
		"OS_PASSWORD":            "secret",
		"OS_PROJECT_NAME":        "research",
		"OS_USER_DOMAIN_NAME":    "Default",
		"OS_PROJECT_DOMAIN_NAME": "Default",
		"OS_REGION_NAME":         "region-one",
	})

	cloud, ok := EnvCloud()
	if !ok {
		t.Fatal("EnvCloud() reported no cloud with OS_AUTH_URL set")
	}

	want := Cloud{
		Name:     EnvCloudName,
		Region:   "region-one",
		Project:  "research",
		Domain:   "Default",
		Identity: "https://keystone.example.com/v3",
	}

	if cloud != want {
		t.Fatalf("EnvCloud() = %+v, want %+v", cloud, want)
	}

	opts := envAuthOptions()

	if opts.IdentityEndpoint != "https://keystone.example.com/v3" {
		t.Fatalf("IdentityEndpoint = %q", opts.IdentityEndpoint)
	}

	if opts.Username != "alice" || opts.Password != "secret" {
		t.Fatalf("Username = %q, Password = %q", opts.Username, opts.Password)
	}

	if opts.DomainName != "Default" || opts.DomainID != "" {
		t.Fatalf("DomainName = %q, DomainID = %q, want Default and empty", opts.DomainName, opts.DomainID)
	}

	if opts.TenantName != "research" || opts.TenantID != "" {
		t.Fatalf("TenantName = %q, TenantID = %q, want research and empty", opts.TenantName, opts.TenantID)
	}

	if !opts.AllowReauth {
		t.Fatal("AllowReauth = false, want true")
	}
}

func TestEnvAuthOptionsHorizonRC(t *testing.T) {
	// Horizon's generated openrc sets a user domain name and a project domain
	// ID; gophercloud refuses auth options that carry both an ID and a name.
	setEnv(t, map[string]string{
		"OS_AUTH_URL":             "https://keystone.example.com/v3",
		"OS_USERNAME":             "alice",
		"OS_PASSWORD":             "secret",
		"OS_PROJECT_ID":           "0123456789abcdef",
		"OS_PROJECT_NAME":         "research",
		"OS_USER_DOMAIN_NAME":     "Default",
		"OS_PROJECT_DOMAIN_ID":    "default",
		"OS_INTERFACE":            "public",
		"OS_IDENTITY_API_VERSION": "3",
	})

	opts := envAuthOptions()

	if opts.DomainName != "Default" || opts.DomainID != "" {
		t.Fatalf("DomainName = %q, DomainID = %q, want Default and empty", opts.DomainName, opts.DomainID)
	}

	if opts.TenantID != "0123456789abcdef" || opts.TenantName != "research" {
		t.Fatalf("TenantID = %q, TenantName = %q", opts.TenantID, opts.TenantName)
	}
}

func TestEnvAuthOptionsDomainIDWins(t *testing.T) {
	setEnv(t, map[string]string{
		"OS_AUTH_URL":         "https://keystone.example.com/v3",
		"OS_USER_DOMAIN_ID":   "d1",
		"OS_USER_DOMAIN_NAME": "Default",
		"OS_DOMAIN_NAME":      "ignored",
	})

	opts := envAuthOptions()

	if opts.DomainID != "d1" || opts.DomainName != "" {
		t.Fatalf("DomainID = %q, DomainName = %q, want d1 and empty", opts.DomainID, opts.DomainName)
	}
}

func TestEnvAuthOptionsApplicationCredential(t *testing.T) {
	setEnv(t, map[string]string{
		"OS_AUTH_URL":                      "https://keystone.example.com/v3",
		"OS_APPLICATION_CREDENTIAL_ID":     "app-id",
		"OS_APPLICATION_CREDENTIAL_SECRET": "app-secret",
	})

	opts := envAuthOptions()

	if opts.ApplicationCredentialID != "app-id" || opts.ApplicationCredentialSecret != "app-secret" {
		t.Fatalf("ApplicationCredentialID = %q, ApplicationCredentialSecret = %q", opts.ApplicationCredentialID, opts.ApplicationCredentialSecret)
	}

	if opts.Username != "" || opts.Password != "" || opts.DomainName != "" || opts.DomainID != "" {
		t.Fatalf("unexpected user fields set: %+v", opts)
	}

	if !opts.AllowReauth {
		t.Fatal("AllowReauth = false, want true")
	}
}

func TestEnvAuthOptionsUserID(t *testing.T) {
	setEnv(t, map[string]string{
		"OS_AUTH_URL": "https://keystone.example.com/v3",
		"OS_USER_ID":  "new",
		"OS_USERID":   "old",
	})

	if got := envAuthOptions().UserID; got != "new" {
		t.Fatalf("UserID = %q, want new", got)
	}

	t.Setenv("OS_USER_ID", "")

	if got := envAuthOptions().UserID; got != "old" {
		t.Fatalf("UserID = %q, want old", got)
	}
}

func TestEnvAuthOptionsToken(t *testing.T) {
	setEnv(t, map[string]string{
		"OS_AUTH_URL": "https://keystone.example.com/v3",
		"OS_TOKEN":    "token",
	})

	opts := envAuthOptions()

	if opts.TokenID != "token" {
		t.Fatalf("TokenID = %q, want token", opts.TokenID)
	}

	if opts.AllowReauth {
		t.Fatal("AllowReauth = true, want false with a token")
	}
}

func TestConnectFromEnvError(t *testing.T) {
	setEnv(t, map[string]string{
		"OS_AUTH_URL":         "http://127.0.0.1:1/v3",
		"OS_USERNAME":         "alice",
		"OS_PASSWORD":         "secret",
		"OS_PROJECT_NAME":     "research",
		"OS_USER_DOMAIN_NAME": "Default",
	})

	c := Context{Cloud: EnvCloudName}

	err := c.ConnectFromEnv(context.Background())
	if err == nil {
		t.Fatal("ConnectFromEnv() succeeded against a closed port")
	}

	if !strings.Contains(err.Error(), `authenticating cloud "envvars"`) {
		t.Fatalf("ConnectFromEnv() error = %q", err)
	}

	if c.Provider != nil {
		t.Fatal("Provider set after a failed connection")
	}
}
