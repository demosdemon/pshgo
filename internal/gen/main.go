package main

import (
	"bytes"
	"flag"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/sqs/goreturns/returns"
	"golang.org/x/tools/imports"
)

var (
	output = flag.String("output", "/dev/stdout", "Define where the generated code is written.")
	local  = flag.String("local", "github.com/demosdemon", "Override the locals prefix for go imports")
)

func writeOutput(schema Render, output, local string) error {
	var buf bytes.Buffer
	err := schema.Render(&buf)
	if err != nil {
		return err
	}

	imports.LocalPrefix = local
	current, err := imports.Process(output, buf.Bytes(), &imports.Options{
		Comments:  true,
		TabIndent: true,
		TabWidth:  4,
	})
	if err != nil {
		return err
	}

	// couldn't find a test case where this returns an error
	current, _ = returns.Process("", output, current, &returns.Options{
		RemoveBareReturns: true,
	})

	return ioutil.WriteFile(output, current, 0666)
}

func main() {
	flag.Parse()
	err := writeOutput(data, *output, *local)
	if err != nil {
		logrus.WithError(err).WithField("output", *output).Fatal("unable to write output")
	}
}
