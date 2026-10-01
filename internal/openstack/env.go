package openstack

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/config"
)

const (
	// EnvCloudName is the name of the context built from the OS_* environment
	// variables. It matches the name openstacksdk uses for the same thing.
	EnvCloudName = "envvars"

	// EnvCloudSource is shown in place of a clouds.yaml path for the context
	// built from the OS_* environment variables.
	EnvCloudSource = "OS_* environment"
)

// EnvCloud reports whether OS_AUTH_URL is set and, if so, returns the display
// data of the context built from the OS_* environment variables.
func EnvCloud() (Cloud, bool) {
	authURL := os.Getenv("OS_AUTH_URL")
	if authURL == "" {
		return Cloud{}, false
	}

	domainID, domainName := envDomain()

	return Cloud{
		Name:     EnvCloudName,
		Region:   os.Getenv("OS_REGION_NAME"),
		Project:  firstNonEmpty(os.Getenv("OS_PROJECT_NAME"), os.Getenv("OS_PROJECT_ID")),
		Domain:   firstNonEmpty(domainName, domainID),
		Identity: authURL,
	}, true
}

// ConnectFromEnv authenticates using the OS_* environment variables.
func (c *Context) ConnectFromEnv(ctx context.Context) error {
	provider, err := config.NewProviderClient(ctx, envAuthOptions())
	if err != nil {
		return fmt.Errorf("authenticating cloud %q: %w", c.Cloud, err)
	}

	c.Provider = provider
	c.CloudsPath = ""

	slog.Info("connected to cloud", "cloud", c.Cloud, "identityEndpoint", provider.IdentityEndpoint)

	return nil
}

// envAuthOptions maps the OS_* environment variables onto gophercloud auth
// options, with the same precedence clouds.Parse applies to clouds.yaml.
func envAuthOptions() gophercloud.AuthOptions {
	domainID, domainName := envDomain()
	token := os.Getenv("OS_TOKEN")

	return gophercloud.AuthOptions{
		IdentityEndpoint:            os.Getenv("OS_AUTH_URL"),
		Username:                    os.Getenv("OS_USERNAME"),
		UserID:                      firstNonEmpty(os.Getenv("OS_USER_ID"), os.Getenv("OS_USERID")),
		Password:                    os.Getenv("OS_PASSWORD"),
		Passcode:                    os.Getenv("OS_PASSCODE"),
		DomainID:                    domainID,
		DomainName:                  domainName,
		TenantID:                    firstNonEmpty(os.Getenv("OS_PROJECT_ID"), os.Getenv("OS_TENANT_ID")),
		TenantName:                  firstNonEmpty(os.Getenv("OS_PROJECT_NAME"), os.Getenv("OS_TENANT_NAME")),
		TokenID:                     token,
		ApplicationCredentialID:     os.Getenv("OS_APPLICATION_CREDENTIAL_ID"),
		ApplicationCredentialName:   os.Getenv("OS_APPLICATION_CREDENTIAL_NAME"),
		ApplicationCredentialSecret: os.Getenv("OS_APPLICATION_CREDENTIAL_SECRET"),
		// gophercloud rejects AllowReauth when a token is passed through as is.
		AllowReauth: token == "",
	}
}

// envDomain picks the domain used for authentication: the user domain first,
// then the project domain, then the plain domain. Within a level the ID wins
// over the name, and only one of the two is ever returned, because gophercloud
// refuses auth options that carry both a domain ID and a domain name.
func envDomain() (id, name string) {
	levels := [][2]string{
		{"OS_USER_DOMAIN_ID", "OS_USER_DOMAIN_NAME"},
		{"OS_PROJECT_DOMAIN_ID", "OS_PROJECT_DOMAIN_NAME"},
		{"OS_DOMAIN_ID", "OS_DOMAIN_NAME"},
	}

	for _, level := range levels {
		if value := os.Getenv(level[0]); value != "" {
			return value, ""
		}

		if value := os.Getenv(level[1]); value != "" {
			return "", value
		}
	}

	return "", ""
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}
