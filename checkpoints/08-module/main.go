package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"
)

func configValue() string { return "08: module-composed app" }
func run(v string)        { fmt.Println(v) }

var configModule = fx.Module("config", fx.Provide(configValue))
var appModule = fx.Module("app", configModule, fx.Invoke(run))

func main() {
	app := fx.New(appModule)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = app.Start(ctx)
	_ = app.Stop(ctx)
}
