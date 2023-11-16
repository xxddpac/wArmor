package global

import (
	"context"
)

const (
	MessageChannel    = "waf"
	DefaultTimeLayout = "2006-01-02 15:04:05"
	Heartbeat         = "heartbeat"
	Rule              = "rule"
	Config            = "config"
	Ip                = "ip"
)

var (
	Ctx    context.Context
	Cancel context.CancelFunc
)

func init() {
	Ctx, Cancel = context.WithCancel(context.Background())
}
