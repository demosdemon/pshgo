package pshgo

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/hashicorp/go-multierror"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

var (
	DefaultProvider Provider = OSProvider{}
)

type (
	Provider interface {
		Lookup(key string) (string, bool)
		Environ() []string
		SetEnv(key, value string) error
		UnsetEnv(key string) error
		GetEnv(key string) string
	}

	PlatformProvider interface {
		Provider
		Prefix() string
	}

	OSProvider struct{}

	MapProvider map[string]string

	LayeredProvider []Provider

	ProviderFunctor func(p Provider) error
)

func ReadEnviron(r io.Reader) (MapProvider, error) {
	hash, err := godotenv.Parse(r)
	return MapProvider(hash), err
}

func ParseEnviron(s []string) (MapProvider, error) {
	var buf bytes.Buffer
	for _, line := range s {
		_, _ = fmt.Fprintln(&buf, line)
	}
	return ReadEnviron(&buf)
}

func CloneProvider(p Provider) Provider {
	hash, _ := ParseEnviron(p.Environ())
	return hash
}

func (OSProvider) Lookup(key string) (string, bool) {
	logrus.WithField("key", key).Trace("OSProvider.Lookup")
	return os.LookupEnv(key)
}

func (OSProvider) Environ() []string {
	logrus.Trace("OSProvider.Environ")
	return os.Environ()
}

func (OSProvider) SetEnv(key, value string) error {
	logrus.
		WithField("key", key).
		WithField("value", value).
		Trace("OSProvider.SetEnv")
	return os.Setenv(key, value)
}

func (OSProvider) UnsetEnv(key string) error {
	logrus.WithField("key", key).Trace("OSProvider.UnsetEnv")
	return os.Unsetenv(key)
}

func (OSProvider) GetEnv(key string) string {
	logrus.WithField("key", key).Trace("OSProvider.GetEnv")
	return os.Getenv(key)
}

func (p MapProvider) Lookup(key string) (string, bool) {
	logrus.WithField("key", key).Trace("MapProvider.Lookup")
	v, ok := p[key]
	return v, ok
}

func (p MapProvider) Environ() []string {
	logrus.Trace("MapProvider.Environ")
	rv := make([]string, 0, len(p))
	for k, v := range p {
		rv = append(rv, fmt.Sprintf("%s=%s", k, v))
	}
	return rv
}

func (p MapProvider) SetEnv(key, value string) error {
	logrus.
		WithField("key", key).
		WithField("value", value).
		Trace("MapProvider.SetEnv")
	p[key] = value
	return nil
}

func (p MapProvider) UnsetEnv(key string) error {
	logrus.WithField("key", key).Trace("MapProvider.UnsetEnv")
	delete(p, key)
	return nil
}

func (p MapProvider) GetEnv(key string) string {
	logrus.WithField("key", key).Trace("MapProvider.GetEnv")
	return p[key]
}

func (lp *LayeredProvider) Push(p Provider) {
	*lp = append(LayeredProvider{p}, *lp...)
}

func (lp *LayeredProvider) Pop() Provider {
	p := (*lp)[0]
	*lp = (*lp)[1:]
	return p
}

func (lp LayeredProvider) ForEach(fn ProviderFunctor) error {
	for _, p := range lp {
		err := fn(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (lp LayeredProvider) First(fn ProviderFunctor) error {
	var result error

	for _, p := range lp {
		err := fn(p)
		if err == nil {
			return nil
		}
		result = multierror.Append(result, err)
	}

	return result
}

func (lp LayeredProvider) Lookup(key string) (rv string, ok bool) {
	logrus.WithField("key", key).Trace("LayeredProvider.Lookup")
	err := lp.First(func(p Provider) error {
		rv, ok = p.Lookup(key)
		if ok {
			return nil
		}
		return errors.New("not found")
	})

	if err != nil {
		return "", false
	}

	return rv, ok
}

func (lp LayeredProvider) Environ() []string {
	logrus.Trace("LayeredProvider.Environ")
	var rv []string

	_ = lp.ForEach(func(p Provider) error {
		rv = append(rv, p.Environ()...)
		return nil
	})

	// reverse list so duplicates are properly masked
	for a, b := 0, len(rv)-1; a < b; a, b = a+1, b-1 {
		rv[a], rv[b] = rv[b], rv[a]
	}

	hash, _ := ParseEnviron(rv)
	return hash.Environ()
}

func (lp LayeredProvider) SetEnv(key, value string) error {
	logrus.
		WithField("key", key).
		WithField("value", value).
		Trace("LayeredProvider.SetEnv")
	return lp.First(func(p Provider) error {
		return p.SetEnv(key, value)
	})
}

func (lp LayeredProvider) UnsetEnv(key string) error {
	logrus.WithField("key", key).Trace("LayeredProvider.UnsetEnv")
	return lp.ForEach(func(p Provider) error {
		return p.UnsetEnv(key)
	})
}

func (lp LayeredProvider) GetEnv(key string) string {
	logrus.WithField("key", key).Trace("LayeredProvider.GetEnv")
	v, _ := lp.Lookup(key)
	return v
}
