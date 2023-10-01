package ast

import "fmt"

// Statement

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

// Program

type IProgram interface {
	IStatement
	GetBody() []IStatement
}

type Program struct {
	statement *Statement
	Body      []IStatement
}

func (self *Program) String() string {
	return fmt.Sprintf("&{%s %s}", self.GetKind(), self.GetBody())
}

func NewProgram() *Program {
	return &Program{
		statement: &Statement{kind: ProgramNode},
		Body:      make([]IStatement, 0),
	}
}

func (self *Program) GetKind() NodeType {
	return self.statement.kind
}

func (self *Program) GetBody() []IStatement {
	return self.Body
}

// Comment

type IComment interface {
	IStatement
	GetComment() string
}

type Comment struct {
	statement *Statement
	comment   string
}

func (self *Comment) String() string {
	return fmt.Sprintf("&{%s '%s'}", self.GetKind(), self.GetComment())
}

func NewComment(comment string) *Comment {
	return &Comment{
		statement: &Statement{kind: CommentNode},
		comment:   comment,
	}
}

func (self *Comment) GetKind() NodeType {
	return self.statement.kind
}

func (self *Comment) GetComment() string {
	return self.comment
}

// VarDeclaration

type IVarDeclaration interface {
	IStatement
	IsConstant() bool
	GetIdentifier() string
	GetValue() IExpr
}

type VarDeclaration struct {
	statement  *Statement
	constant   bool
	identifier string
	value      IExpr
}

func (self *VarDeclaration) String() string {
	return fmt.Sprintf("&{%s %t %s %s}", self.GetKind(), self.IsConstant(), self.GetIdentifier(), self.GetValue())
}

func NewVarDeclaration(constant bool, identifier string, value IExpr) *VarDeclaration {
	return &VarDeclaration{
		statement:  &Statement{kind: VarDeclarationNode},
		constant:   constant,
		identifier: identifier,
		value:      value,
	}
}

func (self *VarDeclaration) GetKind() NodeType {
	return self.statement.kind
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

// FunctionDeclaration

type IFunctionDeclaration interface {
	IStatement
	GetParameters() []string
	GetName() string
	GetBody() []IStatement
	// TODO Return type
}

type FunctionDeclaration struct {
	statement  *Statement
	parameters []string
	name       string
	body       []IStatement
}

func (self *FunctionDeclaration) String() string {
	return fmt.Sprintf("&{%s %s %s %s}", self.GetKind(), self.GetName(), self.GetParameters(), self.GetBody())
}

func NewFunctionDeclaration(name string, parameters []string, body []IStatement) *FunctionDeclaration {
	return &FunctionDeclaration{
		statement:  &Statement{kind: FunctionDeclarationNode},
		name:       name,
		parameters: parameters,
		body:       body,
	}
}

func (self *FunctionDeclaration) GetKind() NodeType {
	return self.statement.kind
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
