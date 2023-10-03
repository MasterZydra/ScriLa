package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"ScriLa/cmd/scrila/lexer"
	"fmt"
)

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

func doTypesMatch(type1 lexer.TokenType, type2 ValueType) bool {
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

// NullVal

type INullVal interface {
	IRuntimeVal
	GetValue() string
}

type NullVal struct {
	runtimeVal *RuntimeVal
	value      string
}

func NewNullVal() *NullVal {
	return &NullVal{
		runtimeVal: NewRuntimeVal(NullValueType),
		value:      "null",
	}
}

func (self *NullVal) GetType() ValueType {
	return self.runtimeVal.GetType()
}

func (self *NullVal) GetValue() string {
	return self.value
}

func (self *NullVal) GetTranspilat() string {
	return self.runtimeVal.transpilat
}

func (self *NullVal) SetTranspilat(transpilat string) {
	self.runtimeVal.transpilat = transpilat
}

func (self *NullVal) ToString() string {
	return self.value
}

// IntVal

type IIntVal interface {
	IRuntimeVal
	GetValue() int64
}

type IntVal struct {
	runtimeVal *RuntimeVal
	value      int64
}

func NewIntVal(value int64) *IntVal {
	return &IntVal{
		runtimeVal: NewRuntimeVal(IntValueType),
		value:      value,
	}
}

func (self *IntVal) GetType() ValueType {
	return self.runtimeVal.GetType()
}

func (self *IntVal) GetValue() int64 {
	return self.value
}

func (self *IntVal) GetTranspilat() string {
	return self.runtimeVal.transpilat
}

func (self *IntVal) SetTranspilat(transpilat string) {
	self.runtimeVal.transpilat = transpilat
}

func (self *IntVal) ToString() string {
	return fmt.Sprintf("%d", self.value)
}

// BoolVal

type IBoolVal interface {
	IRuntimeVal
	GetValue() bool
}

type BoolVal struct {
	runtimeVal *RuntimeVal
	value      bool
}

func (self *BoolVal) String() string {
	return fmt.Sprintf("&{%s %t}", self.GetType(), self.GetValue())
}

func NewBoolVal(value bool) *BoolVal {
	return &BoolVal{
		runtimeVal: NewRuntimeVal(BoolValueType),
		value:      value,
	}
}

func (self *BoolVal) GetType() ValueType {
	return self.runtimeVal.GetType()
}

func (self *BoolVal) GetValue() bool {
	return self.value
}

func (self *BoolVal) GetTranspilat() string {
	return self.runtimeVal.transpilat
}

func (self *BoolVal) SetTranspilat(transpilat string) {
	self.runtimeVal.transpilat = transpilat
}

func (self *BoolVal) ToString() string {
	return fmt.Sprintf("%t", self.value)
}

// ObjVal

type IObjVal interface {
	IRuntimeVal
	GetProperties() map[string]IRuntimeVal
}

type ObjVal struct {
	runtimeVal *RuntimeVal
	properties map[string]IRuntimeVal
}

func NewObjVal() *ObjVal {
	return &ObjVal{
		runtimeVal: NewRuntimeVal(ObjValueType),
		properties: make(map[string]IRuntimeVal),
	}
}

func (self *ObjVal) GetType() ValueType {
	return self.runtimeVal.GetType()
}

func (self *ObjVal) GetProperties() map[string]IRuntimeVal {
	return self.properties
}

func (self *ObjVal) GetTranspilat() string {
	return self.runtimeVal.transpilat
}

func (self *ObjVal) SetTranspilat(transpilat string) {
	self.runtimeVal.transpilat = transpilat
}

func (self *ObjVal) ToString() string {
	return "ObjVal"
}

// StrVal

type IStrVal interface {
	IRuntimeVal
	GetValue() string
}

type StrVal struct {
	runtimeVal *RuntimeVal
	value      string
}

func NewStrVal(value string) *StrVal {
	return &StrVal{
		runtimeVal: NewRuntimeVal(StrValueType),
		value:      value,
	}
}

func (self *StrVal) GetType() ValueType {
	return self.runtimeVal.GetType()
}

func (self *StrVal) GetValue() string {
	return self.value
}

func (self *StrVal) GetTranspilat() string {
	return self.runtimeVal.transpilat
}

func (self *StrVal) SetTranspilat(transpilat string) {
	self.runtimeVal.transpilat = transpilat
}

func (self *StrVal) ToString() string {
	return self.value
}

// NativeFunc

type FunctionCall func(args []ast.IExpr, env *Environment) (IRuntimeVal, error)

type INativeFunc interface {
	IRuntimeVal
	GetCall() FunctionCall
}

type NativeFunc struct {
	runtimeVal *RuntimeVal
	call       FunctionCall
}

func NewNativeFunc(function FunctionCall) *NativeFunc {
	return &NativeFunc{
		runtimeVal: NewRuntimeVal(NativeFnType),
		call:       function,
	}
}

func (self *NativeFunc) GetType() ValueType {
	return self.runtimeVal.GetType()
}

func (self *NativeFunc) GetCall() FunctionCall {
	return self.call
}

func (self *NativeFunc) GetTranspilat() string {
	return self.runtimeVal.transpilat
}

func (self *NativeFunc) SetTranspilat(transpilat string) {
	self.runtimeVal.transpilat = transpilat
}

func (self *NativeFunc) ToString() string {
	return "NativeFunc"
}

// FunctionVal

type IFunctionVal interface {
	IRuntimeVal
	GetName() string
	GetParams() []*ast.Parameter
	GetDeclarationEnv() *Environment
	GetBody() []ast.IStatement
	GetReturnType() lexer.TokenType
}

type FunctionVal struct {
	runtimeVal     *RuntimeVal
	name           string
	params         []*ast.Parameter
	declarationEnv *Environment
	body           []ast.IStatement
	returnType     lexer.TokenType
}

func NewFunctionVal(funcDeclaration ast.IFunctionDeclaration, env *Environment) *FunctionVal {
	return &FunctionVal{
		runtimeVal:     NewRuntimeVal(FunctionValueType),
		name:           funcDeclaration.GetName(),
		params:         funcDeclaration.GetParameters(),
		declarationEnv: env,
		body:           funcDeclaration.GetBody(),
		returnType:     funcDeclaration.GetReturnType(),
	}
}

func (self *FunctionVal) GetType() ValueType {
	return self.runtimeVal.GetType()
}

func (self *FunctionVal) GetName() string {
	return self.name
}

func (self *FunctionVal) GetParams() []*ast.Parameter {
	return self.params
}

func (self *FunctionVal) GetDeclarationEnv() *Environment {
	return self.declarationEnv
}

func (self *FunctionVal) GetBody() []ast.IStatement {
	return self.body
}

func (self *FunctionVal) GetReturnType() lexer.TokenType {
	return self.returnType
}

func (self *FunctionVal) GetTranspilat() string {
	return self.runtimeVal.transpilat
}

func (self *FunctionVal) SetTranspilat(transpilat string) {
	self.runtimeVal.transpilat = transpilat
}

func (self *FunctionVal) ToString() string {
	return "FunctionVal"
}
