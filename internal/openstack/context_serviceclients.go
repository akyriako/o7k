package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/apiversions"
)

func newSharedFileSystemV2(provider *gophercloud.ProviderClient, opts gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error) {
	client, err := openstack.NewSharedFileSystemV2(provider, opts)
	if err != nil {
		return nil, err
	}

	pages, err := apiversions.List(client).AllPages(context.Background())
	if err != nil {
		return nil, fmt.Errorf("listing shared file system API versions: %w", err)
	}

	versions, err := apiversions.ExtractAPIVersions(pages)
	if err != nil {
		return nil, fmt.Errorf("extracting shared file system API versions: %w", err)
	}

	for _, version := range versions {
		if version.ID != "v2.0" {
			continue
		}

		client.Microversion = min(version.Version, "2.60")
		return client, nil
	}

	return nil, fmt.Errorf("shared file system API v2 not found")
}
