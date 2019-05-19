package pshgo

import (
	"net"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Environment struct {
	p      Provider
	prefix string
}

func NewEnvironment(prefix string) *Environment {
	return &Environment{
		prefix: prefix,
	}
}

func (e *Environment) Lookup(key string) (string, bool) {
	return e.provider().Lookup(key)
}

func (e *Environment) Environ() []string {
	return e.provider().Environ()
}

func (e *Environment) SetEnv(key, value string) error {
	return e.provider().SetEnv(key, value)
}

func (e *Environment) UnsetEnv(key string) error {
	return e.provider().UnsetEnv(key)
}

func (e *Environment) GetEnv(key string) string {
	return e.provider().GetEnv(key)
}

func (e *Environment) provider() Provider {
	p := e.p
	if p == nil {
		p = DefaultProvider
	}
	return p
}

func (e *Environment) SetProvider(p Provider) {
	e.p = p
}

func (e *Environment) Prefix() string {
	return e.prefix
}

func (e *Environment) Listener() (net.Listener, error) {
	logrus.Trace("NewListener")

	if socket, ok := e.GetSocket(); ok {
		logrus.WithField("socket", socket).Debug("found SOCKET")
		return net.Listen("unix", socket)
	}

	if port, ok := e.GetPort(); ok {
		logrus.WithField("port", port).Debug("found PORT")
		return net.Listen("tcp", net.JoinHostPort("127.0.0.1", port))
	}

	return nil, errors.New("found neither SOCKET nor PORT")
}

func (e *Environment) Variable(key string) (interface{}, bool) {
	if vars, ok := e.GetVariables(); ok {
		v, ok := vars[key]
		return v, ok
	}
	return nil, false
}
