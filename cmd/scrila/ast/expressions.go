package ast

import "fmt"

// Expr

type IExpr interface {
	IStatement
}
type Expr struct {
	statement *Statement
}

func NewExpr() *Expr {
	return &Expr{
		statement: &Statement{kind: ExprNode},
	}
}

func (self *Expr) GetKind() NodeType {
	return self.statement.kind
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
		expr:    &Expr{statement: &Statement{kind: AssignmentExprNode}},
		assigne: assigne,
		value:   value,
	}
}

func (self *AssignmentExpr) GetKind() NodeType {
	return self.expr.statement.kind
}

func (self *AssignmentExpr) GetAssigne() IExpr {
	return self.assigne
}

func (self *AssignmentExpr) GetValue() IExpr {
	return self.value
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

func NewBinaryExpr(left IExpr, right IExpr, operator string) *BinaryExpr {
	return &BinaryExpr{
		expr:     &Expr{statement: &Statement{kind: BinaryExprNode}},
		left:     left,
		right:    right,
		operator: operator,
	}
}

func (self *BinaryExpr) GetKind() NodeType {
	return self.expr.statement.kind
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
		expr:   &Expr{statement: &Statement{kind: CallExprNode}},
		args:   args,
		caller: caller,
	}
}

func (self *CallExpr) GetKind() NodeType {
	return self.expr.statement.kind
}

func (self *CallExpr) GetCaller() IExpr {
	return self.caller
}

func (self *CallExpr) GetArgs() []IExpr {
	return self.args
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
		expr:       &Expr{statement: &Statement{kind: MemberExprNode}},
		object:     object,
		property:   property,
		isComputed: isComputed,
	}
}

func (self *MemberExpr) GetKind() NodeType {
	return self.expr.statement.kind
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

// Identifier

type IIdentifier interface {
	IExpr
	GetSymbol() string
}

type Identifier struct {
	expr   *Expr
	symbol string
}

func NewIdentifier(symbol string) *Identifier {
	return &Identifier{
		expr:   &Expr{statement: &Statement{kind: IdentifierNode}},
		symbol: symbol,
	}
}

func (self *Identifier) GetKind() NodeType {
	return self.expr.statement.kind
}

func (self *Identifier) GetSymbol() string {
	return self.symbol
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

func NewReturnExpr(value IExpr) *ReturnExpr {
	return &ReturnExpr{
		expr:  &Expr{statement: &Statement{kind: ReturnExprNode}},
		value: value,
	}
}

func (self *ReturnExpr) GetKind() NodeType {
	return self.expr.statement.kind
}

func (self *ReturnExpr) GetValue() IExpr {
	return self.value
}
