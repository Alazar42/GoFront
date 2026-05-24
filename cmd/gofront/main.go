package main

import (
	"os"

	"github.com/Alazar42/GoFront/internal/cli"
)

func main() {
	if err := cli.Execute(os.Args[1:]); err != nil {
		cli.ExitWithError(err)
	}
}
