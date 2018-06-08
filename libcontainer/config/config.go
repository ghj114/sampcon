package config

import "time"

type Command struct {
	Path    string         `json:"path"`
	Args    []string       `json:"args"`
	Env     []string       `json:"env"`
	Dir     string         `json:"dir"`
	Timeout *time.Duration `json:"timeout"`
}

const (
	// EXT_COPYUP is a directive to copy up the contents of a directory when
	// a tmpfs is mounted over it.
	EXT_COPYUP = 1 << iota
)

type Mount struct {
	// Source path for the mount.
	Source string `json:"source"`

	// Destination path for the mount inside the container.
	Destination string `json:"destination"`

	// Device the mount is for.
	Device string `json:"device"`

	// Mount flags.
	Flags int `json:"flags"`

	// Propagation Flags
	PropagationFlags []int `json:"propagation_flags"`

	// Mount data applied to the mount.
	Data string `json:"data"`

	// Relabel source if set, "z" indicates shared, "Z" indicates unshared.
	Relabel string `json:"relabel"`

	// Extensions are additional flags that are specific to runc.
	Extensions int `json:"extensions"`

	// Optional Command to be run before Source is mounted.
	PremountCmds []Command `json:"premount_cmds"`

	// Optional Command to be run after Source is mounted.
	PostmountCmds []Command `json:"postmount_cmds"`
}

type Config struct {
	// Path to a directory containing the container's root filesystem.
	Rootfs string `json:"rootfs"`

	// Readonlyfs will remount the container's rootfs as readonly where only externally mounted
	// bind mounts are writtable.
	Readonlyfs bool `json:"readonlyfs"`

	// Specifies the mount propagation flags to be applied to /.
	//RootPropagation int `json:"rootPropagation"`

	// Mounts specify additional source and destination paths that will be mounted inside the container's
	// rootfs and mount namespace if specified
	Mounts []*Mount `json:"mounts"`

	// The device nodes that should be automatically created within the container upon container start.  Note, make sure that the node is marked as allowed in the cgroup as well!
	//Devices []*Device `json:"devices"`

	//MountLabel string `json:"mount_label"`

	// Hostname optionally sets the container's hostname if provided
	Hostname string `json:"hostname"`

	// Namespaces specifies the container's namespaces that it should setup when cloning the init process
	// If a namespace is not provided that namespace is shared from the container's parent process
	//Namespaces Namespaces `json:"namespaces"`

	// Capabilities specify the capabilities to keep when executing the process inside the container
	// All capabilities not specified will be dropped from the processes capability mask
	//Capabilities *Capabilities `json:"capabilities"`

	// Networks specifies the container's network setup to be created
	//Networks []*Network `json:"networks"`

	// Routes can be specified to create entries in the route table as the container is started
	//Routes []*Route `json:"routes"`

	// Cgroups specifies specific cgroup settings for the various subsystems that the container is
	// placed into to limit the resources the container has available
	//Cgroups *Cgroup `json:"cgroups"`

	// UidMappings is an array of User ID mappings for User Namespaces
	//UidMappings []IDMap `json:"uid_mappings"`

	// GidMappings is an array of Group ID mappings for User Namespaces
	//GidMappings []IDMap `json:"gid_mappings"`
}
