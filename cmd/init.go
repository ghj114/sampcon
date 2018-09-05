package cmd

import (
	"fmt"

	_ "hyphon/sampcon/libcontainer/nsenter"

	"hyphon/sampcon/libcontainer"

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

func StartInitialization() (err error) {
	fmt.Printf("test in init-----------------------------\n")
	err = linuxcontainer.NewContainerInit()
	return err
}
