package cmd

import (
	"errors"
	"fmt"

	"hyphon/sampcon/libcontainer"
	"hyphon/sampcon/libcontainer/config"
	"hyphon/sampcon/libcontainer/utils"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
)

type runner struct {
	container *linuxcontainer.Linuxcontainer
}

func (r *runner) run(process *specs.Process) (int, error) {
	fmt.Printf("process:%+v\n", process)
	cProcess, err := containerProcess(process)
	handler := newSignalHandler(false)
	_, err = r.setIO(cProcess)
	r.container.Run(cProcess)
	//go func() {
	//	status, _ := handler.forward()
	//	fmt.Printf("status:%d\n", status)
	//}()
	status, _ := handler.forward()
	fmt.Printf("status:%d\n", status)
	//select {}
	return 0, err
}

func (r *runner) setIO(process *linuxcontainer.Process) (*tty, error) {
	parent, child, err := utils.NewSockPair("console")
	process.PConsoleSocket = parent
	process.ConsoleSocket = child
	fmt.Printf("consolesocket.parent:%d,child:%d\n", parent, child)
	t := &tty{}
	//t.consoleC = make(chan error, 1)
	go func() {
		//f, err := utils.RecvFd(parent)
		//fmt.Printf("consolesocket.ttymaster:%d\n", f)
		//if err != nil {
		//	fmt.Printf("recvfd err:%s\n", err)
		//}
		//cons, err := console.ConsoleFromFile(f)
		//if err != nil {
		//	fmt.Printf("conslefromfile err:%s\n", err)
		//}
		if err := t.recvtty(parent); err != nil {
			//t.consoleC <- err
			fmt.Printf("recvtty err:%s\n", err)
		}
		//t.consoleC <- nil
		//go func() {
		//	sigchan := make(chan os.Signal, 1)
		//	fmt.Printf("in handleinterrupt\n")
		//	signal.Notify(sigchan, os.Interrupt)
		//	rec := <-sigchan
		//	fmt.Printf("receive sig in recvtty:%v\n", rec)
		//	os.Exit(0)
		//}()
	}()
	return t, err
}

func startContainer(context *cli.Context, spec *specs.Spec) (int, error) {
	id := context.Args().First()
	if id == "" {
		errEmptyID := errors.New("container id can not be empty.")
		return -1, errEmptyID
	}
	container, _ := createContainer(context, id, spec)
	fmt.Printf("container:%+v\n", container)

	r := &runner{
		container: container,
	}
	return r.run(spec.Process)
}

//func loadfactory(context *cli.Context) error {
//	container := linuxcontainer.New()
//}

func createContainer(context *cli.Context, id string, spec *specs.Spec) (*linuxcontainer.Linuxcontainer, error) {
	config, err := config.Specconv(context, spec)
	fmt.Printf("config:%+v\n", config)
	//container := loadfactory(context)
	container, _ := linuxcontainer.New(id, config)
	//container.create()
	return container, err
}

func containerProcess(process *specs.Process) (*linuxcontainer.Process, error) {
	lp := &linuxcontainer.Process{
		Args: process.Args,
		Env:  process.Env,
		Cwd:  process.Cwd,
	}

	return lp, nil
}
