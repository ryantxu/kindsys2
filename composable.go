package kindsys2

type ComposableKind interface {
	Kind

	// can migrate map[string]any ???
	GetComposableType() string

	// panel currently has Options + FieldConfig
	// can we get rid of slots and just have two composable kinds in the same plugin?
	GetComposableSlots() []string

	// Given an object (at a version) check that it is valid
	Validate(obj any, sourceVersion string) error

	// Migrate from one version of the object to another
	Migrate(obj any, sourceVersion string, targetVersion string) (any, error)
}
