package plugins

import "github.com/devimteam/microgen/pkg/microgen"

func init() {
	microgen.RegisterPlugin(loggingPlugin, &loggingMiddlewarePlugin{})
}

const (
	serviceAlias = "service"
	_service_    = "svc"
	_logger_     = "logger"
	_ctx_        = "ctx"
	_next_       = "next"
)
