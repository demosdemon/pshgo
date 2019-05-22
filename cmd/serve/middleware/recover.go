package middleware

import (
	"github.com/go-playground/lars"
	"github.com/sirupsen/logrus"

	"github.com/demosdemon/pshgo/cmd/serve/cpanic"
)

type PanicHandler interface {
	HandlePanic(p *cpanic.Panic)
}

func HandlePanic(c lars.Context, p *cpanic.Panic) {
	if ph, ok := c.(PanicHandler); ok {
		ph.HandlePanic(p)
	}

	res := c.Response()
	if res.Status() != 200 || res.Committed() {
		return
	}

	_ = c.JSON(500, map[string]interface{}{"panic": p.Value})
}

func Recover(c lars.Context) {
	if req, ok := c.Value(RequestContextKey).(*Request); ok {
		defer cpanic.Recover(func(p *cpanic.Panic) {
			req.Panic = p

			var log Logger
			if l, ok := c.(LogContext); ok {
				log = l.Log()
			} else {
				log = logrus.WithFields(req.Fields())
			}

			log.
				WithFields(logrus.Fields{
					"panic":  p.Value,
					"offset": p.Time.Sub(req.Start).String(),
				}).
				Error("recovering from panic")

			HandlePanic(c, p)
		})
	}

	c.Next()
}
