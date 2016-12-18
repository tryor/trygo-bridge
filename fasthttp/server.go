package fasthttp

import (
	"fmt"
	"net"
	"os"

	"github.com/tryor/trygo"
	fhttp "github.com/valyala/fasthttp"
)

type FasthttpServer struct {
	fhttp.Server
}

func (hsl *FasthttpServer) ListenAndServe(app *trygo.App) error {
	app.Prepare()
	hsl.Handler = NewFastHTTPHandler(app.FilterHandler(app, app.Handlers))
	configServer(&hsl.Server, app)

	ln, err := net.Listen("tcp4", app.Config.Listen.Addr)
	if err != nil {
		return err
	}
	return hsl.Server.Serve(app.FilterListener(app, ln))
	//return hsl.Server.ListenAndServe(app.Config.Listen.Addr)
}

type TLSFasthttpServer struct {
	fhttp.Server
	CertFile, KeyFile string
}

func (hsl *TLSFasthttpServer) ListenAndServe(app *trygo.App) error {
	app.Prepare()
	hsl.Handler = NewFastHTTPHandler(app.FilterHandler(app, app.Handlers))
	configServer(&hsl.Server, app)

	ln, err := net.Listen("tcp4", app.Config.Listen.Addr)
	if err != nil {
		return err
	}
	return hsl.Server.ServeTLS(app.FilterListener(app, ln), hsl.CertFile, hsl.KeyFile)

	//return hsl.Server.ListenAndServeTLS(app.Config.Listen.Addr, hsl.CertFile, hsl.KeyFile)
}

type UNIXFasthttpServer struct {
	fhttp.Server
	Mode os.FileMode
}

func (hsl *UNIXFasthttpServer) ListenAndServe(app *trygo.App) error {
	app.Prepare()
	hsl.Handler = NewFastHTTPHandler(app.FilterHandler(app, app.Handlers))
	configServer(&hsl.Server, app)

	addr := app.Config.Listen.Addr

	if err := os.Remove(addr); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unexpected error when trying to remove unix socket file %q: %s", addr, err)
	}
	ln, err := net.Listen("unix", addr)
	if err != nil {
		return err
	}
	if err = os.Chmod(addr, hsl.Mode); err != nil {
		return fmt.Errorf("cannot chmod %#o for %q: %s", hsl.Mode, addr, err)
	}
	return hsl.Server.Serve(app.FilterListener(app, ln))
	//return hsl.Server.ListenAndServeUNIX(app.Config.Listen.Addr, hsl.Mode)
}

func configServer(server *fhttp.Server, app *trygo.App) {
	server.ReadTimeout = app.Config.Listen.ReadTimeout
	server.WriteTimeout = app.Config.Listen.WriteTimeout
	server.Concurrency = app.Config.Listen.Concurrency
	server.MaxRequestBodySize = int(app.Config.MaxRequestBodySize)
	server.Logger = getLogger(app)
}

func getLogger(app *trygo.App) fhttp.Logger {
	if log, ok := app.Logger.(fhttp.Logger); ok {
		return log
	} else {
		return &logger{app.Logger}
	}
}

type logger struct {
	trygo.LoggerInterface
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.LoggerInterface.Info(format, args...)
}
