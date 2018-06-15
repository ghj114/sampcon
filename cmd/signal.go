package cmd

import (
	"fmt"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

const PR_SET_CHILD_SUBREAPER = 36

type signalHandler struct {
	signals chan os.Signal
	//notifySocket *notifySocket
}

func newSignalHandler(reaper bool) *signalHandler {
	if reaper == true {
		err := unix.Prctl(PR_SET_CHILD_SUBREAPER, uintptr(1), 0, 0, 0)
		if err != nil {
			panic("err const PR_SET_CHILD_SUBREAPER = 36")
		}
	}
	s := make(chan os.Signal, 2048)
	signal.Notify(s)
	return &signalHandler{
		signals: s,
	}
}

func (sighdl *signalHandler) forward() (int, error) {
	for s := range sighdl.signals {
		switch s {
		case unix.SIGWINCH:
			fmt.Printf("receive SIGWINCH\n")
		case unix.SIGCHLD:
			fmt.Printf("receive SIGCHLD\n")
		case unix.SIGINT:
			fmt.Printf("receive SIGINT\n")
			signal.Reset(unix.SIGINT)
		default:
			fmt.Printf("receive %v\n", s)
		}
	}

	return 0, nil
}
