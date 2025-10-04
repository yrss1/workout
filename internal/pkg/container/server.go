package container

import (
	"context"

	"github.com/yrss1/workout/internal/pkg/http"
	"golang.org/x/sync/errgroup"
)

type Container struct {
	app        *App
	httpServer *http.Server
}

func NewContainer(configPath string) (*Container, error) {
	app, err := NewApp(configPath)
	if err != nil {
		return nil, err
	}
	recoverer := http.NewPanicRecoverer(app.logger)

	apiRouter := http.NewRouter(app.logger, &app.config.HTTP, recoverer)

	httpServer := http.NewServer(app.config.HTTP, app.logger, apiRouter.Handler())

	return &Container{
		app:        app,
		httpServer: httpServer,
	}, nil
}

func (c *Container) Start(ctx context.Context) error {
	errGroup := errgroup.Group{}
	errGroup.Go(func() error {
		if err := c.httpServer.Start(ctx); err != nil {
			return err
		}
		return nil
	})
	if err := errGroup.Wait(); err != nil {
		return err
	}
	return nil
}

func (c *Container) Stop(ctx context.Context) error {
	if err := c.httpServer.Stop(ctx); err != nil {
		return err
	}
	return nil
}
