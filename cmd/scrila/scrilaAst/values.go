package scrilaAst

type ValueType string

const (
	BoolValueType     ValueType = "bool"
	FunctionValueType ValueType = "function"
	IntValueType      ValueType = "int"
	NativeFnType      ValueType = "native-func"
	NullValueType     ValueType = "null"
	ObjValueType      ValueType = "obj"
	StrValueType      ValueType = "str"
)

var nodeTypeValueTypeMapping = map[ValueType]NodeType{
	BoolValueType: BoolLiteralNode,
	IntValueType:  IntLiteralNode,
	ObjValueType:  ObjectLiteralNode,
	StrValueType:  StrLiteralNode,
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
