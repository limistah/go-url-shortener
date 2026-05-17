package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"
)

type in struct {
	fx.In
	Prefix string `name:"prefix"`
	Value  string `name:"value"`
}

type out struct {
	fx.Out
	Combined string
}

func providePrefix() string { return "04:" }
func provideValue() string  { return " fx.In/fx.Out" }

func combine(v in) out {
	return out{Combined: v.Prefix + v.Value}
}

func printCombined(v string) { fmt.Println(v) }

func main() {
	app := fx.New(
		fx.Provide(
			fx.Annotate(providePrefix, fx.ResultTags(`name:"prefix"`)),
			fx.Annotate(provideValue, fx.ResultTags(`name:"value"`)),
			combine,
		),
		fx.Invoke(printCombined),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = app.Start(ctx)
	_ = app.Stop(ctx)
}
