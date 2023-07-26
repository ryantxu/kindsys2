package kindsys2

import "context"

type KindInfo struct {
	// Organization controlled prefix
	Group string `json:"group"`

	// Must be unique within the group (used as k8s "kind" name)
	Name string `json:"name"`

	// Description of the purpose of this kind
	Description string `json:"description"`

	// This kind depends on composable types
	UsesComposition []string `json:"usesComposition,omitempty"`

	// Indicate where this kind is in the dev cycle
	Maturity Maturity `json:"maturity"`
}

type VersionInfo struct {
	// Must be vMajor-Minor-alpha
	Version string `json:"version"`

	// Human readable descriptions of the changes in this version
	Changelog []string `json:"changelog,omitempty"`

	// The YYYY-MM-DD this version was published (or empty if not yet published)
	Published string `json:"published,omitempty"`

	// JSONSchema hash
	Signature string `json:"signature,omitempty"` // ?? hash of the json schema
}

// Manifest used for all kinds
type Manifest struct {
	KindInfo

	// Only valid for resource types
	MachineNames *MachineNames `json:"machineName,omitempty"`

	// Only valid for composable types
	ComposableType string `json:"type,omitempty"`

	// Only valid for composable types
	// ??? do we want/need multiple slots?  should each slot be a different type?
	ComposableSlots []string `json:"slots,omitempty"`

	// List of version info
	Versions []VersionInfo `json:"versions"`
}

type Kind interface {
	// Get the latest version
	CurrentVersion() string

	// Get all versions
	GetVersions() []VersionInfo

	// Return a JSON schema definition for the selected version
	GetJSONSchema(version string) (string, error)
}

type KindRegistry interface {
	// List the objects that can be saved as k8s style resources
	GetResourceKinds(ctx context.Context) []ResourceKind

	// List the valid composable types ("panel" | "dataquery" | "transformer" | "matcher")
	GetComposableTypes(ctx context.Context) []string

	// Get composable kinds of a given type
	GetComposableKinds(ctx context.Context, composableType string) []ComposableKind
}
