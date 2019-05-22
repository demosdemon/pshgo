package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/go-playground/lars"
	"github.com/sirupsen/logrus"

	"github.com/demosdemon/pshgo"
	"github.com/demosdemon/pshgo/cmd/serve/errors"
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

	Handler = func(ctx *Context) error
)

func New(g *Globals) *Server {
	s := Server{
		LARS: lars.New(),
	}

	tpl := func(*Context) error { return nil }

	s.RegisterContext(newContext(g))
	s.RegisterCustomHandler(Handler(tpl), castContext)

	s.Use(
		// order is important
		middleware.NewRequest,
		middleware.Log,
		middleware.Recover,
	)

	configurators.Configure(s)

	return &s
}

func (c *Context) Log() middleware.Logger {
	v, _ := c.Value(middleware.LogContextKey).(middleware.Logger)
	return v
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
	var err error
	if hdlr, ok := handler.(Handler); ok {
		if ctx, ok := c.(*Context); ok {
			err = hdlr(ctx)
		} else {
			err = errors.InternalServerError("invalid context", nil)
		}
	} else {
		err = errors.InternalServerError("invalid handler", nil)
	}

	if err != nil {
		if c.Response().Committed() {
			logrus.WithError(err).Error("error received after committing the result")
			return
		}

		logrus.WithError(err).Error("handler returned an error")

		var err2 error

		if httpError, ok := err.(errors.HTTPError); ok {
			err2 = c.JSON(httpError.StatusCode, err)
		} else {
			err2 = c.Text(http.StatusInternalServerError, err.Error())
		}

		if err2 != nil {
			logrus.WithError(err2).Error("error while writing error response")
		}
	}
}

func newContext(g *Globals) lars.ContextFunc {
	return func(l *lars.LARS) lars.Context {
		return &Context{
			Ctx:     lars.NewContext(l),
			Globals: g,
		}
	}
}
