package app

import (
	"context"
	"redis-crud-playground/api"
	"redis-crud-playground/api/middleware"
	"redis-crud-playground/internals/app/db"
	"redis-crud-playground/internals/app/handlers"
	"redis-crud-playground/internals/app/services"
	"redis-crud-playground/internals/cfg"

	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

type App struct {
	rdb    *redis.Client
	config cfg.Cfg
}

func NewApp(cfg cfg.Cfg) *App {
	app := &App{rdb: redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddress,
	}),
		config: cfg,
	}
	return app
}

func (a *App) Start(ctx context.Context) error {

	ordersStorage := db.NewOrdersStorage(a.rdb)

	ordersService := services.NewOrdersService(ordersStorage)

	ordersHandler := handlers.NewOrdersHandler(ordersService)

	routes := api.CreateRoutes(ordersHandler)

	routes.Use(middleware.RequestLog)

	server := &http.Server{Addr: fmt.Sprintf(":%d", a.config.ServerPort), Handler: routes}
	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect reddis: %w", err)
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis")
		}
	}()

	fmt.Println("Starting server")

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}

	return nil
}
