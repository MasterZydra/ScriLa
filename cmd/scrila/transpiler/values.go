package transpiler

import (
	"ScriLa/cmd/scrila/ast"
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

func (self *RuntimeVal) GetType() ValueType {
	return self.valueType
}

func (self *RuntimeVal) GetTranspilat() string {
	return self.transpilat
}

func (self *RuntimeVal) SetTranspilat(transpilat string) {
	self.transpilat = transpilat
}

type INullVal interface {
	IRuntimeVal
	GetValue() string
}

type NullVal struct {
	valueType  ValueType
	value      string
	transpilat string
}

func NewNullVal() *NullVal {
	return &NullVal{
		valueType: NullValueType,
		value:     "null",
	}
}

func (self *NullVal) GetType() ValueType {
	return self.valueType
}

func (self *NullVal) GetValue() string {
	return self.value
}

func (self *NullVal) GetTranspilat() string {
	return self.transpilat
}

func (self *NullVal) SetTranspilat(transpilat string) {
	self.transpilat = transpilat
}

func (self *NullVal) ToString() string {
	return self.value
}

type IIntVal interface {
	IRuntimeVal
	GetValue() int64
}

type IntVal struct {
	valueType  ValueType
	value      int64
	transpilat string
}

func NewIntVal(value int64) *IntVal {
	return &IntVal{
		valueType: IntValueType,
		value:     value,
	}
}

func (self *IntVal) GetType() ValueType {
	return self.valueType
}

func (self *IntVal) GetValue() int64 {
	return self.value
}

func (self *IntVal) GetTranspilat() string {
	return self.transpilat
}

func (self *IntVal) SetTranspilat(transpilat string) {
	self.transpilat = transpilat
}

func (self *IntVal) ToString() string {
	return fmt.Sprintf("%d", self.value)
}

type IBoolVal interface {
	IRuntimeVal
	GetValue() bool
}

type BoolVal struct {
	valueType  ValueType
	value      bool
	transpilat string
}

func (self *BoolVal) String() string {
	return fmt.Sprintf("&{%s %t}", self.GetType(), self.GetValue())
}

func NewBoolVal(value bool) *BoolVal {
	return &BoolVal{
		valueType: BoolValueType,
		value:     value,
	}
}

func (self *BoolVal) GetType() ValueType {
	return self.valueType
}

func (self *BoolVal) GetValue() bool {
	return self.value
}

func (self *BoolVal) GetTranspilat() string {
	return self.transpilat
}

func (self *BoolVal) SetTranspilat(transpilat string) {
	self.transpilat = transpilat
}

func (self *BoolVal) ToString() string {
	return fmt.Sprintf("%t", self.value)
}

type IObjVal interface {
	IRuntimeVal
	GetProperties() map[string]IRuntimeVal
}

type ObjVal struct {
	valueType  ValueType
	properties map[string]IRuntimeVal
	transpilat string
}

func NewObjVal() *ObjVal {
	return &ObjVal{
		valueType:  ObjValueType,
		properties: make(map[string]IRuntimeVal),
	}
}

func (self *ObjVal) GetType() ValueType {
	return self.valueType
}

func (self *ObjVal) GetProperties() map[string]IRuntimeVal {
	return self.properties
}

func (self *ObjVal) GetTranspilat() string {
	return self.transpilat
}

func (self *ObjVal) SetTranspilat(transpilat string) {
	self.transpilat = transpilat
}

func (self *ObjVal) ToString() string {
	return "ObjVal"
}

type IStrVal interface {
	IRuntimeVal
	GetValue() string
}

type StrVal struct {
	valueType  ValueType
	value      string
	transpilat string
}

func NewStrVal(value string) *StrVal {
	return &StrVal{
		valueType: StrValueType,
		value:     value,
	}
}

func (self *StrVal) GetType() ValueType {
	return self.valueType
}

func (self *StrVal) GetValue() string {
	return self.value
}

func (self *StrVal) GetTranspilat() string {
	return self.transpilat
}

func (self *StrVal) SetTranspilat(transpilat string) {
	self.transpilat = transpilat
}

func (self *StrVal) ToString() string {
	return self.value
}

type FunctionCall func(args []ast.IExpr, env *Environment) (IRuntimeVal, error)

type INativeFunc interface {
	IRuntimeVal
	GetCall() FunctionCall
}

type NativeFunc struct {
	valueType  ValueType
	call       FunctionCall
	transpilat string
}

func NewNativeFunc(function FunctionCall) *NativeFunc {
	return &NativeFunc{
		valueType: NativeFnType,
		call:      function,
	}
}

func (self *NativeFunc) GetType() ValueType {
	return self.valueType
}

func (self *NativeFunc) GetCall() FunctionCall {
	return self.call
}

func (self *NativeFunc) GetTranspilat() string {
	return self.transpilat
}

func (self *NativeFunc) SetTranspilat(transpilat string) {
	self.transpilat = transpilat
}

func (self *NativeFunc) ToString() string {
	return "NativeFunc"
}

type IFunctionVal interface {
	IRuntimeVal
	GetName() string
	GetParams() []string
	GetDeclarationEnv() *Environment
	GetBody() []ast.IStatement
}

type FunctionVal struct {
	valueType      ValueType
	name           string
	params         []string
	declarationEnv *Environment
	body           []ast.IStatement
	transpilat     string
}

func NewFunctionVal(name string, params []string, declarationEnv *Environment, body []ast.IStatement) *FunctionVal {
	return &FunctionVal{
		valueType:      FunctionValueType,
		name:           name,
		params:         params,
		declarationEnv: declarationEnv,
		body:           body,
	}
}

func (self *FunctionVal) GetType() ValueType {
	return self.valueType
}

func (self *FunctionVal) GetName() string {
	return self.name
}

func (self *FunctionVal) GetParams() []string {
	return self.params
}

func (self *FunctionVal) GetDeclarationEnv() *Environment {
	return self.declarationEnv
}

func (self *FunctionVal) GetBody() []ast.IStatement {
	return self.body
}

func (self *FunctionVal) GetTranspilat() string {
	return self.transpilat
}

func (self *FunctionVal) SetTranspilat(transpilat string) {
	self.transpilat = transpilat
}

func (self *FunctionVal) ToString() string {
	return "FunctionVal"
}
