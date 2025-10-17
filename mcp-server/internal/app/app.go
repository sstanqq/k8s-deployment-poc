package app

import (
	"context"
	"fmt"

	"github.com/sstanqq/k8s-deployment-poc/mcp-server/internal/config"
	"github.com/sstanqq/k8s-deployment-poc/mcp-server/internal/fs"
	"github.com/sstanqq/k8s-deployment-poc/mcp-server/internal/host"
	mcpserver "github.com/sstanqq/k8s-deployment-poc/mcp-server/internal/mcp"
)

type Application struct {
	conf *config.Config

	srv   *mcpserver.Server
	store *fs.RequestStore
}

func NewApplication(conf *config.Config) (*Application, error) {
	store, err := fs.NewRequestStore(conf.LogFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create request store: %w", err)
	}

	h := host.NewSystemHost(conf.HstConfig.NodeName, conf.HstConfig.NodeIP)

	srv := mcpserver.NewServer(&mcpserver.ServerOpts{
		Addr:    fmt.Sprintf("%s:%d", conf.SrvHost, conf.SrvPort),
		Name:    conf.SrvName,
		Version: conf.SrvVersion,
		Store:   store,
		Host:    h,
	})

	return &Application{
		conf:  conf,
		srv:   srv,
		store: store,
	}, nil
}

func (a *Application) Run(ctx context.Context) error {
	return a.srv.Run(ctx)
}

func (a *Application) Stop(ctx context.Context) error {
	if err := a.srv.Stop(ctx); err != nil {
		return err
	}

	if err := a.store.Close(); err != nil {
		return err
	}

	return nil
}
