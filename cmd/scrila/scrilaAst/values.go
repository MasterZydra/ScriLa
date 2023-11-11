package scrilaAst

import "ScriLa/cmd/scrila/lexer"

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

var valueTypeLexerTypeMapping = map[ValueType]lexer.TokenType{
	BoolValueType: lexer.BoolType,
	IntValueType:  lexer.IntType,
	ObjValueType:  lexer.ObjType,
	StrValueType:  lexer.StrType,
}

func DoTypesMatch(type1 lexer.TokenType, type2 ValueType) bool {
	value, ok := valueTypeLexerTypeMapping[type2]
	if !ok {
		return false
	}
	return value == type1
}

// RuntimeVal

type IRuntimeVal interface {
	GetType() ValueType
	ToString() string
	GetTranspilat() string
	SetTranspilat(transpilat string)
}

type RuntimeVal struct {
	valueType  ValueType
	transpilat string
}

func NewRuntimeVal(valueType ValueType) *RuntimeVal {
	return &RuntimeVal{valueType: valueType}
}

func (self *RuntimeVal) GetType() ValueType {
	return self.valueType
}

func (self *RuntimeVal) GetTranspilat() string {
	return self.transpilat
}

func (self *RuntimeVal) SetTranspilat(transpilat string) {
	self.transpilat = transpilat
}
