package main

import (
	"os"

	"github.com/enbiyagoral/sopsctl/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
