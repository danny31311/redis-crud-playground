package main

import (
	"context"
	"redis-crud-playground/internals/app"
	"redis-crud-playground/internals/cfg"

	"fmt"
	"os"
	"os/signal"
)

func main() {
	application := app.NewApp(cfg.LoadConfig())
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := application.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app: ", err)
	}
}
