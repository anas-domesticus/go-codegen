package main

import (
	"flag"
	"fmt"
	"github.com/anas-domesticus/go-codegen/pkg"
	"os"
)

func main() {
	configFile := flag.String("config", "generate.yml", "Path to the config file")
	flag.Parse()

	t, err := pkg.NewTemplaterFromPath(*configFile)
	if err != nil {
		fmt.Printf("Error creating templater: %s\n", err)
		os.Exit(1)
	}
	err = t.GenerateFiles()
	if err != nil {
		fmt.Printf("Error running templater: %s\n", err)
		os.Exit(1)
	}
}
