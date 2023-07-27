package kindsys2

import (
	"bytes"
	"encoding/gob"

	jsoniter "github.com/json-iterator/go"
)

// Verify this implements a real resource
var _ Resource = &GenericResource[any, CustomMetadata, any]{}

// GenericResource
type GenericResource[Spec any, CustomMeta CustomMetadata, Status any] struct {
	StaticMeta StaticMetadata
	CommonMeta CommonMetadata
	CustomMeta CustomMeta
	Spec       Spec
	Status     Status
}

func (u *GenericResource[Spec, CustomMeta, Status]) SpecObject() any {
	return u.Spec
}

func (u *GenericResource[Spec, CustomMeta, Status]) Subresources() map[string]any {
	return map[string]any{
		"status": u.Status,
	}
}

func (u *GenericResource[Spec, CustomMeta, Status]) Copy() Resource {
	dst := &GenericResource[Spec, CustomMeta, Status]{}
	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(u)
	if err != nil {
		return dst // error
	}
	_ = gob.NewDecoder(&buf).Decode(dst)
	return dst
}

// UnmarshalJSON allows unmarshalling Frame from JSON.
func (u *GenericResource[Spec, CustomMeta, Status]) UnmarshalJSON(b []byte) error {
	//iter := jsoniter.ParseBytes(jsoniter.ConfigDefault, b)
	return nil //readResourceJSON(u, iter)
}

// MarshalJSON marshals Frame to JSON.
func (u *GenericResource[Spec, CustomMeta, Status]) MarshalJSON() ([]byte, error) {
	cfg := jsoniter.ConfigCompatibleWithStandardLibrary
	stream := cfg.BorrowStream(nil)
	defer cfg.ReturnStream(stream)

	WriteResourceJSON(u, stream)
	if stream.Error != nil {
		return nil, stream.Error
	}

	return append([]byte(nil), stream.Buffer()...), nil
}

// CommonMetadata returns the object's CommonMetadata
func (u *GenericResource[Spec, CustomMeta, Status]) CommonMetadata() CommonMetadata {
	return u.CommonMeta
}

// SetCommonMetadata overwrites the ObjectMetadata.Common() supplied by BasicMetadataObject.ObjectMetadata()
func (u *GenericResource[Spec, CustomMeta, Status]) SetCommonMetadata(m CommonMetadata) {
	u.CommonMeta = m
}

// StaticMetadata returns the object's StaticMetadata
func (u *GenericResource[Spec, CustomMeta, Status]) StaticMetadata() StaticMetadata {
	return u.StaticMeta
}

// SetStaticMetadata overwrites the StaticMetadata supplied by BasicMetadataObject.StaticMetadata()
func (u *GenericResource[Spec, CustomMeta, Status]) SetStaticMetadata(m StaticMetadata) {
	u.StaticMeta = m
}

// CustomMetadata returns the object's CustomMetadata
func (u *GenericResource[Spec, CustomMeta, Status]) CustomMetadata() CustomMetadata {
	return u.CustomMeta
}
