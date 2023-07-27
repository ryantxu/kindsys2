package kindsys2

import (
	"strings"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

var _ Resource = &UnstructuredResource{}

// UnstructuredResource is an untyped representation of [Resource].
type UnstructuredResource struct {
	BasicMetadataObject
	Spec   map[string]any `json:"spec,omitempty"`
	Status map[string]any `json:"status,omitempty"`
}

func (u *UnstructuredResource) SpecObject() any {
	return u.Spec
}

func (u *UnstructuredResource) Subresources() map[string]any {
	return map[string]any{
		"status": u.Status,
	}
}

func (u *UnstructuredResource) Copy() Resource {
	com := CommonMetadata{
		UID:               u.CommonMeta.UID,
		ResourceVersion:   u.CommonMeta.ResourceVersion,
		CreationTimestamp: u.CommonMeta.CreationTimestamp.UTC(),
		UpdateTimestamp:   u.CommonMeta.UpdateTimestamp.UTC(),
		CreatedBy:         u.CommonMeta.CreatedBy,
		UpdatedBy:         u.CommonMeta.UpdatedBy,
	}

	copy(u.CommonMeta.Finalizers, com.Finalizers)
	if u.CommonMeta.DeletionTimestamp != nil {
		*com.DeletionTimestamp = *(u.CommonMeta.DeletionTimestamp)
	}
	for k, v := range u.CommonMeta.Labels {
		com.Labels[k] = v
	}
	com.ExtraFields = mapcopy(u.CommonMeta.ExtraFields)

	cp := UnstructuredResource{
		Spec:   mapcopy(u.Spec),
		Status: mapcopy(u.Status),
	}

	cp.CommonMeta = com
	cp.CustomMeta = mapcopy(u.CustomMeta)
	return &cp
}

func mapcopy(m map[string]any) map[string]any {
	cp := make(map[string]any)
	for k, v := range m {
		if vm, ok := v.(map[string]any); ok {
			cp[k] = mapcopy(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}

// UnmarshalJSON allows unmarshalling Frame from JSON.
func (u *UnstructuredResource) UnmarshalJSON(b []byte) error {
	iter := jsoniter.ParseBytes(jsoniter.ConfigDefault, b)
	return readResourceJSON(u, iter)
}

// MarshalJSON marshals Frame to JSON.
func (u *UnstructuredResource) MarshalJSON() ([]byte, error) {
	cfg := jsoniter.ConfigCompatibleWithStandardLibrary
	stream := cfg.BorrowStream(nil)
	defer cfg.ReturnStream(stream)

	writeResourceJSON(u, stream)
	if stream.Error != nil {
		return nil, stream.Error
	}

	return append([]byte(nil), stream.Buffer()...), nil
}

func init() { //nolint:gochecknoinits
	jsoniter.RegisterTypeEncoder("kindsys2.UnstructuredResource", &resourceCodec{})
	jsoniter.RegisterTypeDecoder("kindsys2.UnstructuredResource", &resourceCodec{})
}

type resourceCodec struct{}

func (codec *resourceCodec) IsEmpty(ptr unsafe.Pointer) bool {
	//f := (*UnstructuredResource)(ptr)
	return false // f.Fields == nil && f.RefID == "" && f.Meta == nil
}

func (codec *resourceCodec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	f := (*UnstructuredResource)(ptr)
	writeResourceJSON(f, stream)
}

func (codec *resourceCodec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	frame := UnstructuredResource{}
	err := readResourceJSON(&frame, iter)
	if err != nil {
		// keep existing iter error if it exists
		if iter.Error == nil {
			iter.Error = err
		}
		return
	}
	*((*UnstructuredResource)(ptr)) = frame
}

func writeResourceJSON(obj *UnstructuredResource, stream *jsoniter.Stream) {
	stream.WriteObjectStart()
	stream.WriteObjectField("apiVersion")
	stream.WriteString(obj.StaticMeta.Group + "/" + obj.StaticMeta.Version)
	stream.WriteMore()
	stream.WriteObjectField("kind")
	stream.WriteString(obj.StaticMeta.Kind)

	stream.WriteMore()
	stream.WriteObjectField("metadata")
	stream.WriteObjectStart()
	// "name": "ba2eea3b-d42c-4893-b522-6bf3b67efdc5",
	// "namespace": "default",
	// "resourceVersion": "123456",
	// "uid": "63396a47-eed8-46c3-859a-5dbac7c2b241"
	stream.WriteObjectEnd()

	if obj.Spec != nil {
		stream.WriteMore()
		stream.WriteObjectField("spec")
		stream.WriteVal(obj.Spec)
	}

	stream.WriteObjectEnd()
}

func readResourceJSON(obj *UnstructuredResource, iter *jsoniter.Iterator) error {
	for l1Field := iter.ReadObject(); l1Field != ""; l1Field = iter.ReadObject() {
		switch l1Field {
		case "apiVersion":
			vals := strings.SplitN(iter.ReadString(), "/", 2)
			obj.StaticMeta.Group = vals[0]
			obj.StaticMeta.Version = vals[1]

		case "kind":
			obj.StaticMeta.Kind = iter.ReadString()

		case "metadata":
			// TODO!!!!!
			_ = iter.ReadAny()

		case "spec":
			iter.ReadVal(&obj.Spec)

		default:
			iter.ReportError("bind l1", "unexpected field: "+l1Field)
		}
	}
	return iter.Error
}
