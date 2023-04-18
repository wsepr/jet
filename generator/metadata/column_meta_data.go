package metadata

// Column struct
type Column struct {
	Name         string
	IsPrimaryKey bool
	MaxLength    int
	IsNullable   bool
	IsGenerated  bool
	DataType     DataType
	Comment      string
}

// DataTypeKind is database type kind(base, enum, user-defined, array)
type DataTypeKind string

// DataTypeKind possible values
const (
	BaseType        DataTypeKind = "base"
	EnumType        DataTypeKind = "enum"
	UserDefinedType DataTypeKind = "user-defined"
	ArrayType       DataTypeKind = "array"
)

// DataType contains information about column data type
type DataType struct {
	Name       string
	Kind       DataTypeKind
	IsUnsigned bool
}
