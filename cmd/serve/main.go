package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/go-playground/lars/middleware"
	"github.com/joho/godotenv"
	"github.com/octago/sflags/gen/gflag"
	"github.com/sirupsen/logrus"

	. "github.com/demosdemon/pshgo"
	"github.com/demosdemon/pshgo/cmd/serve/ctxutils"
	_ "github.com/demosdemon/pshgo/cmd/serve/routes"
	"github.com/demosdemon/pshgo/cmd/serve/server"
)

func init() {
	logrus.SetLevel(logrus.TraceLevel)
}

func main() {
	Execute(os.Args[1:])
}

func Execute(args []string) {
	cfg, err := NewConfig(args)
	if err != nil {
		logrus.WithError(err).Fatal()
	}

	err = cfg.Execute()
	if err != nil {
		logrus.WithError(err).Fatal()
	}
}

type Config struct {
	Prefix          string        `desc:"the Platform.sh environment prefix"`
	DotEnv          string        `desc:"read the specified .env file if it exists; set to /dev/null to disable"`
	ShutdownTimeout time.Duration `desc:"the amount of time to wait before forcefully terminating the server upon request"`
}

func NewConfig(args []string) (*Config, error) {
	cfg := &Config{
		Prefix: "PLATFORM_",
		DotEnv: ".env",
	}

	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	must(gflag.ParseTo(cfg, fs))

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) Execute() error {
	log := logrus.WithField("config", c)

	dotenv, err := godotenv.Read(c.DotEnv)

	if os.IsNotExist(err) {
		log.Info(".env file not found")
		err = nil
	}

	if err != nil {
		log.WithError(err).Error("unable to read .env file")
		return err
	}

	environ := LayeredProvider{
		MapProvider(dotenv),
		DefaultProvider,
	}

	env := NewEnvironmentWithProvider(c.Prefix, environ)

	s := server.New(&server.Globals{
		Environment: env,
	})

	s.Use(middleware.Gzip)

	l, err := env.Listener()
	if err != nil {
		log.WithError(err).Error("unable to bind listener")
		return err
	}

	ctx, cancel := ctxutils.CancelContextWithSignal(context.Background())
	defer cancel()

	return s.Serve(ctx, l)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
