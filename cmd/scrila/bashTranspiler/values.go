package bashTranspiler

import (
	"ScriLa/cmd/scrila/ast"
	"ScriLa/cmd/scrila/lexer"
	"fmt"
)

// NullVal

type INullVal interface {
	ast.IRuntimeVal
	GetValue() string
}

type NullVal struct {
	runtimeVal *ast.RuntimeVal
	value      string
}

func NewNullVal() *NullVal {
	return &NullVal{
		runtimeVal: ast.NewRuntimeVal(ast.NullValueType),
		value:      "null",
	}
}

func (self *NullVal) GetType() ast.ValueType {
	return self.runtimeVal.GetType()
}

func (self *NullVal) GetValue() string {
	return self.value
}

func (self *NullVal) GetTranspilat() string {
	return self.runtimeVal.GetTranspilat()
}

func (self *NullVal) SetTranspilat(transpilat string) {
	self.runtimeVal.SetTranspilat(transpilat)
}

func (self *NullVal) ToString() string {
	return self.value
}

// IntVal

type IIntVal interface {
	ast.IRuntimeVal
	GetValue() int64
}

type IntVal struct {
	runtimeVal *ast.RuntimeVal
	value      int64
}

func NewIntVal(value int64) *IntVal {
	return &IntVal{
		runtimeVal: ast.NewRuntimeVal(ast.IntValueType),
		value:      value,
	}
}

func (self *IntVal) GetType() ast.ValueType {
	return self.runtimeVal.GetType()
}

func (self *IntVal) GetValue() int64 {
	return self.value
}

func (self *IntVal) GetTranspilat() string {
	return self.runtimeVal.GetTranspilat()
}

func (self *IntVal) SetTranspilat(transpilat string) {
	self.runtimeVal.SetTranspilat(transpilat)
}

func (self *IntVal) ToString() string {
	return fmt.Sprintf("%d", self.value)
}

// BoolVal

type IBoolVal interface {
	ast.IRuntimeVal
	GetValue() bool
}

type BoolVal struct {
	runtimeVal *ast.RuntimeVal
	value      bool
}

func (self *BoolVal) String() string {
	return fmt.Sprintf("&{%s %t}", self.GetType(), self.GetValue())
}

func NewBoolVal(value bool) *BoolVal {
	return &BoolVal{
		runtimeVal: ast.NewRuntimeVal(ast.BoolValueType),
		value:      value,
	}
}

func (self *BoolVal) GetType() ast.ValueType {
	return self.runtimeVal.GetType()
}

func (self *BoolVal) GetValue() bool {
	return self.value
}

func (self *BoolVal) GetTranspilat() string {
	return self.runtimeVal.GetTranspilat()
}

func (self *BoolVal) SetTranspilat(transpilat string) {
	self.runtimeVal.SetTranspilat(transpilat)
}

func (self *BoolVal) ToString() string {
	return fmt.Sprintf("%t", self.value)
}

// ObjVal

type IObjVal interface {
	ast.IRuntimeVal
	GetProperties() map[string]ast.IRuntimeVal
}

type ObjVal struct {
	runtimeVal *ast.RuntimeVal
	properties map[string]ast.IRuntimeVal
}

func NewObjVal() *ObjVal {
	return &ObjVal{
		runtimeVal: ast.NewRuntimeVal(ast.ObjValueType),
		properties: make(map[string]ast.IRuntimeVal),
	}
}

func (self *ObjVal) GetType() ast.ValueType {
	return self.runtimeVal.GetType()
}

func (self *ObjVal) GetProperties() map[string]ast.IRuntimeVal {
	return self.properties
}

func (self *ObjVal) GetTranspilat() string {
	return self.runtimeVal.GetTranspilat()
}

func (self *ObjVal) SetTranspilat(transpilat string) {
	self.runtimeVal.SetTranspilat(transpilat)
}

func (self *ObjVal) ToString() string {
	return "ObjVal"
}

// StrVal

type IStrVal interface {
	ast.IRuntimeVal
	GetValue() string
}

type StrVal struct {
	runtimeVal *ast.RuntimeVal
	value      string
}

func NewStrVal(value string) *StrVal {
	return &StrVal{
		runtimeVal: ast.NewRuntimeVal(ast.StrValueType),
		value:      value,
	}
}

func (self *StrVal) GetType() ast.ValueType {
	return self.runtimeVal.GetType()
}

func (self *StrVal) GetValue() string {
	return self.value
}

func (self *StrVal) GetTranspilat() string {
	return self.runtimeVal.GetTranspilat()
}

func (self *StrVal) SetTranspilat(transpilat string) {
	self.runtimeVal.SetTranspilat(transpilat)
}

func (self *StrVal) ToString() string {
	return self.value
}

// NativeFunc

type FunctionCall func(args []ast.IExpr, env *Environment) (ast.IRuntimeVal, error)

type INativeFunc interface {
	ast.IRuntimeVal
	GetCall() FunctionCall
	GetReturnType() lexer.TokenType
}

type NativeFunc struct {
	runtimeVal *ast.RuntimeVal
	call       FunctionCall
	returnType lexer.TokenType
}

func NewNativeFunc(function FunctionCall, returnType lexer.TokenType) *NativeFunc {
	return &NativeFunc{
		runtimeVal: ast.NewRuntimeVal(ast.NativeFnType),
		call:       function,
		returnType: returnType,
	}
}

func (self *NativeFunc) GetType() ast.ValueType {
	return self.runtimeVal.GetType()
}

func (self *NativeFunc) GetCall() FunctionCall {
	return self.call
}

func (self *NativeFunc) GetReturnType() lexer.TokenType {
	return self.returnType
}

func (self *NativeFunc) GetTranspilat() string {
	return self.runtimeVal.GetTranspilat()
}

func (self *NativeFunc) SetTranspilat(transpilat string) {
	self.runtimeVal.SetTranspilat(transpilat)
}

func (self *NativeFunc) ToString() string {
	return "NativeFunc"
}

// FunctionVal

type IFunctionVal interface {
	ast.IRuntimeVal
	GetName() string
	GetParams() []*ast.Parameter
	GetDeclarationEnv() *Environment
	GetBody() []ast.IStatement
	GetReturnType() lexer.TokenType
}

type FunctionVal struct {
	runtimeVal     *ast.RuntimeVal
	name           string
	params         []*ast.Parameter
	declarationEnv *Environment
	body           []ast.IStatement
	returnType     lexer.TokenType
}

func NewFunctionVal(funcDeclaration ast.IFunctionDeclaration, env *Environment) *FunctionVal {
	return &FunctionVal{
		runtimeVal:     ast.NewRuntimeVal(ast.FunctionValueType),
		name:           funcDeclaration.GetName(),
		params:         funcDeclaration.GetParameters(),
		declarationEnv: env,
		body:           funcDeclaration.GetBody(),
		returnType:     funcDeclaration.GetReturnType(),
	}
}

func (self *FunctionVal) GetType() ast.ValueType {
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
	return self.runtimeVal.GetTranspilat()
}

func (self *FunctionVal) SetTranspilat(transpilat string) {
	self.runtimeVal.SetTranspilat(transpilat)
}

func (self *FunctionVal) ToString() string {
	return "FunctionVal"
}
