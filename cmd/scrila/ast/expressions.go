package ast

import (
	"fmt"
)

// Expr

type IExpr interface {
	IStatement
}
type Expr struct {
	statement *Statement
}

func NewExpr(kind NodeType, ln int, col int) *Expr {
	return &Expr{statement: NewStatement(kind, ln, col)}
}

func NewEmptyExpr() *Expr {
	return NewExpr(ExprNode, 0, 0)
}

func (self *Expr) GetKind() NodeType {
	return self.statement.GetKind()
}

func (self *Expr) GetLn() int {
	return self.statement.GetLn()
}

func (self *Expr) GetCol() int {
	return self.statement.GetCol()
}

// AssignmentExpr

type IAssignmentExpr interface {
	IExpr
	GetAssigne() IExpr
	GetValue() IExpr
}

type AssignmentExpr struct {
	expr    *Expr
	assigne IExpr
	value   IExpr
}

func (self *AssignmentExpr) String() string {
	return fmt.Sprintf("&{%s %s %s}", self.GetKind(), self.GetAssigne(), self.GetValue())
}

func NewAssignmentExpr(assigne IExpr, value IExpr) *AssignmentExpr {
	return &AssignmentExpr{
		expr:    NewExpr(AssignmentExprNode, assigne.GetLn(), assigne.GetCol()),
		assigne: assigne,
		value:   value,
	}
}

func (self *AssignmentExpr) GetKind() NodeType {
	return self.expr.GetKind()
}

func (self *AssignmentExpr) GetAssigne() IExpr {
	return self.assigne
}

func (self *AssignmentExpr) GetValue() IExpr {
	return self.value
}

func (self *AssignmentExpr) GetLn() int {
	return self.expr.GetLn()
}

func (self *AssignmentExpr) GetCol() int {
	return self.expr.GetCol()
}

// BinaryExpr

type IBinaryExpr interface {
	IExpr
	GetLeft() IExpr
	GetRight() IExpr
	GetOperator() string
}

type BinaryExpr struct {
	expr     *Expr
	left     IExpr
	right    IExpr
	operator string
}

func (self *BinaryExpr) String() string {
	return fmt.Sprintf("&{%s %s %s %s}", self.GetKind(), self.GetLeft(), self.GetOperator(), self.GetRight())
}

func NewBinaryExpr(left IExpr, right IExpr, operator string, ln int, col int) *BinaryExpr {
	return &BinaryExpr{
		expr:     NewExpr(BinaryExprNode, ln, col),
		left:     left,
		right:    right,
		operator: operator,
	}
}

func (self *BinaryExpr) GetKind() NodeType {
	return self.expr.GetKind()
}

func (self *BinaryExpr) GetLeft() IExpr {
	return self.left
}

func (self *BinaryExpr) GetRight() IExpr {
	return self.right
}

func (self *BinaryExpr) GetOperator() string {
	return self.operator
}

func (self *BinaryExpr) GetLn() int {
	return self.expr.GetLn()
}

func (self *BinaryExpr) GetCol() int {
	return self.expr.GetCol()
}

// CallExpr

type ICallExpr interface {
	IExpr
	GetArgs() []IExpr
	GetCaller() IExpr
}

type CallExpr struct {
	expr   *Expr
	args   []IExpr
	caller IExpr
}

func (self *CallExpr) String() string {
	return fmt.Sprintf("&{%s %s %s}", self.GetKind(), self.GetCaller(), self.GetArgs())
}

func NewCallExpr(caller IExpr, args []IExpr) *CallExpr {
	return &CallExpr{
		expr:   NewExpr(CallExprNode, caller.GetLn(), caller.GetCol()),
		args:   args,
		caller: caller,
	}
}

func (self *CallExpr) GetKind() NodeType {
	return self.expr.GetKind()
}

func (self *CallExpr) GetCaller() IExpr {
	return self.caller
}

func (self *CallExpr) GetArgs() []IExpr {
	return self.args
}

func (self *CallExpr) GetLn() int {
	return self.expr.GetLn()
}

func (self *CallExpr) GetCol() int {
	return self.expr.GetCol()
}

// MemberExpr

// foo.bar()
// foo["bar"]() <- Computed
// foo[getBar()]() <- Computed

type IMemberExpr interface {
	IExpr
	GetObject() IExpr
	GetProperty() IExpr
	IsComputed() bool
}

type MemberExpr struct {
	expr       *Expr
	object     IExpr
	property   IExpr
	isComputed bool
}

func (self *MemberExpr) String() string {
	return fmt.Sprintf("&{%s %s %s %t}", self.GetKind(), self.GetObject(), self.GetProperty(), self.IsComputed())
}

func NewMemberExpr(object IExpr, property IExpr, isComputed bool) *MemberExpr {
	return &MemberExpr{
		expr:       NewExpr(MemberExprNode, 0, 0),
		object:     object,
		property:   property,
		isComputed: isComputed,
	}
}

func (self *MemberExpr) GetKind() NodeType {
	return self.expr.GetKind()
}

func (self *MemberExpr) GetObject() IExpr {
	return self.object
}

func (self *MemberExpr) GetProperty() IExpr {
	return self.property
}

func (self *MemberExpr) IsComputed() bool {
	return self.isComputed
}

func (self *MemberExpr) GetLn() int {
	return self.expr.GetLn()
}

func (self *MemberExpr) GetCol() int {
	return self.expr.GetCol()
}

// ReturnExpr

type IReturnExpr interface {
	IExpr
	GetValue() IExpr
}

type ReturnExpr struct {
	expr  *Expr
	value IExpr
}

func (self *ReturnExpr) String() string {
	return fmt.Sprintf("&{%s %s}", self.GetKind(), self.GetValue())
}

func NewReturnExpr(value IExpr, ln int, col int) *ReturnExpr {
	return &ReturnExpr{
		expr:  NewExpr(ReturnExprNode, ln, col),
		value: value,
	}
}

func (self *ReturnExpr) GetKind() NodeType {
	return self.expr.GetKind()
}

func (self *ReturnExpr) GetValue() IExpr {
	return self.value
}

func (self *ReturnExpr) GetLn() int {
	return self.expr.GetLn()
}

func (self *ReturnExpr) GetCol() int {
	return self.expr.GetCol()
}
