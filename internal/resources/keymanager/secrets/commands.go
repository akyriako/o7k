package secrets

import (
	"context"
	"encoding/base64"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	barbicansecrets "github.com/gophercloud/gophercloud/v2/openstack/keymanager/v1/secrets"
)

func (r *Resource) show(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.KeyManagerV1()
		if err != nil {
			return resource.DetailsMsg{
				ID:  row.ID,
				Err: err,
			}
		}

		item, err := barbicansecrets.Get(
			context.Background(),
			client,
			row.ID,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				ID:  row.ID,
				Err: err,
			}
		}

		return resource.DetailsMsg{
			ID:      row.ID,
			Content: item,
		}
	}
}

func (r *Resource) payload(row resource.Row) tea.Cmd {
	return func() tea.Msg {
		client, err := r.context.KeyManagerV1()
		if err != nil {
			return resource.DetailsMsg{
				ID:  row.ID,
				Err: err,
			}
		}

		secret, err := barbicansecrets.Get(
			context.Background(),
			client,
			row.ID,
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				ID:  row.ID,
				Err: err,
			}
		}

		contentType := secret.ContentTypes["default"]

		payload, err := barbicansecrets.GetPayload(
			context.Background(),
			client,
			row.ID,
			barbicansecrets.GetPayloadOpts{
				PayloadContentType: contentType,
			},
		).Extract()
		if err != nil {
			return resource.DetailsMsg{
				ID:  row.ID,
				Err: err,
			}
		}

		value := string(payload)
		encoding := "text"

		if !isTextPayload(payload) {
			value = base64.StdEncoding.EncodeToString(payload)
			encoding = "base64"
		}

		return resource.DetailsMsg{
			ID: row.ID,
			Content: map[string]any{
				"content_type": contentType,
				"encoding":     encoding,
				"payload":      value,
			},
		}
	}
}

func isTextPayload(payload []byte) bool {
	if !utf8.Valid(payload) {
		return false
	}

	return strings.IndexFunc(string(payload), func(r rune) bool {
		return !unicode.IsPrint(r) && !unicode.IsSpace(r)
	}) == -1
}
