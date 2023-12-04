package scrilaAst

type ValueType string

const (
	BoolArrayValueType ValueType = "bool-array"
	BoolValueType      ValueType = "bool"
	FunctionValueType  ValueType = "function"
	IntArrayValueType  ValueType = "int-array"
	IntValueType       ValueType = "int"
	NativeFnType       ValueType = "native-func"
	NullValueType      ValueType = "null"
	ObjValueType       ValueType = "obj"
	StrArrayValueType  ValueType = "str-array"
	StrValueType       ValueType = "str"
)

var nodeTypeValueTypeMapping = map[ValueType]NodeType{
	BoolArrayValueType: BoolArrayNode,
	BoolValueType:      BoolLiteralNode,
	IntArrayValueType:  IntArrayNode,
	IntValueType:       IntLiteralNode,
	ObjValueType:       ObjectLiteralNode,
	StrArrayValueType:  StrArrayNode,
	StrValueType:       StrLiteralNode,
}

func DoTypesMatch(type1 NodeType, type2 ValueType) bool {
	value, ok := nodeTypeValueTypeMapping[type2]
	if !ok {
		return false
	}
	return value == type1
}

// RuntimeVal

type IRuntimeVal interface {
	GetType() ValueType
}

type RuntimeVal struct {
	valueType ValueType
}

func NewRuntimeVal(valueType ValueType) *RuntimeVal {
	return &RuntimeVal{valueType: valueType}
}

func (self *RuntimeVal) GetType() ValueType {
	return self.valueType
}
