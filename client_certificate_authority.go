package pshgo

import (
	"crypto/x509"
	"encoding/pem"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ClientCertificateAuthority struct {
	*x509.Certificate
}

func (v ClientCertificateAuthority) MarshalText() ([]byte, error) {
	logrus.Trace("ClientCertificateAuthority.MarshalText")
	var block = pem.Block{
		Type:  "CERTIFICATE",
		Bytes: v.Raw,
	}
	data := pem.EncodeToMemory(&block)
	return data, nil
}

func (v *ClientCertificateAuthority) UnmarshalText(text []byte) error {
	logrus.Trace("ClientCertificateAuthority.UnmarshalText")
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

	*v = ClientCertificateAuthority{Certificate: cert}
	return nil
}
