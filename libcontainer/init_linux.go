package linuxcontainer

import (
	"encoding/json"
	"fmt"
	"hyphon/sampcon/libcontainer/config"
	"hyphon/sampcon/libcontainer/utils"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/containerd/console"
	"golang.org/x/sys/unix"
)

// initConfig is used for transferring parameters from Exec() to Init()
type InitConfig struct {
	Args []string `json:"args"`
	Env  []string `json:"env"`
	Cwd  string   `json:"cwd"`
	//Capabilities     *config.Capabilities `json:"capabilities"`
	ProcessLabel     string         `json:"process_label"`
	AppArmorProfile  string         `json:"apparmor_profile"`
	NoNewPrivileges  bool           `json:"no_new_privileges"`
	User             string         `json:"user"`
	AdditionalGroups []string       `json:"additional_groups"`
	Config           *config.Config `json:"config"`
	//Networks         []*network     `json:"network"`
	PassedFilesCount int    `json:"passed_files_count"`
	ContainerId      string `json:"containerid"`
	//Rlimits          []config.Rlimit       `json:"rlimits"`
	CreateConsole bool   `json:"create_console"`
	ConsoleWidth  uint16 `json:"console_width"`
	ConsoleHeight uint16 `json:"console_height"`
	Rootless      bool   `json:"rootless"`
}

func setconsole(socket *os.File) error {
	defer socket.Close()
	fmt.Printf("socket:%v\n", socket)
	pty, slavePath, err := console.NewPty()
	err = pty.Resize(console.WinSize{
		Height: 30,
		Width:  60,
	})
	defer pty.Close()
	fmt.Printf("ptyname:%s,pty:%v,slavepath:%s\n", pty.Name(), pty.Fd(), slavePath)
	if err := utils.SendFd(socket, pty.Name(), pty.Fd()); err != nil {
		fmt.Printf("sendfd err:%s\n", err)
		return err
	}
	err = dupStdio(slavePath)
	_, err = syscall.Setsid()
	if err := unix.IoctlSetInt(0, unix.TIOCSCTTY, 0); err != nil {
		return err
	}
	return err
}

func execshell() error {
	cmd := exec.Command("/bin/sh")
	//cmd := exec.Command("/bin/ls", "-l")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	fmt.Printf("execshell is over\n")
	return err
}

func NewContainerInit() error {
	//envConsole := os.Getenv("_LIBCONTAINER_INITPIPE")
	envConsole := os.Getenv("_LIBCONTAINER_CONSOLE")
	console, err := strconv.Atoi(envConsole)
	fmt.Printf("receive console:%d\n", console)
	consoleSocket := os.NewFile(uintptr(console), "console")
	var config *InitConfig
	err = json.NewDecoder(consoleSocket).Decode(&config)
	fmt.Printf("receive config:%+v\n", config)
	err = setconsole(consoleSocket)
	err = execshell()
	if err != nil {
		fmt.Printf("execshell err:%s\n", err)
	}
	select {}
	return err
}
