package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"hyphon/sampcon/libcontainer/config"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/urfave/cli"
	"golang.org/x/sys/unix"
)

func startContainer(context *cli.Context, spec *specs.Spec) (int, error) {
	id := context.Args().First()
	if id == "" {
		errEmptyID := errors.New("container id can not be empty.")
		return -1, errEmptyID
	}
	err := createContainer(context, id, spec)

	//r := &runner{}
	//return r.run()
	return 0, err
}

func createContainer(context *cli.Context, id string, spec *specs.Spec) error {
	config, err := specconv(context, spec)
	fmt.Printf("config:%+v\n", config)
	return err
}

// parseMountOptions parses the string and returns the flags, propagation
// flags and any mount data that it contains.
func parseMountOptions(options []string) (int, []int, string, int) {
	var (
		flag     int
		pgflag   []int
		data     []string
		extFlags int
	)
	flags := map[string]struct {
		clear bool
		flag  int
	}{
		"acl":           {false, unix.MS_POSIXACL},
		"async":         {true, unix.MS_SYNCHRONOUS},
		"atime":         {true, unix.MS_NOATIME},
		"bind":          {false, unix.MS_BIND},
		"defaults":      {false, 0},
		"dev":           {true, unix.MS_NODEV},
		"diratime":      {true, unix.MS_NODIRATIME},
		"dirsync":       {false, unix.MS_DIRSYNC},
		"exec":          {true, unix.MS_NOEXEC},
		"iversion":      {false, unix.MS_I_VERSION},
		"lazytime":      {false, unix.MS_LAZYTIME},
		"loud":          {true, unix.MS_SILENT},
		"mand":          {false, unix.MS_MANDLOCK},
		"noacl":         {true, unix.MS_POSIXACL},
		"noatime":       {false, unix.MS_NOATIME},
		"nodev":         {false, unix.MS_NODEV},
		"nodiratime":    {false, unix.MS_NODIRATIME},
		"noexec":        {false, unix.MS_NOEXEC},
		"noiversion":    {true, unix.MS_I_VERSION},
		"nolazytime":    {true, unix.MS_LAZYTIME},
		"nomand":        {true, unix.MS_MANDLOCK},
		"norelatime":    {true, unix.MS_RELATIME},
		"nostrictatime": {true, unix.MS_STRICTATIME},
		"nosuid":        {false, unix.MS_NOSUID},
		"rbind":         {false, unix.MS_BIND | unix.MS_REC},
		"relatime":      {false, unix.MS_RELATIME},
		"remount":       {false, unix.MS_REMOUNT},
		"ro":            {false, unix.MS_RDONLY},
		"rw":            {true, unix.MS_RDONLY},
		"silent":        {false, unix.MS_SILENT},
		"strictatime":   {false, unix.MS_STRICTATIME},
		"suid":          {true, unix.MS_NOSUID},
		"sync":          {false, unix.MS_SYNCHRONOUS},
	}
	propagationFlags := map[string]int{
		"private":     unix.MS_PRIVATE,
		"shared":      unix.MS_SHARED,
		"slave":       unix.MS_SLAVE,
		"unbindable":  unix.MS_UNBINDABLE,
		"rprivate":    unix.MS_PRIVATE | unix.MS_REC,
		"rshared":     unix.MS_SHARED | unix.MS_REC,
		"rslave":      unix.MS_SLAVE | unix.MS_REC,
		"runbindable": unix.MS_UNBINDABLE | unix.MS_REC,
	}
	extensionFlags := map[string]struct {
		clear bool
		flag  int
	}{
		"tmpcopyup": {false, config.EXT_COPYUP},
	}
	for _, o := range options {
		// If the option does not exist in the flags table or the flag
		// is not supported on the platform,
		// then it is a data value for a specific fs type
		if f, exists := flags[o]; exists && f.flag != 0 {
			if f.clear {
				flag &= ^f.flag
			} else {
				flag |= f.flag
			}
		} else if f, exists := propagationFlags[o]; exists && f != 0 {
			pgflag = append(pgflag, f)
		} else if f, exists := extensionFlags[o]; exists && f.flag != 0 {
			if f.clear {
				extFlags &= ^f.flag
			} else {
				extFlags |= f.flag
			}
		} else {
			data = append(data, o)
		}
	}
	return flag, pgflag, strings.Join(data, ","), extFlags
}

func createmount(path string, m *specs.Mount) *config.Mount {
	flags, pgflags, data, ext := parseMountOptions(m.Options)
	source := m.Source
	device := m.Type
	if flags|unix.MS_BIND != 0 {
		if device == "" {
			device = "bind"
		}
		if !filepath.IsAbs(source) {
			source = filepath.Join(path, m.Source)
		}
	}
	mount := config.Mount{
		Device:           device,
		Source:           source,
		Destination:      m.Destination,
		Data:             data,
		Flags:            flags,
		PropagationFlags: pgflags,
		Extensions:       ext,
	}
	return &mount
}

func specconv(context *cli.Context, spec *specs.Spec) (*config.Config, error) {
	dir, _ := os.Getwd()
	abs_dir, err := filepath.Abs(dir)
	rootfsPath := filepath.Join(abs_dir, spec.Root.Path)
	config := &config.Config{
		Rootfs:     rootfsPath,
		Readonlyfs: spec.Root.Readonly,
		Hostname:   spec.Hostname,
	}
	for _, m := range spec.Mounts {
		config.Mounts = append(config.Mounts, createmount(abs_dir, &m))
	}
	return config, err
}
