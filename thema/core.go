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
func NewThemaCoreKind(rt *thema.Runtime, def kindsys.Def[kindsys.CoreProperties]) (*ThemaCoreKind, error) {
	k, err := kindsys.BindCore(rt, def)
	if err != nil {
		return nil, err
	}
	return &ThemaCoreKind{kind: k}, nil
}

func (k *ThemaCoreKind) CoreKind() kindsys.Core {
	return k.kind
}

func (k *ThemaCoreKind) GetMachineNames() kindsys2.MachineNames {
	p := k.kind.Props()
	c := p.Common()
	return kindsys2.MachineNames{
		Plural:   c.PluralName,
		Singular: c.MachineName,
	}
}

func (k *ThemaCoreKind) GetKindInfo() kindsys2.KindInfo {
	p := k.kind.Props()
	c := p.Common()
	return kindsys2.KindInfo{
		Group:       k.kind.Group(),
		Kind:        c.Name,
		Description: c.Description,
	}
}

func (k *ThemaCoreKind) CurrentVersion() string {
	return k.kind.CurrentVersion().String()
}

func (k *ThemaCoreKind) GetVersions() []kindsys2.VersionInfo {
	versions := []kindsys2.VersionInfo{}
	for _, schema := range k.kind.Lineage().All() {
		versions = append(versions, kindsys2.VersionInfo{
			Version: schema.Version().String(),
		})
	}
	return versions
}

func (k *ThemaCoreKind) GetJSONSchema(version string) (string, error) {
	for _, schema := range k.kind.Lineage().All() {
		if version == schema.Version().String() {
			return "", fmt.Errorf("TODO... convert to JSONSchema")
		}
	}
	return "", fmt.Errorf("unknown version")
}

func (k *ThemaCoreKind) Read(reader io.Reader, strict bool) (kindsys2.Resource, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)

	if strict {
		// ?? is this necessary, or part of the FromBytes below?
		err := k.kind.Validate(buf.Bytes(), &encoding.KubernetesJSONDecoder{})
		if err != nil {
			return nil, err
		}
	}

	res, err := k.kind.FromBytes(buf.Bytes(), &encoding.KubernetesJSONDecoder{})
	if err != nil {
		return nil, err
	}

	// TODO!!!! obviously this should be the same base interfaces
	obj := &kindsys2.UnstructuredResource{}
	obj.SetStaticMetadata(kindsys2.StaticMetadata{
		Group:     res.StaticMeta.Group,
		Kind:      res.StaticMeta.Kind,
		Version:   res.StaticMeta.Version,
		Namespace: res.StaticMeta.Namespace,
		Name:      res.StaticMeta.Name,
	})
	obj.SetCommonMetadata(kindsys2.CommonMetadata{
		UID:               res.CommonMeta.UID,
		ResourceVersion:   res.CommonMeta.ResourceVersion,
		Labels:            res.CommonMeta.Labels,
		CreationTimestamp: res.CommonMeta.CreationTimestamp,
		DeletionTimestamp: res.CommonMeta.DeletionTimestamp,
		Finalizers:        res.CommonMeta.Finalizers,
		UpdateTimestamp:   res.CommonMeta.UpdateTimestamp,
		CreatedBy:         res.CommonMeta.CreatedBy,
		UpdatedBy:         res.CommonMeta.UpdatedBy,
		ExtraFields:       res.CommonMeta.ExtraFields,
	})
	obj.Spec = res.Spec
	return obj, nil
}

func (k *ThemaCoreKind) Migrate(obj kindsys2.Resource, targetVersion string) (kindsys2.Resource, error) {
	return nil, fmt.Errorf("TODO")
}
