// Package app 负责依赖装配，将各层组合为可运行的应用。
package app

import (
	"net/http"

	"inventory/internal/config"
	"inventory/internal/handler"
	"inventory/internal/service"
	"inventory/internal/store"
	"inventory/pkg/logger"
)

// App 应用容器，持有装配完成的处理器。
type App struct {
	server *handler.Server
}

// New 装配全部依赖：store -> service -> handler。
func New(cfg *config.Config, log *logger.Logger) (*App, error) {
	st := store.NewMemoryStore()
	svc := service.New(st, log, cfg)
	server := handler.NewServer(svc, log, cfg)

	log.Infof("应用装配完成，配置：%s", cfg.String())
	return &App{server: server}, nil
}

// Routes 返回应用的 HTTP 处理器。
func (a *App) Routes() http.Handler {
	return a.server.Routes()
}
