package openstack

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Cloud struct {
	Name     string
	Region   string
	Project  string
	Domain   string
	Identity string
}

type Clouds struct {
	Path  string
	Items []Cloud
}

type cloudsFile struct {
	Clouds map[string]cloudConfig `yaml:"clouds"`
}

type cloudConfig struct {
	Region string    `yaml:"region_name"`
	Auth   cloudAuth `yaml:"auth"`
}

type cloudAuth struct {
	AuthURL string `yaml:"auth_url"`
	Project string `yaml:"project_name"`
	Domain  string `yaml:"domain_name"`
}

func LoadClouds(path string) (*Clouds, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading clouds file %q: %w", path, err)
	}

	var file cloudsFile

	if err := yaml.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("parsing clouds file %q: %w", path, err)
	}

	if len(file.Clouds) == 0 {
		return nil, fmt.Errorf("no clouds found in %q", path)
	}

	var root yaml.Node

	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("parsing clouds file order %q: %w", path, err)
	}

	names := cloudNames(root)
	items := make([]Cloud, 0, len(names))

	for _, name := range names {
		config := file.Clouds[name]

		region := config.Region
		//if region == "" {
		//	region = "*"
		//}

		items = append(items, Cloud{
			Name:     name,
			Region:   region,
			Project:  config.Auth.Project,
			Domain:   config.Auth.Domain,
			Identity: config.Auth.AuthURL,
		})
	}

	return &Clouds{
		Path:  path,
		Items: items,
	}, nil
}

func cloudNames(root yaml.Node) []string {
	if len(root.Content) == 0 {
		return nil
	}

	document := root.Content[0]

	for i := 0; i+1 < len(document.Content); i += 2 {
		if document.Content[i].Value != "clouds" {
			continue
		}

		clouds := document.Content[i+1]
		names := make([]string, 0, len(clouds.Content)/2)

		for j := 0; j+1 < len(clouds.Content); j += 2 {
			names = append(names, clouds.Content[j].Value)
		}

		return names
	}

	return nil
}
