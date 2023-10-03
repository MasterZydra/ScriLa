package ast

import (
	"ScriLa/cmd/scrila/lexer"
	"fmt"
	"strconv"
)

// IntLiteral

type IIntLiteral interface {
	IExpr
	GetValue() int64
}

type IntLiteral struct {
	expr  *Expr
	value int64
}

func NewIntLiteral(token *lexer.Token) (*IntLiteral, error) {
	intLiteral := &IntLiteral{expr: NewExpr(IntLiteralNode, token.Ln, token.Col)}
	intValue, err := strconv.ParseInt(token.Value, 10, 64)
	if err != nil {
		return intLiteral, err
	}
	intLiteral.value = intValue
	return intLiteral, nil
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

// StrLiteral

type IStrLiteral interface {
	IExpr
	GetValue() string
}

type StrLiteral struct {
	expr  *Expr
	value string
}

func NewStrLiteral(token *lexer.Token) *StrLiteral {
	return &StrLiteral{
		expr:  NewExpr(StrLiteralNode, token.Ln, token.Col),
		value: token.Value,
	}
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

func (self *Property) String() string {
	return fmt.Sprintf("&{%s %s %s}", self.GetKind(), self.GetKey(), self.GetValue())
}

func NewProperty(key string, value IExpr, ln int, col int) *Property {
	return &Property{
		expr:  NewExpr(PropertyNode, ln, col),
		key:   key,
		value: value,
	}
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
		expr:       NewExpr(ObjectLiteralNode, 0, 0),
		properties: properties,
	}
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
