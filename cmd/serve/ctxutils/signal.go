package ctxutils

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/sirupsen/logrus"
)

func CancelContextWithSignal(parent context.Context, s ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	ch := make(chan os.Signal, len(s))

	var once sync.Once

	realCancel := func() {
		once.Do(func() {
			signal.Stop(ch)
			close(ch)
			cancel()
		})
	}

	go func() {
		if sig, ok := <-ch; ok {
			logrus.WithField("signal", sig).Info("signal received")
			realCancel()
		}
	}()

	signal.Notify(ch, s...)

	return ctx, realCancel
}
