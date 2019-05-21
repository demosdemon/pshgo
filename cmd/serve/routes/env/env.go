package env

import (
	"github.com/go-playground/lars"

	"github.com/demosdemon/pshgo"
	"github.com/demosdemon/pshgo/cmd/serve/server"
)

func init() {
	server.RegisterConfigurator("/env", func(g lars.IRouteGroup) {
		g.Get("", GetEnv)
	})
}

func GetEnv(c *server.Context) {
	env := pshgo.CloneProvider(c).(pshgo.MapProvider)
	_ = c.JSON(200, env)
}
