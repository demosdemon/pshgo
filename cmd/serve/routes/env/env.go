package env

import (
	"github.com/go-playground/lars"
	"github.com/joho/godotenv"

	"github.com/demosdemon/pshgo"
	"github.com/demosdemon/pshgo/cmd/serve/server"
)

func init() {
	server.RegisterConfigurator("/env", func(g lars.IRouteGroup) {
		g.Get("", GetEnv)
		g.Get("/export", GetExport)
		g.Get("/application", GetApplication)
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
		return err
	}

	return c.Text(200, marshaled)
}

func GetApplication(c *server.Context) error {
	app := c.GetApplication()
	if app == nil {
		return c.Text(400, "Not Found")
	}
	return c.JSON(200, app)
}
