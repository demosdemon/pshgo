package pshgo

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type TLSVersion uint16

const (
	TLSv10 TLSVersion = iota + 0x0301
	TLSv11
	TLSv12
	TLSv13
	TLSv14
)

var tlsNameMapping = map[TLSVersion]string{
	TLSv10: "TLSv1.0",
	TLSv11: "TLSv1.1",
	TLSv12: "TLSv1.2",
	TLSv13: "TLSv1.3",
	TLSv14: "TLSv1.4",
}

func NewTLSVersion(name string) (TLSVersion, error) {
	logrus.Trace("NewTLSVersion")
	for k, v := range tlsNameMapping {
		if v == name {
			return k, nil
		}
	}

	return 0, fmt.Errorf("unknown TLSVersion %q", name)
}

func (v TLSVersion) String() string {
	logrus.Trace("TLSVersion.String")
	if name, ok := tlsNameMapping[v]; ok {
		return name
	}

	return fmt.Sprintf("unknown TLSVersion 0x%04x", uint16(v))
}

func (v *TLSVersion) UnmarshalText(text []byte) (err error) {
	logrus.Trace("TLSVersion.UnmarshalText")
	*v, err = NewTLSVersion(string(text))
	return err
}

func (v TLSVersion) MarshalText() ([]byte, error) {
	logrus.Trace("TLSVersion.MarshalText")
	if rv, ok := tlsNameMapping[v]; ok {
		return []byte(rv), nil
	}

	return nil, errors.New(v.String())
}
