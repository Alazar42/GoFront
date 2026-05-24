package main

import (
	"os"

	"gofront/internal/cli"
)

func main() {
	if err := cli.Execute(os.Args[1:]); err != nil {
		cli.ExitWithError(err)
	}
}
