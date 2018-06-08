package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

var RunCommand = cli.Command{
	Name:      "run",
	Usage:     "create and run a container.",
	ArgsUsage: "",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "bundle, b",
			Value: "",
			Usage: "",
		},
	},
	Action: func(context *cli.Context) error {
		spec, err := setupSpec(context)
		if err != nil {
			return err
		}
		status, err := startContainer(context, spec)
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", status)
		return err
	},
}
