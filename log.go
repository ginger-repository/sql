package sql

import (
	"context"
	"fmt"
	l "log"
	"strings"
	"time"

	"github.com/ginger-core/log"
	gormLogger "gorm.io/gorm/logger"
)

type LogLevel string

const (
	// Silent silent log level
	Silent LogLevel = "silent"
	// Error error log level
	Error LogLevel = "error"
	// Warn warn log level
	Warn LogLevel = "warn"
	// Info info log level
	Info LogLevel = "info"
)

func getLevel(lvl LogLevel) gormLogger.LogLevel {
	switch LogLevel(strings.ToLower(string(lvl))) {
	case Silent:
		return gormLogger.Silent
	case Error:
		return gormLogger.Error
	case Warn:
		return gormLogger.Warn
	case Info:
		return gormLogger.Info
	}
	panic(fmt.Sprintf("level `%s` not found", lvl))
}

type loggerConfig struct {
	Enabled                   bool
	SlowThreshold             time.Duration
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	LogLevel                  LogLevel
}

type _logger struct {
	interfaces []gormLogger.Interface
	source     log.Logger
}

func (repo *repo) newLogger(source log.Logger) gormLogger.Interface {
	if !repo.config.Logger.Enabled {
		return nil
	}
	writers := source.GetWriters()

	r := &_logger{
		interfaces: make([]gormLogger.Interface, len(writers)),
		source:     source,
	}

	for i, w := range writers {
		v := gormLogger.New(
			l.New(w, "\r\n", l.LstdFlags), // io writer
			gormLogger.Config{
				SlowThreshold:             5 * time.Second,
				LogLevel:                  gormLogger.Warn,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			},
		)
		r.interfaces[i] = v
	}
	return r
}

func (l *_logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	panic("implement me")
}

func (l *_logger) Info(ctx context.Context, s string, i ...interface{}) {
	panic("implement me")
}

func (l *_logger) Warn(ctx context.Context, s string, i ...interface{}) {
	panic("implement me")
}

func (l *_logger) Error(ctx context.Context, s string, i ...interface{}) {
	panic("implement me")
}

func (l *_logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	panic("implement me")
}
