package env

import (
	"github.com/go-playground/lars"

	"github.com/demosdemon/pshgo"
	"github.com/demosdemon/pshgo/cmd/serve/server"
)

func init() {
	server.RegisterConfigurator("/env", func(g lars.IRouteGroup) {
		g.Get("", GetEnv)
		g.Get("/application", GetApplication)
	})
}

func GetEnv(c *server.Context) {
	env := pshgo.CloneProvider(c).(pshgo.MapProvider)
	_ = c.JSON(200, env)
}

func GetApplication(c *server.Context) {
	app := c.GetApplication()
	if app == nil {
		_ = c.Text(400, "Not Found")
	} else {
		_ = c.JSON(200, app)
	}
}
