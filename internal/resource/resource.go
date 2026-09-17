package resource

import "context"

type Column struct {
	Key   string
	Title string
	Width int
}

type Row struct {
	ID     string
	Fields map[string]string
}

type Resource interface {
	Kind() string
	Aliases() []string
	Columns() []Column

	List(ctx context.Context) ([]Row, error)
}
