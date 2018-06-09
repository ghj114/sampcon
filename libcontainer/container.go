package linuxcontainer

import (
	"os"
	"sync"
	"time"
	"path/filepath"

	"hyphon/sampcon/libcontainer/config"
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
    id                   string
    root                 string
    config               *config.Config
    //cgroupManager        cgroups.Manager
    //intelRdtManager      intelrdt.Manager
    initPath             string
    initArgs             []string
    //initProcess          parentProcess
    initProcessStartTime uint64
    //criuPath             string
    newuidmapPath        string
    newgidmapPath        string
    m                    sync.Mutex
    //criuVersion          int
    //state                containerState
    created              time.Time
}

func New(id string) (linuxcon *Linuxcontainer, error){
	containerRoot := filepath.Join("/run/runc", id)
	linuxcon := &linuxContainer{
		id: id
		root:containerRoot,
		initPath:"/proc/self/exe",
		initArgs:  []string{os.Args[0], "init"},

	}
}
