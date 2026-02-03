package main

import (
	"os"

	"github.com/enbiyagoral/sopsy/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
