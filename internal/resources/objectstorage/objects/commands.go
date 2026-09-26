package objects

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/akyriako/o7k/internal/resource"
	tea "github.com/charmbracelet/bubbletea"
	swiftobjects "github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/objects"
)

func (r *Resource) show(row resource.Row) tea.Cmd {
	container := row.Fields["container"]

	return func() tea.Msg {
		client, err := r.context.ObjectStorageV1()
		if err != nil {
			return resource.DetailsMsg{
				ID:  row.ID,
				Err: err,
			}
		}

		item, err := swiftobjects.Get(context.Background(), client, container, row.ID, nil).Extract()
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

const previewSize int64 = 256 * 1024

func (r *Resource) view(row resource.Row) tea.Cmd {
	container := row.Fields["container"]

	return func() tea.Msg {
		client, err := r.context.ObjectStorageV1()
		if err != nil {
			return resource.DetailsMsg{
				ID:  row.ID,
				Err: err,
			}
		}

		size, err := strconv.ParseInt(row.Fields["bytes"], 10, 64)
		if err != nil {
			return resource.DetailsMsg{
				ID:  row.ID,
				Err: err,
			}
		}

		opts := swiftobjects.DownloadOpts{}

		if size > previewSize {
			opts.Range = fmt.Sprintf("bytes=0-%d", previewSize-1)
		}

		result := swiftobjects.Download(context.Background(), client, container, row.ID, opts)

		if result.Err != nil {
			return resource.DetailsMsg{
				ID:  row.ID,
				Err: result.Err,
			}
		}

		content, err := result.ExtractContent()
		if err != nil {
			return resource.DetailsMsg{
				ID:  row.ID,
				Err: err,
			}
		}

		value := string(content)
		encoding := "text"

		if !isTextContent(content) {
			value = base64.StdEncoding.EncodeToString(content)
			encoding = "base64"
		}

		return resource.DetailsMsg{
			ID: row.ID,
			Content: map[string]any{
				"size":            size,
				"displayed_bytes": len(content),
				"truncated":       size > int64(len(content)),
				"encoding":        encoding,
				"content":         value,
			},
		}
	}
}

func isTextContent(content []byte) bool {
	if !utf8.Valid(content) {
		return false
	}

	return strings.IndexFunc(string(content), func(r rune) bool {
		return !unicode.IsPrint(r) && !unicode.IsSpace(r)
	}) == -1
}
