package kindsys2

var _ Resource = &UnstructuredResource{}

// UnstructuredResource is an untyped representation of [Resource].
type UnstructuredResource = GenericResource[map[string]any, SimpleCustomMetadata, map[string]any]
