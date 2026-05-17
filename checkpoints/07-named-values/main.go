package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"
)

type message struct{ text string }

func newMessage() message { return message{text: "checkpoint running"} }

func run(m message) { fmt.Println(m.text) }

func main() {
	app := fx.New(
		fx.Provide(newMessage),
		fx.Invoke(run),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = app.Start(ctx)
	_ = app.Stop(ctx)
}
