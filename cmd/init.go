package cmd

import (
	"fmt"

	_ "hyphon/sampcon/libcontainer/nsenter"

	"github.com/urfave/cli"
)

var InitCommand = cli.Command{
	Name:  "init",
	Usage: "initialize the namespaces and launch the process (do not call it outside of runc)",
	Action: func(context *cli.Context) error {
		err := StartInitialization()
		if err != nil {
			return err
		}
		fmt.Printf("init\n")
		return err
	},
}
