package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"
)

func registerLifecycle(lc fx.Lifecycle) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			fmt.Println("03: OnStart hook")
			return nil
		},
		OnStop: func(context.Context) error {
			fmt.Println("03: OnStop hook")
			return nil
		},
	})
}

func main() {
	app := fx.New(fx.Invoke(registerLifecycle))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = app.Start(ctx)
	_ = app.Stop(ctx)
}
