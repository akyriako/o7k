package resource

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

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

type Command struct {
	Key         string
	Description string
	Default     bool
}

type Resource interface {
	Kind() string
	Title() string
	Aliases() []string
	Columns() []Column
	Commands() []Command
	List(ctx context.Context) ([]Row, error)
	Execute(command Command, row Row) tea.Cmd
}

type scopeContextKey struct{}

func WithScope(ctx context.Context, scope map[string]string) context.Context {
	return context.WithValue(ctx, scopeContextKey{}, scope)
}

func Scope(ctx context.Context) map[string]string {
	scope, _ := ctx.Value(scopeContextKey{}).(map[string]string)
	return scope
}
