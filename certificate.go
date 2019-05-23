package pshgo

import (
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Certificate struct {
	*x509.Certificate
}

func (v Certificate) MarshalText() ([]byte, error) {
	logrus.Trace("Certificate.MarshalText")
	var block = pem.Block{
		Type:  "CERTIFICATE",
		Bytes: v.Raw,
	}
	data := pem.EncodeToMemory(&block)
	return data, nil
}

func (v *Certificate) UnmarshalText(text []byte) error {
	logrus.Trace("Certificate.UnmarshalText")
	block, rest := pem.Decode(text)
	if block == nil {
		return errors.New("invalid PEM data")
	}
	if rest != nil && len(rest) > 0 {
		return errors.New("excess data after decoding the PEM block")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}

	*v = Certificate{Certificate: cert}
	return nil
}
