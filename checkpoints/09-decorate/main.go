package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/fx"
)

func provideMessage() string          { return "09: decorate me" }
func decorateMessage(v string) string { return strings.ToUpper(v) }
func run(v string)                    { fmt.Println(v) }

func main() {
	app := fx.New(
		fx.Provide(provideMessage),
		fx.Decorate(decorateMessage),
		fx.Invoke(run),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = app.Start(ctx)
	_ = app.Stop(ctx)
}
