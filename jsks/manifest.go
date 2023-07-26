package jsks

import (
	"fmt"
	"io/fs"
	"kindsys2"
)

// Internal type that describes what kind of thing we are looking at
type manifest struct {
	// Only valid for composable types
	Type string `json:"type,omitempty"`

	kindsys2.KindInfo

	// Only valid for resource types
	MachineNames *kindsys2.MachineNames `json:"machineName,omitempty"`

	// Only valid for composable types
	ComposableType string `json:"composableType,omitempty"`

	// Only valid for composable types
	// ??? do we want/need multiple slots?  should each slot be a different type?
	ComposableSlots []string `json:"slots,omitempty"`

	// List of version info
	Versions []kindsys2.VersionInfo `json:"versions"`
}

var _ kindsys2.Kind = &kindFromManifest{}

type kindFromManifest struct {
	info     kindsys2.KindInfo
	current  kindsys2.VersionInfo
	versions []kindsys2.VersionInfo
	raw      map[string]string // raw
}

// Load all the schemas
func (m *kindFromManifest) init(sfs fs.FS) (manifest, error) {
	return nil
}

// func init(sfs fs.FS) (*KindFromManifest, error) {
// 	manifest, err := sfs.Open("kind.json")
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to find kind manifest")
// 	}

// 	fmt.Printf("xx:%v", manifest)

// 	m := &kindFromManifest{
// 		raw: make(map[string]string),
// 	}

// 	// sch, err := jsonschema.Compile("testdata/person_schema.json")
// 	// if err != nil {
// 	// 	log.Fatalf("%#v", err)
// 	// }

// 	return m, nil
// }

func (m *kindFromManifest) GetKindInfo() kindsys2.KindInfo {
	return m.info
}

func (m *kindFromManifest) CurrentVersion() string {
	return m.current.Version
}

func (m *kindFromManifest) GetVersions() []kindsys2.VersionInfo {
	return m.versions
}

func (m *kindFromManifest) GetJSONSchema(version string) (string, error) {
	s, ok := m.raw[version]
	if !ok {
		return "", fmt.Errorf("unknown version")
	}
	return s, nil
}
