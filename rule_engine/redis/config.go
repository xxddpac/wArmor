package redis

import (
	"fmt"
	"runtime"
)

const (
	MaxDialTimeout  = 1000 // millisecond
	MaxReadTimeout  = 1000 // millisecond
	MaxWriteTimeout = 3000 // millisecond
	MaxPoolTimeout  = 2    // second
	MinIdleConns    = 3
	MaxRetries      = 1
)

// A Config of go redis
type Config struct {
	Network              string
	Addr                 string
	Passwd               string
	DB                   int
	DialTimeout          int
	ReadTimeout          int
	WriteTimeout         int
	PoolSize             int
	PoolTimeout          int
	MinIdleConns         int
	MaxRetries           int
	TraceIncludeNotFound bool
}

// Name returns client name of the config
func (c *Config) Name() string {
	return fmt.Sprintf("%s(%s/%d)", c.Network, c.Addr, c.DB)
}

// FillWithDefaults apply default values for fields with invalid values.
func (c *Config) FillWithDefaults() {
	if c.DialTimeout <= 0 {
		c.DialTimeout = MaxDialTimeout
	}

	if c.ReadTimeout <= 0 {
		c.ReadTimeout = MaxReadTimeout
	}

	if c.WriteTimeout <= 0 {
		c.WriteTimeout = MaxWriteTimeout
	}

	if c.PoolSize <= 0 {
		c.PoolSize = 10 * runtime.NumCPU()
	}

	if c.PoolTimeout <= 0 {
		c.PoolTimeout = MaxPoolTimeout
	}

	if c.MinIdleConns <= 0 {
		c.MinIdleConns = MinIdleConns
	}

	if c.MaxRetries < 0 {
		c.MaxRetries = MaxRetries
	}
}
