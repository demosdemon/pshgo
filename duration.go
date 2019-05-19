package pshgo

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Duration struct {
	time.Duration
}

func (v Duration) MarshalText() ([]byte, error) {
	logrus.Trace("Duration.MarshalText")
	return []byte(v.String()), nil
}

func (v *Duration) UnmarshalText(text []byte) (err error) {
	logrus.Trace("Duration.UnmarshalText")
	v.Duration, err = time.ParseDuration(string(text))
	return err
}
