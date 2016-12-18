package main

import (
	"fmt"

	"github.com/tryor/trygo"
	"github.com/tryor/trygo-bridge/fasthttp"
)

func main() {
	app := trygo.NewApp()

	app.Get("/", func(ctx *trygo.Context) {

		fmt.Println(ctx.Request.Header)

		ctx.Render("hello world!")
	})

	fmt.Println("ListenAndServe AT ", app.Config.Listen.Addr)
	var server fasthttp.FasthttpServer
	if err := server.ListenAndServe(app); err != nil {
		app.Logger.Critical("ListenAndServe: %v", err)
	}
}
