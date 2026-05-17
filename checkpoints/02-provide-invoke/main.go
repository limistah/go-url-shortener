package main

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/fx"
)

type message struct{ text string }

func newMessage() message    { return message{text: "02: fx.Provide + fx.Invoke"} }
func printMessage(m message) { fmt.Println(m.text) }

func main() {
	app := fx.New(
		fx.Provide(newMessage),
		fx.Invoke(printMessage),
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = app.Start(ctx)
	_ = app.Stop(ctx)
}
