package mcp

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Server struct {
	srv *http.Server
}

type ServerOpts struct {
	Addr    string
	Name    string
	Version string

	Store RequestStore
	Host  HostProvider
}

func NewServer(srvOpts *ServerOpts) *Server {
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    srvOpts.Name,
		Version: srvOpts.Version,
	}, nil)

	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return mcpServer
	}, nil)

	registerTools(mcpServer, srvOpts.Store, srvOpts.Host)

	httpServer := &http.Server{
		Addr:    srvOpts.Addr,
		Handler: handler,
	}

	return &Server{
		srv: httpServer,
	}
}

func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		log.Printf("mcp server listening on %s", s.srv.Addr)
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

func (s *Server) Stop(ctx context.Context) error {
	done := make(chan error, 1)

	go func() {
		done <- s.srv.Shutdown(ctx)
	}()

	select {
	case <-ctx.Done():
		_ = s.srv.Close()
		return ctx.Err()
	case err := <-done:
		return err
	}
}
