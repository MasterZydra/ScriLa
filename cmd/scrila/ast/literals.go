package ast

import "fmt"

// IntLiteral

type IIntLiteral interface {
	IExpr
	GetValue() int64
}

type IntLiteral struct {
	expr  *Expr
	value int64
}

func NewIntLiteral(value int64) *IntLiteral {
	return &IntLiteral{
		expr:  &Expr{statement: &Statement{kind: IntLiteralNode}},
		value: value,
	}
}

func (self *IntLiteral) GetKind() NodeType {
	return self.expr.statement.kind
}

func (self *IntLiteral) GetValue() int64 {
	return self.value
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

func NewStrLiteral(value string) *StrLiteral {
	return &StrLiteral{
		expr:  &Expr{statement: &Statement{kind: StrLiteralNode}},
		value: value,
	}
}

func (self *StrLiteral) GetKind() NodeType {
	return self.expr.statement.kind
}

func (self *StrLiteral) GetValue() string {
	return self.value
}

// Property

type IProperty interface {
	IExpr
	GetKey() string
	GetValue() IExpr
}

type Property struct {
	kind  NodeType
	key   string
	value IExpr
}

func (self *Property) String() string {
	return fmt.Sprintf("&{%s %s %s}", self.GetKind(), self.GetKey(), self.GetValue())
}

func NewProperty(key string, value IExpr) *Property {
	return &Property{
		kind:  PropertyNode,
		key:   key,
		value: value,
	}
}

func (self *Property) GetKind() NodeType {
	return self.kind
}

func (self *Property) GetKey() string {
	return self.key
}

func (self *Property) GetValue() IExpr {
	return self.value
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

func (self *ObjectLiteral) String() string {
	return fmt.Sprintf("&{%s %s}", self.GetKind(), self.GetProperties())
}

func NewObjectLiteral(properties []IProperty) *ObjectLiteral {
	return &ObjectLiteral{
		expr:       &Expr{statement: &Statement{kind: ObjectLiteralNode}},
		properties: properties,
	}
}

func (self *ObjectLiteral) GetKind() NodeType {
	return self.expr.statement.kind
}

func (self *ObjectLiteral) GetProperties() []IProperty {
	return self.properties
}
