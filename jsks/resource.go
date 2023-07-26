package jsks

import (
	"fmt"
	"io"
	"io/fs"
	"kindsys2"
)

var _ kindsys2.ResourceKind = &resourceKindFromManifest{}

type resourceKindFromManifest struct {
	kindFromManifest // the base properties

	//
	machineNames kindsys2.MachineNames
}

// Load a jsonschema based kind from a file system
// the file system will have a manifest that exists
func NewResourceKind(sfs fs.FS) (kindsys2.ResourceKind, error) {

	return nil, nil
}

func (m *resourceKindFromManifest) GetMachineNames() kindsys2.MachineNames {
	return m.machineNames
}

func (m *resourceKindFromManifest) Parse(reader io.Reader) (kindsys2.Resource, error) {
	return nil, fmt.Errorf("TODO")
}

func (m *resourceKindFromManifest) Validate(obj kindsys2.Resource) error {
	return fmt.Errorf("TODO")
}

func (m *resourceKindFromManifest) Migrate(obj kindsys2.Resource, targetVersion string) (kindsys2.Resource, error) {
	return nil, fmt.Errorf("TODO")
}
