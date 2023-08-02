package kindsys2

import (
	"context"
	"io"
)

// Concrete GRD implementation
// An api server will be created with support for the kind and the various
type KindServiceHooks struct {
	// This defines the apiVersion + Kind
	Kind ResourceKind

	// Called before creating a new resource (admission/mutation controller)
	// this can error or mutate
	BeforeAdd func(ctx context.Context, obj Resource) (Resource, error)

	// Called before updating a resource  (admission controller)
	// this can error or mutate
	BeforeUpdate func(ctx context.Context, oldObj Resource, newObj Resource) (Resource, error)

	// Called before deleting a resource
	// this can error
	// ??? is this necessary
	BeforeDelete func(ctx context.Context, obj Resource) error

	// This is called when initialized -- the endpoints will be added to the api server
	// the OpenAPI specs will be exposed in the public API
	GetRawAPIHandlers func(getter ResourceGetter) []RawAPIHandler
}

// This allows access to resources for API handlers
type ResourceGetter = func(ctx context.Context, id StaticMetadata) (Resource, error)

// This is used to answer raw API requests like /logs
type StreamingResponse = func(ctx context.Context, apiVersion, acceptHeader string) (stream io.ReadCloser, flush bool, mimeType string, err error)

// This is used to implement dynamic sub-resources like pods/x/logs
type RawAPIHandler struct {
	Path    string
	OpenAPI string
	Level   RawAPILevel // resource | namespace | group

	// The GET request + response (see the standard /history and /refs)
	Handler func(ctx context.Context, id StaticMetadata) StreamingResponse
}

type RawAPILevel int8

const (
	RawAPILevelResource Maturity = iota
	RawAPILevelNamespace
	RawAPILevelGroup
)
