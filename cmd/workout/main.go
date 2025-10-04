package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/yrss1/workout/internal/pkg/container"
)

func main() {
	app := cli.App{
		Name:    "Workout",
		Usage:   "service",
		Version: "0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
				EnvVars: []string{"HORECA_CONFIG_PATH"},
			},
		},
		Action: runApp,
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func runApp(ctx *cli.Context) error {
	configPath := ctx.String("config")
	container, err := container.NewContainer(configPath)
	if err != nil {
		return err
	}
	if err := container.Start(ctx.Context); err != nil {
		return err
	}
	return nil
}
