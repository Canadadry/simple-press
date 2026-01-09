package config

import (
	"app/pkg/sflag"
	"flag"
	"fmt"
)

var (
	SerializerDateTimeFormat = "2006-01-02T15:04:05-07:00"
)

type Parameters struct {
	Action      string `flag:"action" default:"public" desc:"mode"`
	Port        int    `flag:"port" default:"8080" desc:"port to lsiten to"`
	DatabaseUrl string `flag:"db-url" default:"blog.sqlite" desc:"sqlite database filename"`
	FrontUrl    string `flag:"front-url" default:"http://localhost:5173" desc:"front url"`
	Out         string `flag:"out" default:"" desc:"output filename, stdout if blank"`
}

func New(args []string) (Parameters, error) {
	config := ""
	c := Parameters{}
	f := flag.NewFlagSet("app", flag.ContinueOnError)
	f.StringVar(&config, "config", "", "config file to define app param")
	err := sflag.Parse(f, &c)
	if err != nil {
		return c, fmt.Errorf("cannot parse param struct : %w", err)
	}
	err = f.Parse(args)
	if err != nil {
		return c, fmt.Errorf("cannot cli args : %w", err)
	}
	return c, nil
}
