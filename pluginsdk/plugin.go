package pluginsdk

import (
	"context"
	"encoding/json"
)

type Client interface {
	Metadata() (Metadata, error)
	Resources() ([]Resource, error)
}

type Plugin interface {
	Metadata() Metadata
	Resources() []Resource
	SetHost(Host)
}

type Metadata struct {
	Name    string
	Version string
	Color   string
}

type Context struct {
	Generation uint64
	Cloud      string
	CloudsPath string
	Region     string
}

type Host interface {
	Context(context.Context) (Context, error)
}

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
	Service() string
	Kind() string
	Title() string
	Aliases() []string
	Columns() []Column
	Commands() []Command
	List(context.Context) ([]Row, error)
	Execute(context.Context, Command, Row) (Result, error)
}

type Result struct {
	Details  *Details
	Navigate *Navigate
}

type Details struct {
	ID      string
	Content json.RawMessage
}

type Navigate struct {
	Resource string
	ID       string
	Field    string
	Value    string
	Scope    map[string]string
}

type scopeContextKey struct{}

func WithScope(ctx context.Context, scope map[string]string) context.Context {
	return context.WithValue(ctx, scopeContextKey{}, scope)
}

func Scope(ctx context.Context) map[string]string {
	scope, _ := ctx.Value(scopeContextKey{}).(map[string]string)
	return scope
}
