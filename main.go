package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"

	"hyphon/sampcon/cmd"
)

func main() {
	app := cli.NewApp()
	app.Name = "sampcon"
	app.Usage = "sample container"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "log",
			Value: "/dev/null",
			Usage: "set the log file path.",
		},
	}
	app.Commands = []cli.Command{
		cmd.RunCommand,
		cmd.InitCommand,
	}
	app.Action = func(c *cli.Context) error {
		fmt.Println("Hello friend!")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
