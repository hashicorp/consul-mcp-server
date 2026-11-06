// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package main

import (
	"context"
	"log/slog"
	"maps"

	log "github.com/sirupsen/logrus"
)

// logrusHandler is a slog.Handler that forwards records to logrus at the matching level.
type logrusHandler struct {
	logger *log.Logger
	fields log.Fields
	prefix string
}

func newSlogLogger(logger *log.Logger) *slog.Logger {
	return slog.New(&logrusHandler{logger: logger})
}

func (h *logrusHandler) Enabled(_ context.Context, level slog.Level) bool {
	return h.logger.IsLevelEnabled(logrusLevel(level))
}

func (h *logrusHandler) Handle(_ context.Context, r slog.Record) error {
	fields := make(log.Fields, len(h.fields)+r.NumAttrs())
	maps.Copy(fields, h.fields)
	r.Attrs(func(a slog.Attr) bool {
		addAttr(fields, h.prefix, a)
		return true
	})
	h.logger.WithFields(fields).Log(logrusLevel(r.Level), r.Message)
	return nil
}

func (h *logrusHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	fields := make(log.Fields, len(h.fields)+len(attrs))
	maps.Copy(fields, h.fields)
	for _, a := range attrs {
		addAttr(fields, h.prefix, a)
	}
	return &logrusHandler{logger: h.logger, fields: fields, prefix: h.prefix}
}

func (h *logrusHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &logrusHandler{logger: h.logger, fields: h.fields, prefix: h.prefix + name + "."}
}

func addAttr(fields log.Fields, prefix string, a slog.Attr) {
	v := a.Value.Resolve()
	if v.Kind() == slog.KindGroup {
		if a.Key != "" {
			prefix += a.Key + "."
		}
		for _, ga := range v.Group() {
			addAttr(fields, prefix, ga)
		}
		return
	}
	if a.Key == "" {
		return
	}
	fields[prefix+a.Key] = v.Any()
}

func logrusLevel(level slog.Level) log.Level {
	switch {
	case level >= slog.LevelError:
		return log.ErrorLevel
	case level >= slog.LevelWarn:
		return log.WarnLevel
	case level >= slog.LevelInfo:
		return log.InfoLevel
	default:
		return log.DebugLevel
	}
}
