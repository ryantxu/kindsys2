package kindsys2

type KindInfo struct {
	// Organization controlled prefix
	Group string `json:"group"`

	// Must be unique within the group (used as k8s "kind" name)
	Name string `json:"name"`

	// Description of the purpose of this kind
	Description string `json:"description,omitempty"`

	// Exists when the kind can be used as a component within another kind
	// eg: "panel" | "query" | "transformer" | "matcher"
	ComposableType string `json:"composableType,omitempty"`

	// This kind requires access to other kinds in the registry
	RequiredComposableTypes []string `json:"requiredComposableTypes,omitempty"`

	// Machine names for different contexts
	MachineNames MachineName `json:"machineNames,omitempty"`

	// Maturity????
}

type MachineName struct {
	// This is used in k8s URLs
	Plural string `json:"plural,omitempty"`

	// Used as an alias in the display
	Singular string `json:"singular,omitempty"`

	// Optional shorter names that can be matched in a CLI
	Short []string `json:"short,omitempty"`
}

type VersionInfo struct {
	// The major version must increment when incompatible schema changes happen
	// Changes to the major version will require explicit
	Major int32 `json:"major"`

	// The minor version will increment when new fields are added (but not removed) from the schema
	Minor int32 `json:"minor"`

	// Indicate that this is not a released version
	State string `json:"state,omitempty"` // alpha / beta / dev (empty)

	// Human readable descriptions of the changes in this version
	Changelog []string `json:"changelog,omitempty"`

	// The YYYY-MM-DD this version was published (or empty if not yet published)
	Published string `json:"published,omitempty"`

	// JSONSchema hash
	Signature string `json:"omitempty,omitempty"` // ?? hash of the json schema
}

// ??? or simple struct?  JSON vs YAML... strict or lenient?
type DecoderOption = func(d Decoder)
type EncoderOption = func(e Encoder)

type Decoder interface {
	// Decode bytes, optionally get passed an explicit source version string
	Decode(b []byte, sourceVersionHint string) (any, error)
}

type Encoder interface {
	Encoder(obj any, targetVersion string) ([]byte, error)
}

type Kind interface {
	GetInfo() KindInfo
	CurrentVersion() VersionInfo
	GetVersions() []VersionInfo

	// Return a JSON schema definition for the selected version
	GetJSONSchema(version string) (string, error)

	GetDecoder(opts ...DecoderOption) (Decoder, error)
	GetEncoder(opts ...EncoderOption) (Encoder, error)
}

type KindRegistry interface {
	// List all known kinds
	GetKinds() []Kind

	// Get this list of composable types ("panel" | "query" | "transformer" | "matcher")
	GetComposableTypes() []string

	// List all
	GetComposableKinds(composableType string) []Kind
}
