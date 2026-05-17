package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"
)

func providePrimary() string   { return "primary-store" }
func provideSecondary() string { return "secondary-store" }

func useNamed(in struct {
	fx.In
	Primary   string `name:"primary"`
	Secondary string `name:"secondary"`
}) {
	fmt.Printf("07: %s + %s\n", in.Primary, in.Secondary)
}

func main() {
	app := fx.New(
		fx.Provide(
			fx.Annotate(providePrimary, fx.ResultTags(`name:"primary"`)),
			fx.Annotate(provideSecondary, fx.ResultTags(`name:"secondary"`)),
		),
		fx.Invoke(useNamed),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = app.Start(ctx)
	_ = app.Stop(ctx)
}
