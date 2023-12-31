package kindsys2

import (
	"fmt"
	"strings"
	"time"
)

// A Resource is a single instance of a Grafana Resource kind
//
// The relationship between Resource and [Kind] is similar to the
// relationship between objects and classes in conventional object oriented
// design:
//
// - Objects are instantiated from classes. The name of the class is the type of object.
// - Resources are instantiated from kinds. The name of the kind is the type of resource.
//
// Resource is an interface, rather than a concrete struct, for two reasons:
//
// - Some use cases need to operate generically over resources of any kind.
// - Go generics do not allow the ergonomic expression of certain needed constraints.
//
// The [Core] and [Custom] interfaces are intended for the generic operation
// use case, fulfilling [Resource] using [UnstructuredResource].
//
// For known, specific kinds, it is usually possible to rely on code generation
// to produce a struct that implements [Resource] for each kind. Such a struct
// can be used as the generic type parameter to create a [TypedCore] or [TypedCustom]
type Resource interface {
	// CommonMetadata returns the Resource's CommonMetadata
	CommonMetadata() CommonMetadata

	// SetCommonMetadata overwrites the CommonMetadata of the object.
	// Implementations should always overwrite, rather than attempt merges of the metadata.
	// Callers wishing to merge should get current metadata with CommonMetadata() and set specific values.
	SetCommonMetadata(metadata CommonMetadata)

	// StaticMetadata returns the Resource's StaticMetadata
	StaticMetadata() StaticMetadata

	// SetStaticMetadata overwrites the Resource's StaticMetadata with the provided StaticMetadata.
	// Implementations should always overwrite, rather than attempt merges of the metadata.
	// Callers wishing to merge should get current metadata with StaticMetadata() and set specific values.
	// Note that StaticMetadata is only mutable in an object create context.
	SetStaticMetadata(metadata StaticMetadata)

	// CustomMetadata returns metadata unique to this Resource's kind, as opposed to Common and Static metadata,
	// which are the same across all kinds. An object may have no kind-specific CustomMetadata.
	// CustomMetadata can only be read from this interface, for use with resource.Client implementations,
	// those who wish to set CustomMetadata should use the interface's underlying type.
	CustomMetadata() CustomMetadata

	// SpecObject returns the actual "schema" object, which holds the main body of data
	SpecObject() any

	// Subresources returns a map of subresource name(s) to the object value for that subresource.
	// Spec is not considered a subresource, and should only be returned by SpecObject
	Subresources() map[string]any

	// Copy returns a full copy of the Resource with all its data
	Copy() Resource
}

// CustomMetadata is an interface describing a kindsys.Resource's kind-specific metadata
type CustomMetadata interface {
	// MapFields converts the custom metadata's fields into a map of field key to value.
	// This is used so Clients don't need to engage in reflection for marshaling metadata,
	// as various implementations may not store kind-specific metadata the same way.
	// TODO??? should this be map[string]string... since it needs to land in annotations?
	MapFields() map[string]any
}

// StaticMetadata consists of all non-mutable metadata for an object.
// It is set in the initial Create call for an Resource, then will always remain the same.
type StaticMetadata struct {
	Group     string `json:"group"`
	Version   string `json:"version"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

// GetAPIVersion returns the k8s style group + version
func (s StaticMetadata) GetAPIVersion() string {
	return s.Group + "/" + s.Version
}

func (s *StaticMetadata) setGroupVersionFromAPI(gv string) error {
	// this can be the internal version for the legacy kube types
	// TODO once we've cleared the last uses as strings, this special case should be removed.
	if (len(gv) == 0) || (gv == "/") {
		return nil
	}

	switch strings.Count(gv, "/") {
	case 0:
		s.Group = ""
		s.Version = gv
	case 1:
		i := strings.Index(gv, "/")
		s.Group = gv[:i]
		s.Version = gv[i+1:]
	default:
		return fmt.Errorf("unexpected GroupVersion string: %v", gv)
	}
	return nil
}

// Identifier creates an Identifier struct from the StaticMetadata
func (s StaticMetadata) Identifier() Identifier {
	return Identifier{
		Namespace: s.Namespace,
		Name:      s.Name,
	}
}

type Identifier struct {
	Namespace string
	Name      string
}

// CommonMetadata is the system-defined common metadata associated with a [Resource].
// It combines Kubernetes standard metadata with certain Grafana-specific additions.
//
// It is analogous to [k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta] in vanilla Kubernetes.
//
// TODO review this for optionality
type CommonMetadata struct {
	// UID is the unique ID of the object. This can be used to uniquely identify objects,
	// but is not guaranteed to be usable for lookups.
	UID string `json:"uid"`
	// ResourceVersion is a version string used to identify any and all changes to the object.
	// Any time the object changes in storage, the ResourceVersion will be changed.
	// This can be used to block updates if a change has been made to the object between when the object was
	// retrieved, and when the update was applied.
	ResourceVersion string `json:"resourceVersion,omitempty"`
	// Labels are string key/value pairs attached to the object. They can be used for filtering,
	// or as additional metadata.
	Labels map[string]string `json:"labels"`
	// CreationTimestamp indicates when the resource has been created.
	CreationTimestamp time.Time `json:"creationTimestamp"`
	// DeletionTimestamp indicates that the resource is pending deletion as of the provided time if non-nil.
	// Depending on implementation, this field may always be nil, or it may be a "tombstone" indicator.
	// It may also indicate that the system is waiting on some task to finish before the object is fully removed.
	DeletionTimestamp *time.Time `json:"deletionTimestamp,omitempty"`
	// Finalizers are a list of identifiers of interested parties for delete events for this resource.
	// Once a resource with finalizers has been deleted, the object should remain in the store,
	// DeletionTimestamp is set to the time of the "delete," and the resource will continue to exist
	// until the finalizers list is cleared.
	Finalizers []string `json:"finalizers,omitempty"`
	// UpdateTimestamp is the timestamp of the last update to the resource.
	UpdateTimestamp time.Time `json:"updateTimestamp"`
	// CreatedBy is a string which indicates the user or process which created the resource.
	// Implementations may choose what this indicator should be.
	CreatedBy string `json:"createdBy"`
	// UpdatedBy is a string which indicates the user or process which last updated the resource.
	// Implementations may choose what this indicator should be.
	UpdatedBy string `json:"updatedBy"`
	// Describe where the resource came from
	Origin *ResourceOriginInfo `json:"origin"`

	// ExtraFields stores implementation-specific metadata.
	// This is where the more esoteric k8s metadata will land
	// Not all Client implementations are required to honor all ExtraFields keys.
	// Generally, this field should be shied away from unless you know the specific
	// Client implementation you're working with and wish to track or mutate extra information.
	ExtraFields map[string]any `json:"extraFields"`
}

// ResourceOriginInfo is saved in annotations.  This is used to identify where the resource came from
// This object can model the same data as our existing provisioning table or a more general git sync
type ResourceOriginInfo struct {
	// Name of the origin/provisioning source
	Name string `json:"name,omitempty"`

	// The path within the named origin above (external_id in the existing dashboard provisioning)
	Path string `json:"path,omitempty"`

	// Verification/identification key (check_sum in existing dashboard provisioning)
	Key string `json:"key,omitempty"`

	// Origin modification timestamp when the resource was saved
	// This will be before the resource updated time
	Timestamp *time.Time `json:"time,omitempty"`

	// Avoid extending
	_ interface{}
}

// TODO guard against skew, use indirection through an internal package
// var _ CommonMetadata = encoding.CommonMetadata{}

// SimpleCustomMetadata is an implementation of CustomMetadata
type SimpleCustomMetadata map[string]any

// MapFields returns a map of string->value for all CustomMetadata fields
func (s SimpleCustomMetadata) MapFields() map[string]any {
	return s
}
