package main

import (
	"fmt"
	"os"
	"time"

	"github.com/tryor/trygo"
	"github.com/tryor/trygo-bridge/graceful"
)

func main() {
	ListenAndServe(":8080", &graceful.GracefulServer{Timeout: 10 * time.Second})
	//ListenAndServe(":4333", &graceful.TLSGracefulServer{Timeout: 10 * time.Second, CertFile: "cert.pem", KeyFile: "key.pem"})
}

func ListenAndServe(addr string, server trygo.Server) {
	app := trygo.NewApp()
	app.Config.Listen.Addr = addr
	app.Get("/", func(ctx *trygo.Context) {
		d, err := time.ParseDuration(ctx.Input.GetValue("d"))
		if err != nil {
			ctx.ResponseWriter.Error(err.Error(), 400)
			return
		}
		w := ctx.ResponseWriter
		fmt.Fprintf(w, "going to sleep %s with pid %d\n", d, os.Getpid())
		w.Flush()
		time.Sleep(d)
		fmt.Fprintf(w, "slept %s with pid %d\n", d, os.Getpid())
	})

	fmt.Println("ListenAndServe AT ", app.Config.Listen.Addr)
	if err := server.ListenAndServe(app); err != nil {
		app.Logger.Critical("ListenAndServe: %v", err)
	}

}
