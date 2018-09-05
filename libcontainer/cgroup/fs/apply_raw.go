package fs

import (
	"fmt"
	"hyphon/sampcon/libcontainer/cgroup"
	"hyphon/sampcon/libcontainer/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type subsystemSet []subsystem
type subsystem interface {
	// Name returns the name of the subsystem.
	Name() string
	// Returns the stats, as 'stats', corresponding to the cgroup under 'path'.
	//GetStats(path string, stats *cgroup.Stats) error
	// Removes the cgroup represented by 'cgroupData'.
	Remove(*cgroupData) error
	// Creates and joins the cgroup represented by 'cgroupData'.
	Apply(*cgroupData) error
	// Set the cgroup represented by cgroup.
	Set(path string, cgroup *config.Cgroup) error
}

var (
	subsystems = subsystemSet{
		&CpusetGroup{},
		&DevicesGroup{},
		//&MemoryGroup{},
		&CpuGroup{},
		&CpuacctGroup{},
		&PidsGroup{},
		//&BlkioGroup{},
		//&HugetlbGroup{},
		//&NetClsGroup{},
		//&NetPrioGroup{},
		//&PerfEventGroup{},
		//&FreezerGroup{},
		&NameGroup{GroupName: "name=systemd", Join: true},
	}
	//HugePageSizes, _ = cgroups.GetHugePageSize()
)

type Manager struct {
	mu      sync.Mutex
	Cgroups *config.Cgroup
	Paths   map[string]string
}

type cgroupData struct {
	root      string
	innerPath string
	config    *config.Cgroup
	pid       int
}

func (raw *cgroupData) join(subsystem string) (string, error) {
	path, err := raw.path(subsystem)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", err
	}
	if err := cgroup.WriteCgroupProc(path, raw.pid); err != nil {
		return "", err
	}
	return path, nil
}

func (raw *cgroupData) path(subsystem string) (string, error) {
	mnt, err := cgroup.FindCgroupMountpoint(subsystem)
	// If we didn't mount the subsystem, there is no point we make the path.
	if err != nil {
		return "", err
	}

	// If the cgroup name/path is absolute do not look relative to the cgroup of the init process.
	if filepath.IsAbs(raw.innerPath) {
		// Sometimes subsystems can be mounted together as 'cpu,cpuacct'.
		return filepath.Join(raw.root, filepath.Base(mnt), raw.innerPath), nil
	}

	// Use GetOwnCgroupPath instead of GetInitCgroupPath, because the creating
	// process could in container and shared pid namespace with host, and
	// /proc/1/cgroup could point to whole other world of cgroups.
	parentPath, err := cgroup.GetOwnCgroupPath(subsystem)
	if err != nil {
		return "", err
	}

	return filepath.Join(parentPath, raw.innerPath), nil
}

func getCgroupData(c *config.Cgroup, pid int) (*cgroupData, error) {
	//root, err := getCgroupRoot()
	//if err != nil {
	//	return nil, err
	//}

	if (c.Name != "" || c.Parent != "") && c.Path != "" {
		return nil, fmt.Errorf("cgroup: either Path or Name and Parent should be used")
	}

	// XXX: Do not remove this code. Path safety is important! -- cyphar
	cgPath := cgroup.CleanPath(c.Path)
	cgParent := cgroup.CleanPath(c.Parent)
	cgName := cgroup.CleanPath(c.Name)

	innerPath := cgPath
	if innerPath == "" {
		innerPath = filepath.Join(cgParent, cgName)
	}

	return &cgroupData{
		root:      "/sys/fs/cgroup/",
		innerPath: innerPath,
		config:    c,
		pid:       pid,
	}, nil
}

func (m *Manager) Apply(pid int) (err error) {
	if m.Cgroups == nil {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	//var c = m.Cgroups

	d, err := getCgroupData(m.Cgroups, pid)
	fmt.Printf("hyph.cgroup:%+v\n", d)
	if err != nil {
		return err
	}

	m.Paths = make(map[string]string)
	//if c.Paths != nil {
	//	for name, path := range c.Paths {
	//		_, err := d.path(name)
	//		if err != nil {
	//			if cgroups.IsNotFound(err) {
	//				continue
	//			}
	//			return err
	//		}
	//		m.Paths[name] = path
	//	}
	//	return cgroups.EnterPid(m.Paths, pid)
	//}

	for _, sys := range subsystems {
		// TODO: Apply should, ideally, be reentrant or be broken up into a separate
		// create and join phase so that the cgroup hierarchy for a container can be
		// created then join consists of writing the process pids to cgroup.procs
		p, err := d.path(sys.Name())
		//fmt.Printf("hyph.path=%s, sysName=%s\n", p, sys.Name())
		if err != nil {
			// The non-presence of the devices subsystem is
			// considered fatal for security reasons.
			if cgroup.IsNotFound(err) && sys.Name() != "devices" {
				continue
			}
			return err
		}
		m.Paths[sys.Name()] = p
		fmt.Printf("hyph.cgroup:name=%s,path=%s\n", sys.Name(), p)

		if err := sys.Apply(d); err != nil {
			if os.IsPermission(err) && m.Cgroups.Path == "" {
				// If we didn't set a cgroup path, then let's defer the error here
				// until we know whether we have set limits or not.
				// If we hadn't set limits, then it's ok that we couldn't join this cgroup, because
				// it will have the same limits as its parent.
				delete(m.Paths, sys.Name())
				continue
			}
			return err
		}

	}
	return nil
}

func writeFile(dir, file, data string) error {
	// Normally dir should not be empty, one case is that cgroup subsystem
	// is not mounted, we will get empty dir, and we want it fail here.
	if dir == "" {
		return fmt.Errorf("no such directory for %s", file)
	}
	if err := ioutil.WriteFile(filepath.Join(dir, file), []byte(data), 0700); err != nil {
		return fmt.Errorf("failed to write %v to %v: %v", data, file, err)
	}
	return nil
}

func readFile(dir, file string) (string, error) {
	data, err := ioutil.ReadFile(filepath.Join(dir, file))
	return string(data), err
}

func removePath(p string, err error) error {
	if err != nil {
		return err
	}
	if p != "" {
		return os.RemoveAll(p)
	}
	return nil
}
