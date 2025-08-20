package log

import (
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"
)

type LogHandle func(time time.Time, msg, level string, source *slog.Source, attrs []slog.Attr)

type LogOptions struct {
	Level  string
	Handle LogHandle
}

func Init(options LogOptions) {
	if options.Handle != nil {
		opts := PrettyHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: options.parseLevel(),
			},
			handle: options.Handle,
		}
		slog.SetDefault(slog.New(NewPrettyHandler(os.Stdout, opts)))
	} else {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: options.parseLevel()})))
	}
}

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
	handle   LogHandle
}

type PrettyHandler struct {
	opts   PrettyHandlerOptions
	l      *log.Logger
	attrs  []slog.Attr
	groups []string
	handle LogHandle
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	if h.handle == nil {
		return nil
	}
	var allAttrs []slog.Attr
	allAttrs = append(allAttrs, h.attrs...)
	r.Attrs(func(a slog.Attr) bool {
		allAttrs = append(allAttrs, a)
		return true
	})
	h.handle(r.Time, r.Message, r.Level.String(), source(&r), allAttrs)
	return nil
}

func NewPrettyHandler(
	out io.Writer,
	opts PrettyHandlerOptions,
) *PrettyHandler {
	h := &PrettyHandler{
		opts:   opts,
		l:      log.New(out, "", 0),
		attrs:  []slog.Attr{},
		groups: []string{},
		handle: opts.handle,
	}
	return h
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.SlogOpts.Level.Level()
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &PrettyHandler{
		opts:   h.opts,
		l:      h.l,
		attrs:  newAttrs,
		groups: h.groups,
		handle: h.handle,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = name

	return &PrettyHandler{
		opts:   h.opts,
		l:      h.l,
		attrs:  h.attrs,
		groups: newGroups,
		handle: h.handle,
	}
}

func (o *LogOptions) parseLevel() slog.Level {
	switch strings.ToUpper(o.Level) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}


func source(r *slog.Record) *slog.Source {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()
	return &slog.Source{
		Function: f.Function,
		File:     f.File,
		Line:     f.Line,
	}
}