package clusters

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/certificates"
	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/clusters"
	"gopkg.in/yaml.v3"
)

func (r *Resource) clusterTemplate(row resource.Row) tea.Cmd {
	clusterTemplateID := row.Fields["cluster_template_id"]

	return func() tea.Msg {
		return resource.NavigateFilteredMsg{
			Resource: "clustertemplates",
			Field:    "id",
			Value:    clusterTemplateID,
		}
	}
}

func (r *Resource) nodeGroups(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		return resource.NavigateScopedMsg{
			Resource: "nodegroups",
			Scope: map[string]string{
				"cluster_id": row.ID,
			},
		}
	}
}

func (r *Resource) show(id string) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.ContainerInfraV1()
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		cluster, err := clusters.Get(
			context.Background(),
			client,
			id,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting cluster %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      cluster.UUID,
			Content: cluster,
		}
	}
}
func (r *Resource) kubeconfig(id string) tea.Cmd {
	return func() tea.Msg {
		data, err := r.getKubeconfigData(context.Background(), id)
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		var content any
		if err := yaml.Unmarshal(data, &content); err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("decoding kubeconfig for cluster %q: %w", id, err),
			}
		}

		return resource.DetailsMsg{
			ID:      id,
			Content: content,
		}
	}
}

func (r *Resource) getKubeconfigData(ctx context.Context, id string) ([]byte, error) {
	client, err := r.context.ContainerInfraV1()
	if err != nil {
		return nil, fmt.Errorf("getting container infrastructure client: %w", err)
	}

	cluster, err := clusters.Get(ctx, client, id).Extract()
	if err != nil {
		return nil, fmt.Errorf("getting cluster %q: %w", id, err)
	}

	ca, err := certificates.Get(ctx, client, id).Extract()
	if err != nil {
		return nil, fmt.Errorf("getting CA certificate for cluster %q: %w", id, err)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generating private key: %w", err)
	}

	csrDER, err := x509.CreateCertificateRequest(
		rand.Reader,
		&x509.CertificateRequest{
			Subject: pkix.Name{
				CommonName:   "admin",
				Organization: []string{"system:masters"},
			},
		},
		privateKey,
	)
	if err != nil {
		return nil, fmt.Errorf("generating certificate request: %w", err)
	}

	csr := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrDER,
	})

	cert, err := certificates.Create(ctx, client, certificates.CreateOpts{
		ClusterUUID: id,
		CSR:         string(csr),
	}).Extract()
	if err != nil {
		return nil, fmt.Errorf("creating client certificate for cluster %q: %w", id, err)
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	clusterName := cluster.Name
	if clusterName == "" {
		clusterName = id
	}

	encode := func(data []byte) string {
		return base64.StdEncoding.EncodeToString(data)
	}

	data := fmt.Sprintf(`apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: %s
    server: %s
  name: %s
contexts:
- context:
    cluster: %s
    user: admin
  name: default
current-context: default
kind: Config
preferences: {}
users:
- name: admin
  user:
    client-certificate-data: %s
    client-key-data: %s
`,
		encode([]byte(strings.TrimSpace(ca.PEM))),
		cluster.APIAddress,
		clusterName,
		clusterName,
		encode([]byte(strings.TrimSpace(cert.PEM))),
		encode(privateKeyPEM),
	)

	return []byte(data), nil
}

func (r *Resource) getKubeconfig(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		data, err := r.getKubeconfigData(context.Background(), row.ID)
		if err != nil {
			return resource.DetailsMsg{Err: err}
		}

		home, err := os.UserHomeDir()
		if err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("getting user home directory: %w", err),
			}
		}

		kubeDir := filepath.Join(home, ".kube")
		if err := os.MkdirAll(kubeDir, 0700); err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("creating kubeconfig directory: %w", err),
			}
		}

		clusterName := row.Fields["name"]
		if clusterName == "" {
			clusterName = row.ID
		}

		filename := fmt.Sprintf("o7k-generated_%s_%s_magnum-%s.yaml", r.context.Region, r.context.Cloud, clusterName)
		path := filepath.Join(kubeDir, filename)

		if err := os.WriteFile(path, data, 0600); err != nil {
			return resource.DetailsMsg{
				Err: fmt.Errorf("writing kubeconfig %q: %w", path, err),
			}
		}

		return nil
	}
}
