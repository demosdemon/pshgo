package middleware

import (
	"time"

	"github.com/go-playground/lars"
	"github.com/sirupsen/logrus"
)

const (
	LogContextKey = "github.com/demosdemon/pshgo/cmd/serve/middleware/LogContextKey"
)

type (
	Logger = *logrus.Entry

	LogContext interface {
		Log() Logger
	}
)

func Log(c lars.Context) {
	log := logrus.WithFields(logrus.Fields{
		"start": time.Now(),
		"url":   c.Request().URL.String(),
	})

	req, ok := c.Value(RequestContextKey).(*Request)

	if ok {
		log = log.WithField("request_id", req.ID)
		c.Set(LogContextKey, log)

		log = log.WithFields(req.Fields())
	}

	log.Info("start")

	tick := time.NewTicker(time.Second * 30)
	defer tick.Stop()
	go func() {
		for t := range tick.C {
			log.WithField("delay", t.Sub(req.Start).String()).Info("pending")
		}
	}()

	c.Next()

	if ok {
		req.UpdateLARS(c.Response())
		log = log.WithFields(req.Fields())
	}

	log.Info("end")
}
