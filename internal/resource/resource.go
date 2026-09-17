package resource

import "context"

type Column struct {
	Key      string
	Title    string
	MinWidth int
	Flex     int
}

type Row struct {
	ID     string
	Fields map[string]string
}

type Resource interface {
	Kind() string
	Aliases() []string
	Title() string
	Columns() []Column
	Commands() []Command
	List(ctx context.Context) ([]Row, error)
}

type Command struct {
	Key         string
	Description string
}
