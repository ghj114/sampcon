package linuxcontainer

import (
	"fmt"
	"hyphon/sampcon/libcontainer/cgroup"
	"hyphon/sampcon/libcontainer/utils"
	"io"
	"os"
	"os/exec"
)

type initProcess struct {
	cmd        *exec.Cmd
	parentPipe *os.File
	childPipe  *os.File
	//config     *initConfig
	cgroupManager cgroup.Manager
	//intelRdtManager intelrdt.Manager
	container     *Linuxcontainer
	fds           []string
	process       *Process
	bootstrapData io.Reader
	sharePidns    bool
}

func (p *initProcess) wait() (*os.ProcessState, error) {
	err := p.cmd.Wait()
	return p.cmd.ProcessState, err
}

func (p *initProcess) pid() int {
	return p.cmd.Process.Pid
}

func (p *initProcess) start() error {
	defer p.parentPipe.Close()
	fmt.Printf("cmd:%+v\n", p.cmd)
	err := p.cmd.Start()
	p.childPipe.Close()
	p.cgroupManager.Apply(p.pid())

	cfg := &InitConfig{
		//Config:           r.container.config,
		//Args: "abc",
		//Env:              p.Env,
		//User:             p.User,
		//AdditionalGroups: p.AdditionalGroups,
		Cwd: "xyz",
		//Capabilities:     process.Capabilities,
		//PassedFilesCount: len(p.ExtraFiles),
		//ContainerId:      r.container.config.ID(),
		//NoNewPrivileges:  r.container.config.NoNewPrivileges,
		//Rootless:         r.container.config.Rootless,
		//AppArmorProfile:  r.container.config.AppArmorProfile,
		//ProcessLabel:     r.container.config.ProcessLabel,
		//Rlimits:          c.config.Rlimits,
	}
	//err = utils.WriteJSON(p.parentPipe, cfg)
	err = utils.WriteJSON(p.process.PConsoleSocket, cfg)

	go func() {
		state, _ := p.wait()
		fmt.Printf("state:%+v\n", state)
	}()
	//return err

	//cmd := exec.Command("/bin/sh", "-c", "ls -l")
	//cmd := exec.Command("touch a.out")
	//cmd := exec.Command("/proc/self/exe", "init")
	//cmd := exec.Command("ls", "-l")
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//fmt.Printf("cmd:%+v\n", cmd)
	//cmd.Start()
	//err = cmd.Wait()
	//state = cmd.ProcessState
	//fmt.Printf("state:%v\n", state)
	return err
}
