package thema

import (
	"bytes"
	"fmt"
	"io"
	"kindsys2"

	"github.com/grafana/kindsys"
	"github.com/grafana/kindsys/encoding"
	"github.com/grafana/thema"
)

var _ kindsys2.ResourceKind = &ThemaCoreKind{}

type ThemaCoreKind struct {
	kind kindsys.Core
}

// Load a jsonschema based kind from a file system
// the file system will have a manifest that exists
func NewThemaCoreKind(rt *thema.Runtime, def kindsys.Def[kindsys.CoreProperties], opts ...thema.BindOption) (*ThemaCoreKind, error) {
	k, err := kindsys.BindCore(rt, def)
	if err != nil {
		return nil, err
	}
	return &ThemaCoreKind{kind: k}, nil
}

func (m *ThemaCoreKind) CoreKind() kindsys.Core {
	return m.kind
}

func (m *ThemaCoreKind) GetMachineNames() kindsys2.MachineNames {
	p := m.kind.Props()
	c := p.Common()
	return kindsys2.MachineNames{
		Plural:   c.PluralName,
		Singular: c.MachineName,
	}
}

func (m *ThemaCoreKind) GetKindInfo() kindsys2.KindInfo {
	p := m.kind.Props()
	c := p.Common()
	return kindsys2.KindInfo{
		Group:       m.kind.Group(),
		Kind:        c.Name,
		Description: c.Description,
	}
}

func (m *ThemaCoreKind) CurrentVersion() string {
	return m.kind.CurrentVersion().String()
}

func (m *ThemaCoreKind) GetVersions() []kindsys2.VersionInfo {
	versions := []kindsys2.VersionInfo{}
	for _, schema := range m.kind.Lineage().All() {
		versions = append(versions, kindsys2.VersionInfo{
			Version: schema.Version().String(),
		})
	}
	return versions
}

func (m *ThemaCoreKind) GetJSONSchema(version string) (string, error) {
	return "", fmt.Errorf("TODO")
}

func (m *ThemaCoreKind) Read(reader io.Reader, strict bool) (kindsys2.Resource, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	res, err := m.kind.FromBytes(buf.Bytes(), &encoding.KubernetesJSONDecoder{})
	if err != nil {
		return nil, err
	}
	fmt.Printf("GOT: %v", res)
	// TODO!!!!
	obj := &kindsys2.UnstructuredResource{}
	obj.Spec = res.Spec
	//	obj.CommonMeta = res.CommonMetadata()
	return obj, nil
}

func (m *ThemaCoreKind) Migrate(obj kindsys2.Resource, targetVersion string) (kindsys2.Resource, error) {
	return nil, fmt.Errorf("TODO")
}
