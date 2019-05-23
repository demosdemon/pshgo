package env

import (
	"github.com/demosdemon/pshgo"
	"github.com/demosdemon/pshgo/cmd/serve/errors"
	"github.com/demosdemon/pshgo/cmd/serve/server"
	"github.com/go-playground/lars"
	"github.com/joho/godotenv"
)

func init() {
	server.RegisterConfigurator("/env", func(g lars.IRouteGroup) {
		g.Get("", GetEnv)
		g.Get("/export", GetExport)
		g.Get("/application", GetApplication)
		g.Get("/routes", GetRoutes)
	})
}

func GetEnv(c *server.Context) error {
	env := pshgo.CloneProvider(c).(pshgo.MapProvider)
	return c.JSON(200, env)
}

func GetExport(c *server.Context) error {
	env := pshgo.CloneProvider(c).(pshgo.MapProvider)
	marshaled, err := godotenv.Marshal(env)
	if err != nil {
		return errors.InternalServerError("unable to marshal environment", err)
	}

	return c.Text(200, marshaled)
}

func GetApplication(c *server.Context) error {
	app := c.GetApplication()
	if app == nil {
		return errors.NotFound("Not Found", nil)
	}
	return c.JSON(200, app)
}

func GetRoutes(c *server.Context) error {
	routes := c.GetRoutes()
	if routes == nil {
		return errors.NotFound("Not Found", nil)
	}
	return c.JSON(200, routes)
}
