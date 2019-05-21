package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-playground/lars"

	"github.com/demosdemon/pshgo"
	"github.com/demosdemon/pshgo/cmd/serve/middleware"
)

var (
	DefaultShutdownTimeout = 15 * time.Second
)

type (
	Globals struct {
		*pshgo.Environment
	}

	Context struct {
		*lars.Ctx
		*Globals
	}

	Server struct {
		*lars.LARS
	}

	Handler func(ctx *Context)
)

func New(g *Globals) *Server {
	s := Server{
		LARS: lars.New(),
	}

	tpl := func(*Context) {}

	s.RegisterContext(newContext(g))
	s.RegisterCustomHandler(tpl, castContext)
	s.RegisterCustomHandler(Handler(tpl), castContext)

	s.Use(
		// order is important
		middleware.Log,
		middleware.Recover,
	)

	configurators.Configure(s)

	return &s
}

func (s *Server) Serve(ctx context.Context, l net.Listener) error {
	done := make(chan error)

	srv := http.Server{Handler: s.LARS.Serve()}

	go func() {
		done <- srv.Serve(l)
	}()

	go func() {
		// wait for the context
		<-ctx.Done()

		newCtx, cancel := context.WithTimeout(context.Background(), DefaultShutdownTimeout)
		defer cancel()

		// gracefully shutdown server within 15 seconds
		err := srv.Shutdown(newCtx)

		// if graceful shutdown fails, use force
		if err != nil {
			_ = srv.Close()
		}
	}()

	// wait for server to return
	err := <-done

	// pass only interesting errors back
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}

func castContext(c lars.Context, handler lars.Handler) {
	var hdlr Handler
	if h, ok := handler.(Handler); ok {
		hdlr = h
	} else {
		hdlr = handler.(func(*Context))
	}

	ctx := c.(*Context)
	hdlr(ctx)
}

func newContext(g *Globals) lars.ContextFunc {
	return func(l *lars.LARS) lars.Context {
		return &Context{
			Ctx:     lars.NewContext(l),
			Globals: g,
		}
	}
}
