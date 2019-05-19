package pshgo

import (
	"fmt"
	"os"
	"strings"
)

type Provider interface {
	Lookup(key string) (string, bool)
	Environ() []string
	SetEnv(key, value string) error
	UnsetEnv(key string) error
	GetEnv(key string) string
}

type PlatformProvider interface {
	Provider
	Prefix() string
}

type OSProvider struct{}

func (OSProvider) Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

func (OSProvider) Environ() []string {
	return os.Environ()
}

func (OSProvider) SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

func (OSProvider) UnsetEnv(key string) error {
	return os.Unsetenv(key)
}

func (OSProvider) GetEnv(key string) string {
	return os.Getenv(key)
}

type MapProvider map[string]string

func (p MapProvider) Lookup(key string) (string, bool) {
	v, ok := p[key]
	return v, ok
}

func (p MapProvider) Environ() []string {
	rv := make([]string, 0, len(p))
	for k, v := range p {
		rv = append(rv, fmt.Sprintf("%s=%s", k, v))
	}
	return rv
}

func (p MapProvider) SetEnv(key, value string) error {
	p[key] = value
	return nil
}

func (p MapProvider) UnsetEnv(key string) error {
	delete(p, key)
	return nil
}

func (p MapProvider) GetEnv(key string) string {
	return p[key]
}

var DefaultProvider Provider = OSProvider{}

func CloneProvider(p Provider) Provider {
	environ := p.Environ()
	rv := make(MapProvider, len(environ))
	for _, env := range environ {
		if idx := strings.Index(env, "="); idx > 0 {
			rv[env[:idx]] = env[idx+1:]
		}
	}
	return rv
}
