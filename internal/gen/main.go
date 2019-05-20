package main

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"os"

	"github.com/octago/sflags/gen/gflag"
	"github.com/pkg/errors"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/sirupsen/logrus"
	"github.com/sqs/goreturns/returns"
	"golang.org/x/tools/imports"
)

var (
	ErrNoChange = errors.New("no change detected")
)

type Config struct {
	Path              string `desc:"path of the generated file"`
	Diff              bool   `desc:"display a diff instead of rewriting files"`
	Write             bool   `desc:"write the generated file"`
	ExitCode          bool   `desc:"exit with a failure code if no change was detected"`
	NoClobber         bool   `desc:"fail to write a file if it already exists"`
	Imports           bool   `desc:"run goimports on the file post generation"`
	Returns           bool   `desc:"run goreturns on the file post generation"`
	Local             string `desc:"put imports beginning with this string after 3rd-party packages (see goimports)"`
	AllErrors         bool   `desc:"report all errors (not just the first 10 on different lines)"`
	Comments          bool   `desc:"keep comments"`
	TabIndent         bool   `desc:"use tabs for indent"`
	TabWidth          int    `desc:"set tab width"`
	FormatOnly        bool   `desc:"disable the insertion and deletions of imports"`
	PrintErrors       bool   `desc:"print non-fatal typechecking errors to stderr"`
	RemoveBareReturns bool   `desc:"remove bare returns"`
}

func NewConfig(args []string) (*Config, error) {
	cfg := &Config{
		Path:      "/dev/stdout",
		Imports:   true,
		Returns:   true,
		Local:     "github.com/demosdemon",
		Comments:  true,
		TabIndent: true,
		TabWidth:  8,
	}

	fs := flag.NewFlagSet("gen", flag.ContinueOnError)
	must(gflag.ParseTo(cfg, fs))

	err := fs.Parse(args)
	if err == flag.ErrHelp {
		fs.PrintDefaults()
	}
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func main() {
	Execute(os.Args[1:], data)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func Execute(args []string, data Render) {
	cfg, err := NewConfig(args)
	if err != nil {
		logrus.WithError(err).Fatal()
	}

	err = cfg.Execute(data)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
}

func (c *Config) ImportsOptions() *imports.Options {
	return &imports.Options{
		Fragment:   false,
		AllErrors:  c.AllErrors,
		Comments:   c.Comments,
		TabIndent:  c.TabIndent,
		TabWidth:   c.TabWidth,
		FormatOnly: c.FormatOnly,
	}
}

func (c *Config) ReturnsOptions() *returns.Options {
	return &returns.Options{
		Fragment:          false,
		PrintErrors:       c.PrintErrors,
		AllErrors:         c.AllErrors,
		RemoveBareReturns: c.RemoveBareReturns,
	}
}

func (c *Config) ReadPrevious() ([]byte, error) {
	switch c.Path {
	case "/dev/stdout", "/dev/stderr":
		return nil, nil
	default:
		data, err := ioutil.ReadFile(c.Path)
		if os.IsNotExist(err) {
			return nil, nil
		}
		return data, err
	}
}

func (c *Config) Render(r Render) ([]byte, error) {
	var buf bytes.Buffer
	err := r.Render(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Config) GoImports(data []byte) ([]byte, error) {
	if c.Imports {
		return imports.Process(c.Path, data, c.ImportsOptions())
	}
	return data, nil
}

func (c *Config) GoReturns(data []byte) ([]byte, error) {
	if c.Returns {
		return returns.Process("", c.Path, data, c.ReturnsOptions())
	}
	return data, nil
}

func (c *Config) PrintDiff(w io.Writer, prev, current []byte) error {
	if !c.Diff {
		return nil
	}

	return difflib.WriteUnifiedDiff(w, difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(prev)),
		B:        difflib.SplitLines(string(current)),
		FromFile: "a",
		FromDate: "",
		ToFile:   "b",
		ToDate:   "",
		Context:  5,
	})
}

func (c *Config) WriteFile(prev, current []byte) error {
	if bytes.Equal(prev, current) {
		return ErrNoChange
	}

	if !c.Write {
		return nil
	}

	if c.NoClobber {
		_, err := os.Stat(c.Path)
		if err == nil {
			return &os.PathError{
				Op:   "WriteFile",
				Path: c.Path,
				Err:  os.ErrExist,
			}
		}
	}

	return ioutil.WriteFile(c.Path, current, 0666)
}

func (c *Config) Execute(r Render) error {
	prev, err := c.ReadPrevious()
	if err != nil {
		return errors.Wrap(err, "error reading existing file")
	}

	current, err := c.Render(r)
	if err != nil {
		return errors.Wrap(err, "error executing render step")
	}

	current, err = c.GoImports(current)
	if err != nil {
		return errors.Wrap(err, "error executing goimports")
	}

	current, err = c.GoReturns(current)
	if err != nil {
		return errors.Wrap(err, "error executing goreturns")
	}

	_ = c.PrintDiff(os.Stdout, prev, current)

	return c.WriteFile(prev, current)
}
