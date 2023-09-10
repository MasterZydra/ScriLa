package ast

import "fmt"

type NodeType string

const (
	// Statements
	ProgramNode        NodeType = "Program"
	VarDeclarationNode NodeType = "VarDeclaration"

	// Expressions
	BinaryExprNode          NodeType = "BinaryExpr"
	CallExprNode            NodeType = "CallExpr"
	FunctionDeclarationNode NodeType = "FunctionDeclaration"
	IdentifierNode          NodeType = "Identifier"
	IntLiteralNode          NodeType = "IntLiteral"
	UnaryExprNode           NodeType = "UnaryExpr"
)

type IStatement interface {
	GetKind() NodeType
}

type Statement struct {
	kind NodeType
}

func (self *Statement) GetKind() NodeType {
	return self.kind
}

type IProgram interface {
	IStatement
	GetBody() []IStatement
}

type Program struct {
	kind NodeType
	Body []IStatement
}

func (self *Program) String() string {
	return fmt.Sprintf("&{%s %s}", self.GetKind(), self.GetBody())
}

func NewProgram() *Program {
	return &Program{
		kind: ProgramNode,
		Body: make([]IStatement, 0),
	}
}

func (self *Program) GetKind() NodeType {
	return self.kind
}

func (self *Program) GetBody() []IStatement {
	return self.Body
}

type IVarDeclaration interface {
	IStatement
	IsConstant() bool
	GetIdentifier() string
	GetValue() IExpr
}

type VarDeclaration struct {
	kind       NodeType
	constant   bool
	identifier string
	value      IExpr
}

func (self *VarDeclaration) String() string {
	return fmt.Sprintf("&{%s %t %s %s}", self.GetKind(), self.IsConstant(), self.GetIdentifier(), self.GetValue())
}

func NewVarDeclaration(constant bool, identifier string, value IExpr) *VarDeclaration {
	return &VarDeclaration{
		kind:       VarDeclarationNode,
		constant:   constant,
		identifier: identifier,
		value:      value,
	}
}

func (self *VarDeclaration) GetKind() NodeType {
	return self.kind
}

func (self *VarDeclaration) IsConstant() bool {
	return self.constant
}

func (self *VarDeclaration) GetIdentifier() string {
	return self.identifier
}

func (self *VarDeclaration) GetValue() IExpr {
	return self.value
}

type IExpr interface {
	IStatement
}
type Expr struct {
	kind NodeType
}

func (self *Expr) GetKind() NodeType {
	return self.kind
}

type IBinaryExpr interface {
	IExpr
	GetLeft() IExpr
	GetRight() IExpr
	GetOperator() string
}

type BinaryExpr struct {
	kind     NodeType
	left     IExpr
	right    IExpr
	operator string
}

func (self *BinaryExpr) String() string {
	return fmt.Sprintf("&{%s %s %s %s}", self.GetKind(), self.GetLeft(), self.GetOperator(), self.GetRight())
}

func NewBinaryExpr(left IExpr, right IExpr, operator string) *BinaryExpr {
	return &BinaryExpr{
		kind:     BinaryExprNode,
		left:     left,
		right:    right,
		operator: operator,
	}
}

func (self *BinaryExpr) GetKind() NodeType {
	return self.kind
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

type IIdentifier interface {
	IExpr
	GetSymbol() string
}

type Identifier struct {
	kind   NodeType
	symbol string
}

func NewIdentifier(symbol string) *Identifier {
	return &Identifier{
		kind:   IdentifierNode,
		symbol: symbol,
	}
}

func (self *Identifier) GetKind() NodeType {
	return self.kind
}

func (self *Identifier) GetSymbol() string {
	return self.symbol
}

type IIntLiteral interface {
	IExpr
	GetValue() int64
}

type IntLiteral struct {
	kind  NodeType
	value int64
}

func NewIntLiteral(value int64) *IntLiteral {
	return &IntLiteral{
		kind:  IntLiteralNode,
		value: value,
	}
}

func (self *IntLiteral) GetKind() NodeType {
	return self.kind
}

func (self *IntLiteral) GetValue() int64 {
	return self.value
}
