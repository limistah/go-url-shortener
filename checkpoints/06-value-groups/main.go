package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"
)

type out struct {
	fx.Out
	Route string `group:"routes"`
}

func routeA() out { return out{Route: "POST /shorten"} }
func routeB() out { return out{Route: "GET /{slug}"} }

func printRoutes(in struct {
	fx.In
	Routes []string `group:"routes"`
}) {
	fmt.Printf("06: %v\n", in.Routes)
}

func main() {
	app := fx.New(
		fx.Provide(routeA, routeB),
		fx.Invoke(printRoutes),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = app.Start(ctx)
	_ = app.Stop(ctx)
}
