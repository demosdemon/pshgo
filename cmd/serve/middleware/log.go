package middleware

import (
	"time"

	"github.com/go-playground/lars"
	"github.com/sirupsen/logrus"
)

func Log(c lars.Context) {
	req, ok := c.Value(RequestContextKey).(*Request)

	if ok {
		logrus.WithField("request", req).Info("start request")

		tick := time.NewTicker(time.Second)
		defer tick.Stop()

		go func() {
			for t := range tick.C {
				logrus.
					WithFields(logrus.Fields{
						"request_id": req.ID,
						"url":        req.URL,
						"delay":      t.Sub(req.Start),
					}).
					Info("in-progress request")
			}
		}()
	}

	c.Next()

	if ok {
		req.UpdateLARS(c.Response())
		logrus.WithField("request", req).Info("end request")
	}
}
