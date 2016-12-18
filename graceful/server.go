package graceful

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/tryor/trygo"
	g "github.com/tylerb/graceful"
)

type GracefulServer struct {
	g.Server
	Network string
	Timeout time.Duration
}

func (s *GracefulServer) ListenAndServe(app *trygo.App) error {
	app.Prepare()
	if s.Server.Server == nil {
		s.Server.Server = &http.Server{}
	}
	s.Server.Server.ReadTimeout = app.Config.Listen.ReadTimeout
	s.Server.Server.WriteTimeout = app.Config.Listen.WriteTimeout
	s.Server.Server.Addr = app.Config.Listen.Addr
	s.Server.Server.Handler = app.FilterHandler(app, trygo.DefaultFilterHandler(app, app.Handlers))
	s.Server.Timeout = s.Timeout

	if w, ok := app.Logger.(io.Writer); ok {
		s.ErrorLog = log.New(w, "[HTTP]", 0)
		s.Logger = s.ErrorLog
	}
	if s.Network == "" {
		s.Network = "tcp"
	}
	l, err := net.Listen(s.Network, s.Addr)
	if err != nil {
		return err
	}
	return s.Serve(app.FilterListener(app, trygo.DefaultFilterListener(app, l)))
}

//TLS
type TLSGracefulServer struct {
	g.Server
	CertFile, KeyFile string
	Timeout           time.Duration
}

func (s *TLSGracefulServer) ListenAndServe(app *trygo.App) error {
	app.Prepare()
	if s.Server.Server == nil {
		s.Server.Server = &http.Server{}
	}
	s.Server.Server.ReadTimeout = app.Config.Listen.ReadTimeout
	s.Server.Server.WriteTimeout = app.Config.Listen.WriteTimeout
	s.Server.Server.Addr = app.Config.Listen.Addr
	s.Server.Server.Handler = app.FilterHandler(app, trygo.DefaultFilterHandler(app, app.Handlers))
	s.Server.Timeout = s.Timeout

	if w, ok := app.Logger.(io.Writer); ok {
		s.ErrorLog = log.New(w, "[HTTPS]", 0)
		s.Logger = s.ErrorLog
	}

	config, err := s.tlsConfig()
	if err != nil {
		return err
	}

	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	tlsListener := tls.NewListener(app.FilterListener(app, trygo.DefaultFilterListener(app, l)), config)
	return s.Serve(tlsListener)
}

func strSliceContains(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func (s *TLSGracefulServer) tlsConfig() (*tls.Config, error) {
	config := &tls.Config{}
	if !strSliceContains(config.NextProtos, "http/1.1") {
		config.NextProtos = append(config.NextProtos, "http/1.1")
	}
	configHasCert := len(config.Certificates) > 0 || config.GetCertificate != nil
	if !configHasCert || s.CertFile != "" || s.KeyFile != "" {
		var err error
		config.Certificates = make([]tls.Certificate, 1)
		config.Certificates[0], err = tls.LoadX509KeyPair(s.CertFile, s.KeyFile)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}
