package pshgo_test

import (
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/demosdemon/pshgo"
)

const rootCertificatePEM = `-----BEGIN CERTIFICATE-----
MIIEHzCCAwegAwIBAgIQbvVcMf/jHqiQXeb4tqfG3zANBgkqhkiG9w0BAQsFADCB
mDELMAkGA1UEBhMCVVMxEjAQBgNVBAgTCUxvdWlzaWFuYTESMBAGA1UEBxMJTGFm
YXlldHRlMRswGQYDVQQKExJMZUJsYW5jIENvZGVzLCBMTEMxETAPBgNVBAsTCFNl
Y3VyaXR5MTEwLwYDVQQDEyhMZUJsYW5jIENvZGVzLCBMTEMgQ2VydGlmaWNhdGUg
QXV0aG9yaXR5MB4XDTE5MDUwOTAxMjAzNVoXDTIwMDUwODAxMjAzNFowgZgxCzAJ
BgNVBAYTAlVTMRIwEAYDVQQIEwlMb3Vpc2lhbmExEjAQBgNVBAcTCUxhZmF5ZXR0
ZTEbMBkGA1UEChMSTGVCbGFuYyBDb2RlcywgTExDMREwDwYDVQQLEwhTZWN1cml0
eTExMC8GA1UEAxMoTGVCbGFuYyBDb2RlcywgTExDIENlcnRpZmljYXRlIEF1dGhv
cml0eTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALVot8Kpef49RSLl
NZjJopXIYJvs4hxO0GV1Q+A9oU7jnn27dDEoxZTSzdc/wu9rMtwUDEJKhRJVmvcs
4wnmbB7kqoNi/2iyxa5JCw4HArHXr4b+b+lKlHYNng0JEDOJIFV+A/orlmR/MiKq
LFN5qEcUPuGYBFUbwTbgQonpw99WQGN7Bu+9yDDAA2CXaJZhJkxmfUMeW7urg0ZQ
5pUV3JDn0k+yrPALu0xJOXd0RiYALql6CJDeU4tJsZo5Mu3aWNJww0it/sr4bJ/T
/4fE60rboOIMsgM5RtZItZmYFZ+pavBFmXPuX4XyIjqi7TbJKRhbRfpzekxnmlTP
9KB3o7UCAwEAAaNjMGEwDgYDVR0PAQH/BAQDAgGGMA8GA1UdEwEB/wQFMAMBAf8w
HQYDVR0OBBYEFOuVJ3ry0B49g93rUDMsjqaNVK0hMB8GA1UdIwQYMBaAFOuVJ3ry
0B49g93rUDMsjqaNVK0hMA0GCSqGSIb3DQEBCwUAA4IBAQA+osWvK2xQiuuxocbm
YCPOzh85UOMFhhzYOPZO3IEfLJo9+jp7uBOjqziMc1BILCBcVJ63dsidRtgRZ69f
/qqsv/lz4LZwZqJhoNXOmj++bZLeIOjzexe/tXKA5RPnv/luVQl+eVLxfW17R09B
0cMRGKNhrrKsnDY8tWL1pYMM+8sN+6UZ5cb9F1/Xd78xFHpN0yOX/54bEhY2SJsq
2vBfEaQ8XlBj6nKdX2BxO+sEeUmAcnPcZbgOOX3WqELsZjFCp46FEuytR9NW4Wi1
2MP7IIJusPcrpccX1QoxaV1YPzAgJkJwA0BSG+ArCNz03x/vnjHOKxu4/m0wzxBJ
zjSq
-----END CERTIFICATE-----
`

const intermediateCertificatePEM = `-----BEGIN CERTIFICATE-----
MIIEFzCCAv+gAwIBAgIQPiFyJe+Z2C86i9lbTiNgYjANBgkqhkiG9w0BAQsFADCB
mDELMAkGA1UEBhMCVVMxEjAQBgNVBAgTCUxvdWlzaWFuYTESMBAGA1UEBxMJTGFm
YXlldHRlMRswGQYDVQQKExJMZUJsYW5jIENvZGVzLCBMTEMxETAPBgNVBAsTCFNl
Y3VyaXR5MTEwLwYDVQQDEyhMZUJsYW5jIENvZGVzLCBMTEMgQ2VydGlmaWNhdGUg
QXV0aG9yaXR5MB4XDTE5MDUwOTAxMjI1N1oXDTIwMDUwODAxMjI1N1owgZAxCzAJ
BgNVBAYTAlVTMRIwEAYDVQQIEwlMb3Vpc2lhbmExEjAQBgNVBAcTCUxhZmF5ZXR0
ZTEbMBkGA1UEChMSTGVCbGFuYyBDb2RlcywgTExDMREwDwYDVQQLEwhTZWN1cml0
eTEpMCcGA1UEAxMgTGVCbGFuYyBDb2RlcywgTExDIC0gSW50ZXJuYWwgQ0EwggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDOj84uQo+yV5AMaKaBNoOwnHh9
yLSVda0ro/eoB4wrp+IhEwFulEYP1Ybv/Maw2CFLUktscu0yWf6R1uxh+AnqLH+W
Lk3CWm4cSa55ExH6ZJCLgAk1k1ic4jMXAkDI1mjK9qetFI+yGGgyV0bxa3J+qxnh
KXlZGbn8hrD2j964LSpqy/oG8pwFrsryHIUagItzkLDVXZU1xWTR4CopCpFMkqNd
wvGY2zt35e7xHRLbOcfTh3QmzzY2IEcajcv7LJtGo4C0nhuHgFXmhmoJz02wvws2
t4YIKqXqGhLSWQlXfwKMjf5Uu4LiIHoYoNP8bGWsVzAuIYGtcCcyr0nqc4n3AgMB
AAGjYzBhMA4GA1UdDwEB/wQEAwIBhjAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQW
BBTQdbiiuFfRiqLwcRHvvmB8tbLdkTAfBgNVHSMEGDAWgBTrlSd68tAePYPd61Az
LI6mjVStITANBgkqhkiG9w0BAQsFAAOCAQEACPBSFdCEutcTehyZdO2OrYFwcAXd
ef6olDcL4bZ1MTmMk4+p7InJsP+piFQfpNyARVokovcivn2Lu/dbGQn7OjTv1iBV
uNpCcj+tuqM2K1xxoY/smXRni8CyCsOxSk75Nry+Ed0RApbL8anpqd11CIV0zVcP
OJTId7gokE5aqhss2nbo5J31X0qXbkWL2SHM6Mj/PXviqDklRmOWTRV/zFHzWB8u
zN6Ez98Foev69YyTdIOPUoZdWyjpk9sCsmDVh16rE0/FTXe8HOHIztJ/quHgKlox
uFg6+IaAqYd0EDr0F+apr3ho8ZGikG+3mg0sBlqaeaFGbpG+OldL8WHw3Q==
-----END CERTIFICATE-----
`

// openssl genrsa -3 256
const privateKey256PEM = `-----BEGIN RSA PRIVATE KEY-----
MIGqAgEAAiEA0QMC15WQJUucA9kZNd9Z51AkiJwa+kIwVsKI6kDpNXECAQMCIQCL
V1c6Y7VuMmgCkLt5P5FDqqLrz8XPMWfRV/2vZGXASwIRAPsVnobZgYRKSYaC+eCy
460CEQDVGohdmMHzylM4CWlJnbFVAhEAp2O/BJEBAtwxBFdRQHdCcwIRAI4RsD5l
1qKG4iVbm4ZpIOMCEGMHqRKMao5Iq+ibujgoRVc=
-----END RSA PRIVATE KEY-----
`

var (
	rootCertificate         *x509.Certificate
	intermediateCertificate *x509.Certificate
)

func init() {
	var err error
	rootCertificateBlock, _ := pem.Decode([]byte(rootCertificatePEM))
	rootCertificate, err = x509.ParseCertificate(rootCertificateBlock.Bytes)
	if err != nil {
		panic(err)
	}

	intermediateCertificateBlock, _ := pem.Decode([]byte(intermediateCertificatePEM))
	intermediateCertificate, err = x509.ParseCertificate(intermediateCertificateBlock.Bytes)
	if err != nil {
		panic(err)
	}
}

func TestCertificate_MarshalText(t *testing.T) {
	cases := []struct {
		name    string
		v       Certificate
		want    string
		wantErr bool
	}{
		{
			name:    "RootCertificate",
			v:       Certificate{Certificate: rootCertificate},
			want:    rootCertificatePEM,
			wantErr: false,
		},
		{
			name:    "IntermediateCertificate",
			v:       Certificate{Certificate: intermediateCertificate},
			want:    intermediateCertificatePEM,
			wantErr: false,
		},
		{
			name:    "Nil",
			wantErr: true,
		},
	}

	t.Parallel()
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				r := recover()
				assert.True(t, (r != nil) == c.wantErr)
			}()

			data, err := c.v.MarshalText()
			assert.NoError(t, err)
			assert.Equal(t, c.want, string(data))
		})
	}
}

func TestCertificate_UnmarshalText(t *testing.T) {
	cases := []struct {
		name    string
		text    string
		want    Certificate
		wantErr bool
	}{
		{
			name:    "RootCertificate",
			text:    rootCertificatePEM,
			want:    Certificate{Certificate: rootCertificate},
			wantErr: false,
		},
		{
			name:    "IntermediateCertificate",
			text:    intermediateCertificatePEM,
			want:    Certificate{Certificate: intermediateCertificate},
			wantErr: false,
		},
		{
			name:    "Invalid PEM Data",
			text:    "-----HELLO WORLD-----\n",
			want:    Certificate{},
			wantErr: true,
		},
		{
			name:    "Excess data after PEM Data",
			text:    rootCertificatePEM + intermediateCertificatePEM,
			want:    Certificate{},
			wantErr: true,
		},
		{
			name:    "Invalid Certificate",
			text:    privateKey256PEM,
			want:    Certificate{},
			wantErr: true,
		},
	}

	t.Parallel()
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			var got Certificate
			err := got.UnmarshalText([]byte(c.text))
			assert.True(t, (err != nil) == c.wantErr)
			assert.EqualValues(t, c.want, got)
		})
	}
}
