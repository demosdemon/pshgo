package middleware

import (
	"net/http"
	"time"

	"github.com/go-playground/lars"
	"github.com/sirupsen/logrus"
)

func Log(c lars.Context) {
	start := time.Now()

	c.Next()

	logrus.
		WithFields(logrus.Fields{
			"start":     start,
			"delay":     time.Since(start),
			"client_ip": c.ClientIP(),
		}).
		WithFields(requestFields(c.Request())).
		WithFields(responseFields(c.Response())).
		Info("request")
}

func requestFields(req *http.Request) logrus.Fields {
	return logrus.Fields{
		"remote_addr": req.RemoteAddr,
		"method":      req.Method,
		"url":         req.URL,
		"proto":       req.Proto,
		"user_agent":  req.UserAgent(),
	}
}

func responseFields(res *lars.Response) logrus.Fields {
	return logrus.Fields{
		"status": res.Status(),
		"size":   res.Size(),
	}
}
