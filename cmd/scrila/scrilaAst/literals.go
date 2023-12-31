package scrilaAst

import (
	"fmt"
)

// Array

type IArray interface {
	IExpr
	AddValue(value IExpr)
	GetValues() []IExpr
}

type Array struct {
	expr   *Expr
	values []IExpr
}

func (self *Array) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - id: %d", self.GetKind(), self.GetId())
	for i, val := range self.values {
		str += fmt.Sprintf("\n%s%d: %s", indent(), i, val)
	}
	indentDepth--
	return str + "}"
}

func NewArray(ln int, col int) *Array {
	return &Array{expr: NewExpr(ArrayLiteralNode, ln, col)}
}

func (self *Array) AddValue(value IExpr) {
	self.values = append(self.values, value)
}

func (self *Array) GetId() int {
	return self.expr.GetId()
}

func (self *Array) GetKind() NodeType {
	return self.expr.GetKind()
}

func (self *Array) GetValues() []IExpr {
	return self.values
}

func (self *Array) GetLn() int {
	return self.expr.GetLn()
}

func (self *Array) GetCol() int {
	return self.expr.GetCol()
}

func (self *Array) GetResult() IRuntimeVal {
	return self.expr.GetResult()
}

func (self *Array) SetResult(value IRuntimeVal) {
	self.expr.SetResult(value)
}

// BoolLiteral

type IBoolLiteral interface {
	IExpr
	GetValue() bool
}

type BoolLiteral struct {
	expr  *Expr
	value bool
}

func (self *BoolLiteral) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - id: %d, value: %t}", self.GetKind(), self.GetId(), self.GetValue())
	indentDepth--
	return str
}

func NewBoolLiteral(value bool, ln int, col int) *BoolLiteral {
	return &BoolLiteral{
		expr:  NewExpr(BoolLiteralNode, ln, col),
		value: value,
	}
}

func (self *BoolLiteral) GetId() int {
	return self.expr.GetId()
}

func (self *BoolLiteral) GetKind() NodeType {
	return self.expr.GetKind()
}

func (self *BoolLiteral) GetValue() bool {
	return self.value
}

func (self *BoolLiteral) GetLn() int {
	return self.expr.GetLn()
}

func (self *BoolLiteral) GetCol() int {
	return self.expr.GetCol()
}

func (self *BoolLiteral) GetResult() IRuntimeVal {
	return self.expr.GetResult()
}

func (self *BoolLiteral) SetResult(value IRuntimeVal) {
	self.expr.SetResult(value)
}

// IntLiteral

type IIntLiteral interface {
	IExpr
	GetValue() int64
}

type IntLiteral struct {
	expr  *Expr
	value int64
}

func (self *IntLiteral) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - id: %d, value: %d}", self.GetKind(), self.GetId(), self.GetValue())
	indentDepth--
	return str
}

func NewIntLiteral(value int64, ln int, col int) *IntLiteral {
	return &IntLiteral{
		expr:  NewExpr(IntLiteralNode, ln, col),
		value: value,
	}
}

func (self *IntLiteral) GetId() int {
	return self.expr.GetId()
}

func (self *IntLiteral) GetKind() NodeType {
	return self.expr.GetKind()
}

func (self *IntLiteral) GetValue() int64 {
	return self.value
}

func (self *IntLiteral) GetLn() int {
	return self.expr.GetLn()
}

func (self *IntLiteral) GetCol() int {
	return self.expr.GetCol()
}

func (self *IntLiteral) GetResult() IRuntimeVal {
	return self.expr.GetResult()
}

func (self *IntLiteral) SetResult(value IRuntimeVal) {
	self.expr.SetResult(value)
}

// StrLiteral

type IStrLiteral interface {
	IExpr
	GetValue() string
}

type StrLiteral struct {
	expr  *Expr
	value string
}

func (self *StrLiteral) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - id: %d, value: '%s'}", self.GetKind(), self.GetId(), self.GetValue())
	indentDepth--
	return str
}

func NewStrLiteral(value string, ln int, col int) *StrLiteral {
	return &StrLiteral{
		expr:  NewExpr(StrLiteralNode, ln, col),
		value: value,
	}
}

func (self *StrLiteral) GetId() int {
	return self.expr.GetId()
}

func (self *StrLiteral) GetKind() NodeType {
	return self.expr.GetKind()
}

func (self *StrLiteral) GetValue() string {
	return self.value
}

func (self *StrLiteral) GetLn() int {
	return self.expr.GetLn()
}

func (self *StrLiteral) GetCol() int {
	return self.expr.GetCol()
}

func (self *StrLiteral) GetResult() IRuntimeVal {
	return self.expr.GetResult()
}

func (self *StrLiteral) SetResult(value IRuntimeVal) {
	self.expr.SetResult(value)
}

// Property

type IProperty interface {
	IExpr
	GetKey() string
	GetValue() IExpr
}

type Property struct {
	expr  *Expr
	key   string
	value IExpr
}

func NewProperty(key string, value IExpr, ln int, col int) *Property {
	return &Property{
		expr:  NewExpr(PropertyNode, ln, col),
		key:   key,
		value: value,
	}
}

func (self *Property) GetId() int {
	return self.expr.GetId()
}

func (self *Property) GetKind() NodeType {
	return self.expr.GetKind()
}

func (self *Property) GetKey() string {
	return self.key
}

func (self *Property) GetValue() IExpr {
	return self.value
}

func (self *Property) GetLn() int {
	return self.expr.GetLn()
}

func (self *Property) GetCol() int {
	return self.expr.GetCol()
}

func (self *Property) GetResult() IRuntimeVal {
	return self.expr.GetResult()
}

func (self *Property) SetResult(value IRuntimeVal) {
	self.expr.SetResult(value)
}

// ObjectLiteral

type IObjectLiteral interface {
	IExpr
	GetProperties() []IProperty
}

type ObjectLiteral struct {
	expr       *Expr
	properties []IProperty
}

func NewObjectLiteral(properties []IProperty) *ObjectLiteral {
	return &ObjectLiteral{
		expr:       NewExpr(ObjectLiteralNode, 0, 0),
		properties: properties,
	}
}

func (self *ObjectLiteral) GetId() int {
	return self.expr.GetId()
}

func (self *ObjectLiteral) GetKind() NodeType {
	return self.expr.GetKind()
}

func (self *ObjectLiteral) GetProperties() []IProperty {
	return self.properties
}

func (self *ObjectLiteral) GetLn() int {
	return self.expr.GetLn()
}

func (self *ObjectLiteral) GetCol() int {
	return self.expr.GetCol()
}

func (self *ObjectLiteral) GetResult() IRuntimeVal {
	return self.expr.GetResult()
}

func (self *ObjectLiteral) SetResult(value IRuntimeVal) {
	self.expr.SetResult(value)
}

// Identifier

type IIdentifier interface {
	IExpr
	GetSymbol() string
}

type Identifier struct {
	expr   *Expr
	symbol string
}

func (self *Identifier) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - id: %d, symbol: '%s'}", self.GetKind(), self.GetId(), self.GetSymbol())
	indentDepth--
	return str
}

func NewIdentifier(symbol string, ln int, col int) *Identifier {
	return &Identifier{
		expr:   NewExpr(IdentifierNode, ln, col),
		symbol: symbol,
	}
}

func (self *Identifier) GetId() int {
	return self.expr.GetId()
}

func (self *Identifier) GetKind() NodeType {
	return self.expr.GetKind()
}

func (self *Identifier) GetSymbol() string {
	return self.symbol
}

func (self *Identifier) GetLn() int {
	return self.expr.GetLn()
}

func (self *Identifier) GetCol() int {
	return self.expr.GetCol()
}

func (self *Identifier) GetResult() IRuntimeVal {
	return self.expr.GetResult()
}

func (self *Identifier) SetResult(value IRuntimeVal) {
	self.expr.SetResult(value)
}
