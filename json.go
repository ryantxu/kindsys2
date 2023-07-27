package kindsys2

import (
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type JSONResourceBuilder struct {
	SetStaticMetadata func(v StaticMetadata)
	SetCommonMetadata func(v CommonMetadata)
	SetAnnotation     func(key string, val string)

	ReadSpec   func(iter *jsoniter.Iterator) error
	ReadStatus func(iter *jsoniter.Iterator) error
	ReadSub    func(name string, iter *jsoniter.Iterator) error
}

func ReadResourceJSON(reader io.Reader, builder JSONResourceBuilder) error {
	iter := jsoniter.Parse(jsoniter.ConfigDefault, reader, 1024)

	static := StaticMetadata{}
	common := CommonMetadata{}

	for l1Field := iter.ReadObject(); l1Field != ""; l1Field = iter.ReadObject() {
		err := iter.Error
		switch l1Field {
		case "apiVersion":
			vals := strings.SplitN(iter.ReadString(), "/", 2)
			static.Group = vals[0]
			static.Version = vals[1]

		case "kind":
			static.Kind = iter.ReadString()

		case "metadata":
			for l2Field := iter.ReadObject(); l2Field != ""; l2Field = iter.ReadObject() {
				switch l2Field {
				case "namespace":
					static.Namespace = iter.ReadString()
				case "name":
					static.Name = iter.ReadString()
				case "annotations":
					for anno := iter.ReadObject(); anno != ""; anno = iter.ReadObject() {
						val := iter.ReadString()
						switch anno {
						default:
							fmt.Printf("anno> %s = %v\n", anno, val)
						}
						if iter.Error != nil {
							return iter.Error
						}
					}
				default:
					tt := iter.ReadAny()
					fmt.Printf("meta> %s = %v\n", l2Field, tt)
				}
				if iter.Error != nil {
					return iter.Error
				}
			}
		case "spec":
			err = builder.ReadSpec(iter)
		case "status":
			err = builder.ReadStatus(iter)
		default:
			err = builder.ReadSub(l1Field, iter)
		}
		if err != nil {
			return err
		}
		if iter.Error != nil {
			return iter.Error
		}
	}

	builder.SetStaticMetadata(static)
	builder.SetCommonMetadata(common)
	return iter.Error
}

func WriteResourceJSON(obj Resource, stream *jsoniter.Stream) {
	isMore := false
	static := obj.StaticMetadata()
	common := obj.CommonMetadata()
	custom := obj.CustomMetadata() // ends up in annotations
	spec := obj.SpecObject()

	stream.WriteObjectStart()
	stream.WriteObjectField("apiVersion")
	stream.WriteString(static.Group + "/" + static.Version)
	stream.WriteMore()
	stream.WriteObjectField("kind")
	stream.WriteString(static.Kind)

	stream.WriteMore()
	stream.WriteObjectField("metadata")
	stream.WriteObjectStart()
	isMore = writeOptionalString(false, "name", static.Name, stream)
	isMore = writeOptionalString(isMore, "namespace", static.Namespace, stream)

	if isMore {
		stream.WriteMore()
	}
	stream.WriteObjectField("annotations")
	stream.WriteObjectStart()
	isMore = writeOptionalString(false, "grafana.com/createdBy", common.CreatedBy, stream)

	if custom != nil {
		for k, v := range custom.MapFields() {
			if isMore {
				stream.WriteMore()
			}
			stream.WriteObjectField(k)
			stream.WriteVal(v)
			isMore = true
		}
	}

	stream.WriteObjectEnd()
	isMore = writeOptionalString(false, "resourceVersion", common.ResourceVersion, stream)
	if len(common.Labels) > 0 {
		if isMore {
			stream.WriteMore()
		}
		stream.WriteObjectField("labels")
		stream.WriteVal(common.Labels)
		isMore = true
	}
	_ = writeOptionalString(isMore, "uid", common.UID, stream)
	stream.WriteObjectEnd()

	if spec != nil {
		stream.WriteMore()
		stream.WriteObjectField("spec")
		stream.WriteVal(spec)
	}

	stream.WriteObjectEnd()
}

func writeOptionalString(isMore bool, key string, val string, stream *jsoniter.Stream) bool {
	if val == "" {
		return isMore
	}
	if isMore {
		stream.WriteMore()
	}
	stream.WriteObjectField(key)
	stream.WriteString(val)
	return true
}
