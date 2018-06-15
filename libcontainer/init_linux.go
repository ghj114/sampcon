package linuxcontainer

import (
	"encoding/json"
	"fmt"
	"hyphon/sampcon/libcontainer/config"
	"os"
	"strconv"
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

func setconsole() error {
	return nil
}
func NewContainerInit() error {
	fmt.Printf("itest------------1\n")
	envConsole := os.Getenv("_LIBCONTAINER_INITPIPE")
	console, err := strconv.Atoi(envConsole)
	fmt.Printf("receive console:%d\n", console)
	consoleSocket := os.NewFile(uintptr(console), "init")
	var config *InitConfig
	fmt.Printf("itest------------2\n")
	err = json.NewDecoder(consoleSocket).Decode(&config)
	fmt.Printf("receive config:%+v\n", config)
	err = setconsole()
	return err
}
