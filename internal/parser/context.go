package parser

// FieldDef object
type FieldDef struct {
	// MethodName in model instance which returns this field
	MethodName string
}

// ModelDef options
type ModelDef struct {
	// Table name
	Table string

	// Fields map by their model instance getter methods names
	Fields map[string]*FieldDef

	// TypeDefined in parsed file, and no need generate it struct
	TypeDefined bool
}

// NewFieldDef instance constructo, name should be method name not name from table
func NewFieldDef(fieldName string) *FieldDef {
	return &FieldDef{
		// MethodName which return field from model
		MethodName: fieldName,
	}
}

// NewModelDef instance constructor
func NewModelDef(modelName string) *ModelDef {
	return &ModelDef{
		Table:  modelName,
		Fields: make(map[string]*FieldDef),
	}
}

// Context for parsing
type Context struct {
	// currentDef for parsing
	currentDef *ModelDef

	// SchemaPkgNames contains all imported names for schema package
	SchemaPkgNames map[string]struct{}

	// DefList map
	DefList map[string]*ModelDef

	// PkgFile path for parse
	PkgFile string

	// PkgName contains name of parsed package
	PkgName string

	// ModelImport name, by default `model`
	ModelImport string

	// FieldImport name, by default `field`
	FieldImport string
}
