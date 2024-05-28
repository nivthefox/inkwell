package main

import (
	"flag"
	"github.com/nivthefox/inkwell/config"
	"github.com/nivthefox/inkwell/processor"
)

func main() {
	var err error

	path := flag.String("config", "", "path to the config file")
	flag.Parse()

	if *path == "" {
		panic("no config file provided")
	}

	// Load the config file
	var cfg *config.InkwellConfig
	cfg, err = config.NewInkwellConfig(*path)
	if err != nil {
		panic(err)
	}

	// Process the config
	err = processor.ProcessBook(*cfg)
	if err != nil {
		panic(err)
	}
}
