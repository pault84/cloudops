package cloudops

import (
	"errors"
	"fmt"
)

const (
	// SetIdentifierNone is a default identifier to group all disks from a
	// particular set
	SetIdentifierNone = "None"
)

// Custom storage operation error codes.
const (
	_ = iota + 5000
	// ErrVolDetached is code for a volume is detached on the instance
	ErrVolDetached
	// ErrVolInval is the code for a invalid volume
	ErrVolInval
	// ErrVolAttachedOnRemoteNode is code when a volume is not attached locally
	// but attached on a remote node
	ErrVolAttachedOnRemoteNode
	// ErrVolNotFound is code when a volume is not found
	ErrVolNotFound
	// ErrInvalidDevicePath is code when a volume/disk has invalid device path
	ErrInvalidDevicePath
)

var (
	// ErrUnsupported operation is unsupported.
	ErrUnsupported = errors.New("Unsupported Operation")
)

// ErrNotSupported is returned when a particular operation is not supported
var ErrNotSupported = fmt.Errorf("operation not supported")

// StorageError error returned for storage operations
type StorageError struct {
	// Code is one of storage operation driver error codes.
	Code int
	// Msg is human understandable error message.
	Msg string
	// Instance provides more information on the error.
	Instance string
}

// CloudResourceInfo provides metadata information on a cloud resource.
type CloudResourceInfo struct {
	// Name of the cloud resource.
	Name string
	// ID of the cloud resource.
	ID string
	// Labels on the cloud resource.
	Labels map[string]string
	// Zone where the resource exists.
	Zone string
	// Region where the resource exists.
	Region string
}

// InstanceGroupInfo encapsulates info for a cloud instance group. In AWS this
// maps to ASG.
type InstanceGroupInfo struct {
	CloudResourceInfo
	// AutoscalingEnabled is true if auto scaling is turned on
	AutoscalingEnabled bool
	// Min number of nodes in the instance group
	Min *int64
	// Max number of nodes in the instance group
	Max *int64
	// Zones that the instance group is part of
	Zones []string
}

// InstanceInfo encapsulates info for a cloud instance
type InstanceInfo struct {
	CloudResourceInfo
}

// Compute interface to manage compute instances.
type Compute interface {
	// InstanceID of instance where command is executed.
	InstanceID() string
	// InspectInstance inspects the node with the given instance ID
	InspectInstance(instanceID string) (*InstanceInfo, error)
	// InspectInstanceGroupForInstance inspects the instance group to which the
	// cloud instance with given ID belongs
	InspectInstanceGroupForInstance(instanceID string) (*InstanceGroupInfo, error)
}

// Storage interface to manage storage operations.
type Storage interface {
	// Create volume based on input template volume and also apply given labels.
	Create(template interface{}, labels map[string]string) (interface{}, error)
	// GetDeviceID returns ID/Name of the given device/disk or snapshot
	GetDeviceID(template interface{}) (string, error)
	// Attach volumeID.
	// Return attach path.
	Attach(volumeID string) (string, error)
	// Detach volumeID.
	Detach(volumeID string) error
	// DetachFrom detaches the disk/volume with given ID from the given instance ID
	DetachFrom(volumeID, instanceID string) error
	// Delete volumeID.
	Delete(volumeID string) error
	// DeleteFrom deletes the given volume/disk from the given instanceID
	DeleteFrom(volumeID, instanceID string) error
	// Desribe an instance
	Describe() (interface{}, error)
	// FreeDevices returns free block devices on the instance.
	// blockDeviceMappings is a data structure that contains all block devices on
	// the instance and where they are mapped to
	FreeDevices(blockDeviceMappings []interface{}, rootDeviceName string) ([]string, error)
	// Inspect volumes specified by volumeID
	Inspect(volumeIds []*string) ([]interface{}, error)
	// DeviceMappings returns map[local_attached_volume_path]->volume ID/NAME
	DeviceMappings() (map[string]string, error)
	// Enumerate volumes that match given filters. Organize them into
	// sets identified by setIdentifier.
	// labels can be nil, setIdentifier can be empty string.
	Enumerate(volumeIds []*string,
		labels map[string]string,
		setIdentifier string,
	) (map[string][]interface{}, error)
	// DevicePath for the given volume i.e path where it's attached
	DevicePath(volumeID string) (string, error)
	// Snapshot the volume with given volumeID
	Snapshot(volumeID string, readonly bool) (interface{}, error)
	// SnapshotDelete deletes the snapshot with given ID
	SnapshotDelete(snapID string) error
	// ApplyTags will apply given labels/tags on the given volume
	ApplyTags(volumeID string, labels map[string]string) error
	// RemoveTags removes labels/tags from the given volume
	RemoveTags(volumeID string, labels map[string]string) error
	// Tags will list the existing labels/tags on the given volume
	Tags(volumeID string) (map[string]string, error)
}

// Ops interface to perform basic cloud operations.
type Ops interface {
	// Name returns name of the cloud operations driver
	Name() string
	// Storage operations in the cloud
	Storage
	// Compute operations in the cloud
	Compute
}

// NewStorageError creates a new custom storage error instance
func NewStorageError(code int, msg string, instance string) error {
	return &StorageError{Code: code, Msg: msg, Instance: instance}
}

func (e *StorageError) Error() string {
	return e.Msg
}
