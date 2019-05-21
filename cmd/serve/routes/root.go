package routes

import (
	"time"

	"github.com/go-playground/lars"

	"github.com/demosdemon/pshgo/cmd/serve/server"
)

func init() {
	server.RegisterConfigurator("", func(g lars.IRouteGroup) {
		g.Get("/", Home)
		g.Any("/ping", Ping)
		g.Any("/panic", Panic)
	})
}

func Home(c *server.Context) {
	_ = c.Text(200, "Hello, World!\n")
}

func Ping(c *server.Context) {
	rv := struct {
		Message   string    `json:"msg"`
		Timestamp time.Time `json:"ts"`
	}{
		Message:   "pong",
		Timestamp: time.Now(),
	}

	_ = c.JSON(200, rv)
}

func Panic(_ *server.Context) {
	panic("this route always panics")
}
