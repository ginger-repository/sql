package sql

import (
	"time"
)

type config struct {
	Debug bool

	Dialect          string
	ConnectionString string

	ConnMaxLifetime    *time.Duration
	MaxIdleConnections *int
	MaxOpenConnections *int

	Logger loggerConfig
}

func (c *config) initialize() {
}
