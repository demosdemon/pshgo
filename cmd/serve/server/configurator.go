package server

import (
	"sort"

	"github.com/go-playground/lars"
	"github.com/sirupsen/logrus"
)

var configurators = make(Configurators)

type (
	Configurators map[string]Configurator
	Configurator  func(g lars.IRouteGroup)
)

func RegisterConfigurator(path string, fn Configurator) {
	if _, ok := configurators[path]; ok {
		logrus.WithField("path", path).Panic("duplicate configurator for path")
	}

	configurators[path] = fn
}

func (c Configurators) Configure(g lars.IRouteGroup) {
	paths := make([]string, 0, len(c))
	for k := range c {
		paths = append(paths, k)
	}

	sort.Strings(paths)

	for _, k := range paths {
		fn := c[k]
		fn(g.Group(k))
	}
}
