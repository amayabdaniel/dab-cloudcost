package main

import (
	"os"

	"github.com/amayabdaniel/dab-cloudcost/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
