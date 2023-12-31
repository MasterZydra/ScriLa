package bashTranspiler

import (
	"ScriLa/cmd/scrila/scrilaAst"
	"fmt"
)

// Array

type IArrayVal interface {
	scrilaAst.IRuntimeVal
}

func NewArrayVal(valueType scrilaAst.ValueType) *scrilaAst.RuntimeVal {
	return scrilaAst.NewRuntimeVal(valueType)
}

// NullVal

type INullVal interface {
	scrilaAst.IRuntimeVal
	GetValue() string
}

type NullVal struct {
	runtimeVal *scrilaAst.RuntimeVal
	value      string
}

func NewNullVal() *NullVal {
	return &NullVal{
		runtimeVal: scrilaAst.NewRuntimeVal(scrilaAst.NullValueType),
		value:      "null",
	}
}

func (self *NullVal) GetType() scrilaAst.ValueType {
	return self.runtimeVal.GetType()
}

func (self *NullVal) GetValue() string {
	return self.value
}

// IntVal

type IIntVal interface {
	scrilaAst.IRuntimeVal
	GetValue() int64
}

type IntVal struct {
	runtimeVal *scrilaAst.RuntimeVal
	value      int64
}

func NewIntVal(value int64) *IntVal {
	return &IntVal{
		runtimeVal: scrilaAst.NewRuntimeVal(scrilaAst.IntValueType),
		value:      value,
	}
}

func (self *IntVal) GetType() scrilaAst.ValueType {
	return self.runtimeVal.GetType()
}

func (self *IntVal) GetValue() int64 {
	return self.value
}

// BoolVal

type IBoolVal interface {
	scrilaAst.IRuntimeVal
	GetValue() bool
}

type BoolVal struct {
	runtimeVal *scrilaAst.RuntimeVal
	value      bool
}

func (self *BoolVal) String() string {
	return fmt.Sprintf("&{%s %t}", self.GetType(), self.GetValue())
}

func NewBoolVal(value bool) *BoolVal {
	return &BoolVal{
		runtimeVal: scrilaAst.NewRuntimeVal(scrilaAst.BoolValueType),
		value:      value,
	}
}

func (self *BoolVal) GetType() scrilaAst.ValueType {
	return self.runtimeVal.GetType()
}

func (self *BoolVal) GetValue() bool {
	return self.value
}

// ObjVal

type IObjVal interface {
	scrilaAst.IRuntimeVal
	GetProperties() map[string]scrilaAst.IRuntimeVal
}

type ObjVal struct {
	runtimeVal *scrilaAst.RuntimeVal
	properties map[string]scrilaAst.IRuntimeVal
}

func NewObjVal() *ObjVal {
	return &ObjVal{
		runtimeVal: scrilaAst.NewRuntimeVal(scrilaAst.ObjValueType),
		properties: make(map[string]scrilaAst.IRuntimeVal),
	}
}

func (self *ObjVal) GetType() scrilaAst.ValueType {
	return self.runtimeVal.GetType()
}

func (self *ObjVal) GetProperties() map[string]scrilaAst.IRuntimeVal {
	return self.properties
}

// StrVal

type IStrVal interface {
	scrilaAst.IRuntimeVal
	GetValue() string
}

type StrVal struct {
	runtimeVal *scrilaAst.RuntimeVal
	value      string
}

func NewStrVal(value string) *StrVal {
	return &StrVal{
		runtimeVal: scrilaAst.NewRuntimeVal(scrilaAst.StrValueType),
		value:      value,
	}
}

func (self *StrVal) GetType() scrilaAst.ValueType {
	return self.runtimeVal.GetType()
}

func (self *StrVal) GetValue() string {
	return self.value
}

// NativeFunc

type FunctionCall func(args []scrilaAst.IExpr, env *Environment) (scrilaAst.IRuntimeVal, error)

type INativeFunc interface {
	scrilaAst.IRuntimeVal
	GetCall() FunctionCall
	GetReturnType() scrilaAst.NodeType
}

type NativeFunc struct {
	runtimeVal *scrilaAst.RuntimeVal
	call       FunctionCall
	returnType scrilaAst.NodeType
}

func NewNativeFunc(function FunctionCall, returnType scrilaAst.NodeType) *NativeFunc {
	return &NativeFunc{
		runtimeVal: scrilaAst.NewRuntimeVal(scrilaAst.NativeFnType),
		call:       function,
		returnType: returnType,
	}
}

func (self *NativeFunc) GetType() scrilaAst.ValueType {
	return self.runtimeVal.GetType()
}

func (self *NativeFunc) GetCall() FunctionCall {
	return self.call
}

func (self *NativeFunc) GetReturnType() scrilaAst.NodeType {
	return self.returnType
}

// FunctionVal

type IFunctionVal interface {
	scrilaAst.IRuntimeVal
	GetName() string
	GetParams() []*scrilaAst.Parameter
	GetDeclarationEnv() *Environment
	GetBody() []scrilaAst.IStatement
	GetReturnType() scrilaAst.NodeType
}

type FunctionVal struct {
	runtimeVal     *scrilaAst.RuntimeVal
	name           string
	params         []*scrilaAst.Parameter
	declarationEnv *Environment
	body           []scrilaAst.IStatement
	returnType     scrilaAst.NodeType
}

func NewFunctionVal(funcDeclaration scrilaAst.IFunctionDeclaration, env *Environment) *FunctionVal {
	return &FunctionVal{
		runtimeVal:     scrilaAst.NewRuntimeVal(scrilaAst.FunctionValueType),
		name:           funcDeclaration.GetName(),
		params:         funcDeclaration.GetParameters(),
		declarationEnv: env,
		body:           funcDeclaration.GetBody(),
		returnType:     funcDeclaration.GetReturnType(),
	}
}

func (self *FunctionVal) GetType() scrilaAst.ValueType {
	return self.runtimeVal.GetType()
}

func (self *FunctionVal) GetName() string {
	return self.name
}

func (self *FunctionVal) GetParams() []*scrilaAst.Parameter {
	return self.params
}

func (self *FunctionVal) GetDeclarationEnv() *Environment {
	return self.declarationEnv
}

func (self *FunctionVal) GetBody() []scrilaAst.IStatement {
	return self.body
}

func (self *FunctionVal) GetReturnType() scrilaAst.NodeType {
	return self.returnType
}
