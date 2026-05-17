package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"
)

type Greeter interface{ Greet() string }

type englishGreeter struct{}

func (englishGreeter) Greet() string     { return "05: fx.Annotate + fx.As" }
func newEnglishGreeter() *englishGreeter { return &englishGreeter{} }
func run(g Greeter)                      { fmt.Println(g.Greet()) }

func main() {
	app := fx.New(
		fx.Provide(fx.Annotate(newEnglishGreeter, fx.As(new(Greeter)))),
		fx.Invoke(run),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = app.Start(ctx)
	_ = app.Stop(ctx)
}
