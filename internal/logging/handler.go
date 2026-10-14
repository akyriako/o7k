package logging

import (
	"context"
	"log/slog"
	"strings"

	"github.com/hashicorp/go-hclog"
)

type slogHandler struct {
	hl     hclog.Logger
	groups []string
}

func newSlogHandler(hl hclog.Logger) slog.Handler {
	return &slogHandler{hl: hl}
}

func toHclogLevel(l slog.Level) hclog.Level {
	switch {
	case l < slog.LevelDebug:
		return hclog.Trace
	case l < slog.LevelInfo:
		return hclog.Debug
	case l < slog.LevelWarn:
		return hclog.Info
	case l < slog.LevelError:
		return hclog.Warn
	default:
		return hclog.Error
	}
}

func (h *slogHandler) Enabled(_ context.Context, l slog.Level) bool {
	switch toHclogLevel(l) {
	case hclog.Trace:
		return h.hl.IsTrace()
	case hclog.Debug:
		return h.hl.IsDebug()
	case hclog.Info:
		return h.hl.IsInfo()
	case hclog.Warn:
		return h.hl.IsWarn()
	default:
		return h.hl.IsError()
	}
}

func (h *slogHandler) Handle(_ context.Context, r slog.Record) error {
	args := make([]any, 0, r.NumAttrs()*2)
	r.Attrs(func(a slog.Attr) bool {
		args = h.appendAttr(args, a)
		return true
	})
	h.hl.Log(toHclogLevel(r.Level), r.Message, args...)
	return nil
}

func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	args := make([]any, 0, len(attrs)*2)
	for _, a := range attrs {
		args = h.appendAttr(args, a)
	}
	return &slogHandler{hl: h.hl.With(args...), groups: h.groups}
}

func (h *slogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &slogHandler{
		hl:     h.hl,
		groups: append(append([]string{}, h.groups...), name),
	}
}

func (h *slogHandler) appendAttr(args []any, a slog.Attr) []any {
	a.Value = a.Value.Resolve()
	if a.Equal(slog.Attr{}) {
		return args
	}

	if a.Value.Kind() == slog.KindGroup {
		sub := &slogHandler{hl: h.hl, groups: h.groups}
		if a.Key != "" {
			sub.groups = append(append([]string{}, h.groups...), a.Key)
		}
		for _, ga := range a.Value.Group() {
			args = sub.appendAttr(args, ga)
		}
		return args
	}

	key := a.Key
	if len(h.groups) > 0 {
		key = strings.Join(h.groups, ".") + "." + key
	}
	return append(args, key, a.Value.Any())
}
