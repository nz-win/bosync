package main

import (
	"backorder_updater/cmd/bosync/commands"
	"backorder_updater/internal/pkg"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	envRoot, err := os.Getwd()
	pkg.CheckAndPanic(err)
	err = pkg.LoadEnv(envRoot, ".env")
	pkg.CheckAndPanic(err)

	app := &cli.App{
		Name: "Backorder Sync - Tools for managing backorder sync",
		Commands: []*cli.Command{
			commands.LogsCommand,
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
