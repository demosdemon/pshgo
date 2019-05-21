package routes

import (
	"encoding/json"
	"io"
	"math/rand"
	"time"

	"github.com/go-playground/lars"

	"github.com/demosdemon/pshgo/cmd/serve/server"
)

func init() {
	server.RegisterConfigurator("", func(g lars.IRouteGroup) {
		g.Get("/", Home)
		g.Any("/ping", Ping)
		g.Any("/panic", Panic)
		g.Get("/drip", Drip)
	})
}

func Home(c *server.Context) error {
	return c.Text(200, "Hello, World!\n")
}

func Ping(c *server.Context) error {
	rv := struct {
		Message   string    `json:"msg"`
		Timestamp time.Time `json:"ts"`
	}{
		Message:   "pong",
		Timestamp: time.Now(),
	}

	return c.JSON(200, rv)
}

func Panic(_ *server.Context) error {
	panic("this route always panics")
}

func Drip(c *server.Context) error {
	c.Response().Header().Set("Content-Type", "application/stream+json")

	count := rand.Intn(100)
	var err error
	c.Stream(func(w io.Writer) bool {
		if count <= 0 {
			return false
		}
		count--

		drop := struct {
			Message string `json:"msg"`
			Count   int    `json:"count"`
		}{
			Message: "drip",
			Count:   count,
		}

		data, _ := json.Marshal(drop)

		_, err = w.Write(data)
		if err != nil {
			return false
		}

		_, err = w.Write([]byte("\n"))
		if err != nil {
			return false
		}

		time.Sleep(time.Second)

		return true
	})

	return err
}
