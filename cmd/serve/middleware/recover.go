package middleware

import (
	"time"

	"github.com/go-playground/lars"
	"github.com/sirupsen/logrus"

	"github.com/demosdemon/pshgo/cmd/serve/cpanic"
)

func Recover(c lars.Context) {
	start := time.Now()

	defer cpanic.Recover(func(p *cpanic.Panic) {
		logrus.
			WithFields(requestFields(c.Request())).
			WithFields(logrus.Fields{
				"start": start,
				"delay": time.Since(start),
			}).
			Error("recovering from panic")

		logrus.Printf("\n%s", p)

		resp := c.Response()
		if resp.Status() != 200 || resp.Committed() {
			return
		}

		text, _ := p.MarshalText()
		_ = c.TextBytes(500, text)
	})()

	c.Next()
}
