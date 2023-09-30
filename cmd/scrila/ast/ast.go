package ast

import "fmt"

type NodeType string

const (
	// Statements
	StatementNode           NodeType = "Statement"
	CommentNode             NodeType = "Comment"
	ProgramNode             NodeType = "Program"
	VarDeclarationNode      NodeType = "VarDeclaration"
	FunctionDeclarationNode NodeType = "FunctionDeclaration"

	// Expressions
	ExprNode           NodeType = "Expr"
	AssignmentExprNode NodeType = "AssignmentExpr"
	BinaryExprNode     NodeType = "BinaryExpr"
	UnaryExprNode      NodeType = "UnaryExpr"
	CallExprNode       NodeType = "CallExpr"
	MemberExprNode     NodeType = "MemberExpr"

	// Literals
	PropertyNode      NodeType = "Property"
	ObjectLiteralNode NodeType = "ObjectLiteral"
	IdentifierNode    NodeType = "Identifier"
	IntLiteralNode    NodeType = "IntLiteral"
	StrLiteralNode    NodeType = "StrLiteral"
)

type IStatement interface {
	GetKind() NodeType
}

type Statement struct {
	kind NodeType
}

func NewStatement() *Statement {
	return &Statement{
		kind: StatementNode,
	}
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

type IComment interface {
	IStatement
	GetComment() string
}

type Comment struct {
	kind    NodeType
	comment string
}

func (self *Comment) String() string {
	return fmt.Sprintf("&{%s '%s'}", self.GetKind(), self.GetComment())
}

func NewComment(comment string) *Comment {
	return &Comment{
		kind:    CommentNode,
		comment: comment,
	}
}

func (self *Comment) GetKind() NodeType {
	return self.kind
}

func (self *Comment) GetComment() string {
	return self.comment
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

type IFunctionDeclaration interface {
	IStatement
	GetParameters() []string
	GetName() string
	GetBody() []IStatement
	// TODO Return type
}

type FunctionDeclaration struct {
	kind       NodeType
	parameters []string
	name       string
	body       []IStatement
}

func (self *FunctionDeclaration) String() string {
	return fmt.Sprintf("&{%s %s %s %s}", self.GetKind(), self.GetName(), self.GetParameters(), self.GetBody())
}

func NewFunctionDeclaration(name string, parameters []string, body []IStatement) *FunctionDeclaration {
	return &FunctionDeclaration{
		kind:       FunctionDeclarationNode,
		name:       name,
		parameters: parameters,
		body:       body,
	}
}

func (self *FunctionDeclaration) GetKind() NodeType {
	return self.kind
}

func (self *FunctionDeclaration) GetName() string {
	return self.name
}

func (self *FunctionDeclaration) GetParameters() []string {
	return self.parameters
}

func (self *FunctionDeclaration) GetBody() []IStatement {
	return self.body
}

type IExpr interface {
	IStatement
}
type Expr struct {
	kind NodeType
}

func NewExpr() *Expr {
	return &Expr{
		kind: ExprNode,
	}
}

func (self *Expr) GetKind() NodeType {
	return self.kind
}

type IAssignmentExpr interface {
	IExpr
	GetAssigne() IExpr
	GetValue() IExpr
}

type AssignmentExpr struct {
	kind    NodeType
	assigne IExpr
	value   IExpr
}

func (self *AssignmentExpr) String() string {
	return fmt.Sprintf("&{%s %s %s}", self.GetKind(), self.GetAssigne(), self.GetValue())
}

func NewAssignmentExpr(assigne IExpr, value IExpr) *AssignmentExpr {
	return &AssignmentExpr{
		kind:    AssignmentExprNode,
		assigne: assigne,
		value:   value,
	}
}

func (self *AssignmentExpr) GetKind() NodeType {
	return self.kind
}

func (self *AssignmentExpr) GetAssigne() IExpr {
	return self.assigne
}

func (self *AssignmentExpr) GetValue() IExpr {
	return self.value
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

type ICallExpr interface {
	IExpr
	GetArgs() []IExpr
	GetCaller() IExpr
}

type CallExpr struct {
	kind   NodeType
	args   []IExpr
	caller IExpr
}

func (self *CallExpr) String() string {
	return fmt.Sprintf("&{%s %s %s}", self.GetKind(), self.GetCaller(), self.GetArgs())
}

func NewCallExpr(caller IExpr, args []IExpr) *CallExpr {
	return &CallExpr{
		kind:   CallExprNode,
		args:   args,
		caller: caller,
	}
}

func (self *CallExpr) GetKind() NodeType {
	return self.kind
}

func (self *CallExpr) GetCaller() IExpr {
	return self.caller
}

func (self *CallExpr) GetArgs() []IExpr {
	return self.args
}

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
	kind       NodeType
	object     IExpr
	property   IExpr
	isComputed bool
}

func (self *MemberExpr) String() string {
	return fmt.Sprintf("&{%s %s %s %t}", self.GetKind(), self.GetObject(), self.GetProperty(), self.IsComputed())
}

func NewMemberExpr(object IExpr, property IExpr, isComputed bool) *MemberExpr {
	return &MemberExpr{
		kind:       MemberExprNode,
		object:     object,
		property:   property,
		isComputed: isComputed,
	}
}

func (self *MemberExpr) GetKind() NodeType {
	return self.kind
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

type IStrLiteral interface {
	IExpr
	GetValue() string
}

type StrLiteral struct {
	kind  NodeType
	value string
}

func NewStrLiteral(value string) *StrLiteral {
	return &StrLiteral{
		kind:  StrLiteralNode,
		value: value,
	}
}

func (self *StrLiteral) GetKind() NodeType {
	return self.kind
}

func (self *StrLiteral) GetValue() string {
	return self.value
}

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

type IObjectLiteral interface {
	IExpr
	GetProperties() []IProperty
}

type ObjectLiteral struct {
	kind       NodeType
	properties []IProperty
}

func (self *ObjectLiteral) String() string {
	return fmt.Sprintf("&{%s %s}", self.GetKind(), self.GetProperties())
}

func NewObjectLiteral(properties []IProperty) *ObjectLiteral {
	return &ObjectLiteral{
		kind:       ObjectLiteralNode,
		properties: properties,
	}
}

func (self *ObjectLiteral) GetKind() NodeType {
	return self.kind
}

func (self *ObjectLiteral) GetProperties() []IProperty {
	return self.properties
}
