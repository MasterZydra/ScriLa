package runtime

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
}

type RuntimeVal struct {
	valueType ValueType
}

func (self *RuntimeVal) GetType() ValueType {
	return self.valueType
}

type INullVal interface {
	IRuntimeVal
	GetValue() string
}

type NullVal struct {
	valueType ValueType
	value     string
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

func (self *NullVal) ToString() string {
	return self.value
}

type IIntVal interface {
	IRuntimeVal
	GetValue() int64
}

type IntVal struct {
	valueType ValueType
	value     int64
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

func (self *IntVal) ToString() string {
	return fmt.Sprintf("%d", self.value)
}

type IBoolVal interface {
	IRuntimeVal
	GetValue() bool
}

type BoolVal struct {
	valueType ValueType
	value     bool
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

func (self *ObjVal) ToString() string {
	return "ObjVal"
}

type FunctionCall func(args []IRuntimeVal, env *Environment) IRuntimeVal

type INativeFunc interface {
	IRuntimeVal
	GetCall() FunctionCall
}

type NativeFunc struct {
	valueType ValueType
	call      FunctionCall
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

func (self *FunctionVal) ToString() string {
	return "FunctionVal"
}
