package jsks

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"kindsys2"
	"strings"

	jsoniter "github.com/json-iterator/go"
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

func (m *resourceKindFromManifest) Read(reader io.Reader, strict bool) (kindsys2.Resource, error) {
	obj := &kindsys2.UnstructuredResource{}
	err := kindsys2.ReadResourceJSON(reader, kindsys2.JSONResourceBuilder{
		SetStaticMetadata: func(v kindsys2.StaticMetadata) { obj.StaticMeta = v },
		SetCommonMetadata: func(v kindsys2.CommonMetadata) { obj.CommonMeta = v },
		ReadSpec: func(iter *jsoniter.Iterator) error {
			obj.Spec = make(map[string]any)
			iter.ReadVal(&obj.Spec)
			return iter.Error
		},
		SetAnnotation: func(key, val string) {
			fmt.Printf("??? unknown")
		},
		ReadStatus: func(iter *jsoniter.Iterator) error {
			obj.Status = make(map[string]any)
			iter.ReadVal(&obj.Status)
			return iter.Error
		},
		ReadSub: func(name string, iter *jsoniter.Iterator) error {
			return fmt.Errorf("unsupported sub resource")
		},
	})
	if err != nil {
		return obj, err
	}

	if strict {
		meta := obj.StaticMetadata()
		if meta.Group != m.info.Group {
			return obj, fmt.Errorf("wrong group")
		}
		if meta.Kind != m.info.Kind {
			return obj, fmt.Errorf("wrong kind")
		}

		schema, ok := m.parsed[meta.Version]
		if !ok || schema == nil {
			return obj, fmt.Errorf("unknown version")
		}

		// TODO!!! schema is right now just on the spec!!!
		doc := obj.SpecObject()
		// TODO: need to make sure the doc+resource are ones that we can parse ()
		err = schema.ValidateInterface(doc)
	}
	return obj, err
}

func (m *resourceKindFromManifest) Migrate(obj kindsys2.Resource, targetVersion string) (kindsys2.Resource, error) {
	return nil, fmt.Errorf("TODO")
}
