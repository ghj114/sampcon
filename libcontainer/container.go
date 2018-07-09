package linuxcontainer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"hyphon/sampcon/libcontainer/config"
	"hyphon/sampcon/libcontainer/utils"

	"golang.org/x/sys/unix"
)

const (
	execFifoFilename = "exec.fifo"
)

//type Linuxcontainer struct{
//	// Root directory for the factory to store state.
//	Root string
//	// InitPath is the path for calling the init responsibilities for spawning a container.
//	InitPath string
//	// InitArgs are arguments for calling the init responsibilities for spawning a container.
//	InitArgs []string
//	// New{u,g}uidmapPath is the path to the binaries used for mapping with rootless containers.
//	NewuidmapPath string
//	NewgidmapPath string
//	// Validator provides validation to container configurations.
//	//Validator validate.Validator
//	// NewCgroupsManager returns an initialized cgroups manager for a single container.
//	//NewCgroupsManager func(config *config.Cgroup, paths map[string]string) cgroups.Manager
//}

type Linuxcontainer struct {
	id     string
	root   string
	config *config.Config
	//cgroupManager        cgroups.Manager
	//intelRdtManager      intelrdt.Manager
	initPath string
	initArgs []string
	//initProcess          parentProcess
	initProcessStartTime uint64
	//criuPath             string
	newuidmapPath string
	newgidmapPath string
	m             sync.Mutex
	//criuVersion          int
	//state                containerState
	created time.Time
}

func (c *Linuxcontainer) createExecinfo() error {
	//rootuid, err := c.Config().HostRootUID()
	//if err != nil {
	//	return err
	//}
	//rootgid, err := c.Config().HostRootGID()
	//if err != nil {
	//	return err
	//}

	fmt.Printf("fifo:%s, %s\n", c.root, execFifoFilename)
	fifoName := filepath.Join(c.root, execFifoFilename)
	fmt.Printf("fifoname:%s\n", fifoName)
	if _, err := os.Stat(fifoName); err == nil {
		return fmt.Errorf("exec fifo %s already exists", fifoName)
	}
	oldMask := unix.Umask(0000)
	if err := unix.Mkfifo(fifoName, 0622); err != nil {
		unix.Umask(oldMask)
		fmt.Printf("err:%s\n", err)
		return err
	}
	unix.Umask(oldMask)
	if err := os.Chown(fifoName, 0, 0); err != nil {
		return err
	}
	return nil
}
func (c *Linuxcontainer) execCmd(p *Process, childpipe *os.File) (*exec.Cmd, error) {
	fmt.Printf("execCmd in path:%s, args:%v\n", c.initPath, c.initArgs)
	cmd := exec.Command(c.initPath, c.initArgs[1:]...)
	//cmd := exec.Command("/proc/self/exe", "init")
	cmd.Args[0] = c.initArgs[0]
	//cmd.Stdin = p.Stdin
	//cmd.Stdout = p.Stdout
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = c.config.Rootfs
	cmd.ExtraFiles = append(cmd.ExtraFiles, p.ConsoleSocket, childpipe)
	cmd.Env = append(cmd.Env,
		fmt.Sprintf("_LIBCONTAINER_CONSOLE=%d", 3+len(cmd.ExtraFiles)-2),
		fmt.Sprintf("_LIBCONTAINER_INITPIPE=%d", 3+len(cmd.ExtraFiles)-1),
	)
	fmt.Printf("cmd.ENV:%s\n", cmd.Env)
	return cmd, nil
}

func (c *Linuxcontainer) newParentProcess(p *Process) (*initProcess, error) {
	parentPipe, childPipe, err := utils.NewSockPair("init")
	cmd, err := c.execCmd(p, childPipe)
	initProcess := &initProcess{
		cmd:        cmd,
		childPipe:  childPipe,
		parentPipe: parentPipe,
		//manager:         c.cgroupManager,
		//intelRdtManager: c.intelRdtManager,
		//config:          c.newInitConfig(p),
		container: c,
		process:   p,
		//bootstrapData: data,
		//sharePidns:    sharePidns,
	}

	return initProcess, err
}

func (c *Linuxcontainer) start(p *Process) error {
	newparent, err := c.newParentProcess(p)
	err = newparent.start()
	return err
}

func (c *Linuxcontainer) Start(p *Process) error {
	err := c.createExecinfo()
	if err != nil {
		fmt.Printf("err:%s\n", err)
	}
	c.start(p)
	return nil
}

func (c *Linuxcontainer) Run(p *Process) error {
	c.Start(p)
	//c.Exec()
	return nil
}

func New(id string, config *config.Config) (linuxcon *Linuxcontainer, err error) {
	containerRoot := filepath.Join("/run/runc", id)
	if _, err := os.Stat(containerRoot); err == nil {
		return nil, fmt.Errorf("container with id exists: %v", id)
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	if err := os.MkdirAll(containerRoot, 0711); err != nil {
		return nil, err
	}
	if err := os.Chown(containerRoot, unix.Geteuid(), unix.Getegid()); err != nil {
		return nil, err
	}
	linuxcon = &Linuxcontainer{
		id:       id,
		root:     containerRoot,
		config:   config,
		initPath: "/proc/self/exe",
		initArgs: []string{os.Args[0], "init"},
	}
	return linuxcon, nil
}
