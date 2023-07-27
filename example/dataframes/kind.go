package dataframes

import (
	"embed"
	"fmt"
	"io"
	"kindsys2"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	jsoniter "github.com/json-iterator/go"
)

var _ kindsys2.Resource = &DataFramesResource{}
var _ kindsys2.ResourceKind = &DataFramesResourceKind{}

// DataFramesResource holds a collection of DataFrames
type DataFramesResource = kindsys2.GenericResource[
	data.Frames, // << the frame array
	kindsys2.SimpleCustomMetadata,
	map[string]any, // any status
]

type DataFramesResourceKind struct{}

//go:embed schema-*.json
var dir embed.FS

func (k *DataFramesResourceKind) GetKindInfo() kindsys2.KindInfo {
	return kindsys2.KindInfo{
		Group:       "ext.dataframes.grafana.com",
		Kind:        "DataFrames",
		Description: "A collection of DataFrames",
		Maturity:    kindsys2.MaturityExperimental,
	}
}

var currentVersion = kindsys2.VersionInfo{
	Version: "v1-0",
}

func (k *DataFramesResourceKind) CurrentVersion() string {
	return currentVersion.Version
}

func (k *DataFramesResourceKind) GetVersions() []kindsys2.VersionInfo {
	return []kindsys2.VersionInfo{currentVersion}
}

func (k *DataFramesResourceKind) GetJSONSchema(version string) (string, error) {
	if version != currentVersion.Version {
		return "", fmt.Errorf("unknown version")
	}
	s, err := dir.ReadFile("schema-v1.json")
	if err != nil {
		return "", err
	}
	return string(s), nil
}

func (k *DataFramesResourceKind) GetMachineNames() kindsys2.MachineNames {
	return kindsys2.MachineNames{
		Plural:   "dataframes", // already plural
		Singular: "dataframes", // the same i suppose
	}
}

func (k *DataFramesResourceKind) Read(reader io.Reader, strict bool) (kindsys2.Resource, error) {
	return k.ReadFrames(reader, strict) // From well typed to
}

func (k *DataFramesResourceKind) ReadFrames(reader io.Reader, strict bool) (*DataFramesResource, error) {
	obj := &DataFramesResource{}
	err := kindsys2.ReadResourceJSON(reader, kindsys2.JSONResourceBuilder{
		SetStaticMetadata: func(v kindsys2.StaticMetadata) { obj.StaticMeta = v },
		SetCommonMetadata: func(v kindsys2.CommonMetadata) { obj.CommonMeta = v },
		ReadSpec: func(iter *jsoniter.Iterator) error {
			obj.Spec = data.Frames{}
			for iter.ReadArray() {
				frame := &data.Frame{}
				iter.ReadVal(frame)
				if iter.Error != nil {
					return iter.Error
				}
				obj.Spec = append(obj.Spec, frame)
			}
			return nil
		},
		SetAnnotation: func(key, val string) {
			fmt.Printf("??? unknown")
		},
		ReadStatus: func(iter *jsoniter.Iterator) error {
			obj.Status = make(map[string]any)
			iter.ReadVal(obj.Status)
			return iter.Error
		},
		ReadSub: func(name string, iter *jsoniter.Iterator) error {
			return fmt.Errorf("unsupported sub resource")
		},
	})

	if strict {
		// TODO?? check data plane contracts??

		// Check that the lengths are all OK
		for idx, frame := range obj.Spec {
			_, err := frame.RowLen()
			if err != nil {
				return obj, fmt.Errorf("invalid rows on frame[%d] %w", idx, err)
			}
		}
	}

	return obj, err
}

func (k *DataFramesResourceKind) Migrate(obj kindsys2.Resource, targetVersion string) (kindsys2.Resource, error) {
	if targetVersion == "" || obj.StaticMetadata().Version == targetVersion {
		return obj, nil // noop
	}
	return nil, fmt.Errorf("migrations are not yet supported")
}
