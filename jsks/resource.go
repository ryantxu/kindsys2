package jsks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"kindsys2"
	"strings"
)

var _ kindsys2.ResourceKind = &resourceKindFromManifest{}

type resourceKindFromManifest struct {
	kindFromManifest // the base properties

	//
	names kindsys2.MachineNames
}

// JUST USED while testing things -- with private package
func HackGetManifestJSON(k kindsys2.ResourceKind) ([]byte, error) {
	names := k.GetMachineNames()
	info := &manifest{
		KindInfo:     k.GetKindInfo(),
		Versions:     k.GetVersions(),
		MachineNames: &names,
	}
	return json.MarshalIndent(info, "", "  ")
}

// Load a jsonschema based kind from a file system
// the file system will have a manifest that exists
func NewResourceKind(sfs fs.FS) (kindsys2.ResourceKind, error) {
	m := &resourceKindFromManifest{}
	info, err := m.init(sfs)
	if err != nil {
		return m, err
	}
	if info.ComposableType != "" || len(info.ComposableSlots) > 0 {
		return nil, fmt.Errorf("invalid info in the manifest (should not have composable types)")
	}

	if info.MachineNames != nil {
		m.names = *info.MachineNames
	}
	if m.names.Singular == "" {
		m.names.Singular = strings.ToLower(info.Kind)
	}
	if m.names.Plural == "" {
		m.names.Plural = m.names.Singular + "s"
	}
	return m, nil
}

func (m *resourceKindFromManifest) GetMachineNames() kindsys2.MachineNames {
	return m.names
}

func (m *resourceKindFromManifest) Parse(reader io.Reader) (kindsys2.Resource, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	obj := &kindsys2.UnstructuredResource{}
	err = json.Unmarshal(buf.Bytes(), obj)
	return obj, err
}

func (m *resourceKindFromManifest) Validate(obj kindsys2.Resource) error {
	meta := obj.StaticMetadata()
	if meta.Group != m.info.Group {
		return fmt.Errorf("wrong group")
	}
	if meta.Kind != m.info.Kind {
		return fmt.Errorf("wrong kind")
	}

	schema, ok := m.parsed[meta.Version]
	if !ok || schema == nil {
		return fmt.Errorf("unknown version")
	}

	// TODO!!! schema is right now just on the spec!!!
	doc := obj.SpecObject()
	// TODO: need to make sure the doc+resource are ones that we can parse ()
	return schema.ValidateInterface(doc)
}

func (m *resourceKindFromManifest) Migrate(obj kindsys2.Resource, targetVersion string) (kindsys2.Resource, error) {
	return nil, fmt.Errorf("TODO")
}
