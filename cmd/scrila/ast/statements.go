package ast

import (
	"ScriLa/cmd/scrila/lexer"
	"fmt"
)

// Statement

type IStatement interface {
	GetKind() NodeType
	GetLn() int
	GetCol() int
}

type Statement struct {
	kind NodeType
	ln   int
	col  int
}

func NewStatement(kind NodeType, ln int, col int) *Statement {
	return &Statement{kind: kind, ln: ln, col: col}
}

func NewEmptyStatement() *Statement {
	return NewStatement(StatementNode, 0, 0)
}

func (self *Statement) GetKind() NodeType {
	return self.kind
}

func (self *Statement) GetLn() int {
	return self.ln
}

func (self *Statement) GetCol() int {
	return self.col
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
		statement: NewStatement(ProgramNode, 0, 0),
		Body:      make([]IStatement, 0),
	}
}

func (self *Program) GetKind() NodeType {
	return self.statement.GetKind()
}

func (self *Program) GetBody() []IStatement {
	return self.Body
}

func (self *Program) GetLn() int {
	return self.statement.GetLn()
}

func (self *Program) GetCol() int {
	return self.statement.GetCol()
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

func NewComment(token *lexer.Token) *Comment {
	return &Comment{
		statement: NewStatement(CommentNode, token.Ln, token.Col),
		comment:   token.Value,
	}
}

func (self *Comment) GetKind() NodeType {
	return self.statement.GetKind()
}

func (self *Comment) GetComment() string {
	return self.comment
}

func (self *Comment) GetLn() int {
	return self.statement.GetLn()
}

func (self *Comment) GetCol() int {
	return self.statement.GetCol()
}

// VarDeclaration

type IVarDeclaration interface {
	IStatement
	GetVarType() lexer.TokenType
	IsConstant() bool
	GetIdentifier() string
	GetValue() IExpr
}

type VarDeclaration struct {
	statement  *Statement
	varType    lexer.TokenType
	constant   bool
	identifier string
	value      IExpr
}

func (self *VarDeclaration) String() string {
	return fmt.Sprintf("&{%s %s %t %s %s}", self.GetKind(), self.GetVarType(), self.IsConstant(), self.GetIdentifier(), self.GetValue())
}

func NewVarDeclaration(varType lexer.TokenType, constant bool, identifier string, value IExpr, ln int, col int) *VarDeclaration {
	return &VarDeclaration{
		statement:  NewStatement(VarDeclarationNode, ln, col),
		varType:    varType,
		constant:   constant,
		identifier: identifier,
		value:      value,
	}
}

func (self *VarDeclaration) GetKind() NodeType {
	return self.statement.GetKind()
}

func (self *VarDeclaration) GetVarType() lexer.TokenType {
	return self.varType
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

func (self *VarDeclaration) GetLn() int {
	return self.statement.GetLn()
}

func (self *VarDeclaration) GetCol() int {
	return self.statement.GetCol()
}

// FunctionDeclaration

type Parameter struct {
	name      string
	paramType lexer.TokenType
}

func NewParameter(name string, paramType lexer.TokenType) *Parameter {
	return &Parameter{
		name:      name,
		paramType: paramType,
	}
}

func (self *Parameter) GetName() string {
	return self.name
}

func (self *Parameter) GetParamType() lexer.TokenType {
	return self.paramType
}

func (self *Parameter) String() string {
	return fmt.Sprintf("&{Parameter %s %s}", self.GetName(), self.GetParamType())
}

type IFunctionDeclaration interface {
	IStatement
	GetParameters() []*Parameter
	GetName() string
	GetBody() []IStatement
	GetReturnType() lexer.TokenType
}

type FunctionDeclaration struct {
	statement  *Statement
	parameters []*Parameter
	name       string
	body       []IStatement
	returnType lexer.TokenType
}

func (self *FunctionDeclaration) String() string {
	return fmt.Sprintf("&{%s %s %s %s}", self.GetKind(), self.GetName(), self.GetParameters(), self.GetBody())
}

func NewFunctionDeclaration(name string, parameters []*Parameter, body []IStatement, returnType lexer.TokenType, ln int, col int) *FunctionDeclaration {
	return &FunctionDeclaration{
		statement:  NewStatement(FunctionDeclarationNode, ln, col),
		name:       name,
		parameters: parameters,
		body:       body,
		returnType: returnType,
	}
}

func (self *FunctionDeclaration) GetKind() NodeType {
	return self.statement.GetKind()
}

func (self *FunctionDeclaration) GetName() string {
	return self.name
}

func (self *FunctionDeclaration) GetParameters() []*Parameter {
	return self.parameters
}

func (self *FunctionDeclaration) GetBody() []IStatement {
	return self.body
}

func (self *FunctionDeclaration) GetReturnType() lexer.TokenType {
	return self.returnType
}

func (self *FunctionDeclaration) GetLn() int {
	return self.statement.GetLn()
}

func (self *FunctionDeclaration) GetCol() int {
	return self.statement.GetCol()
}

// IfStatement

type IIfStatement interface {
	IStatement
	GetCondition() IExpr
	GetBody() []IStatement
}

type IfStatement struct {
	statement *Statement
	condition IExpr
	body      []IStatement
}

func (self *IfStatement) String() string {
	return fmt.Sprintf("&{%s %s %s}", self.GetKind(), self.GetCondition(), self.GetBody())
}

func NewIfStatement(condition IExpr, body []IStatement, ln int, col int) *IfStatement {
	return &IfStatement{
		statement: NewStatement(IfStatementNode, ln, col),
		condition: condition,
		body:      body,
	}
}

func (self *IfStatement) GetKind() NodeType {
	return self.statement.GetKind()
}

func (self *IfStatement) GetCondition() IExpr {
	return self.condition
}

func (self *IfStatement) GetBody() []IStatement {
	return self.body
}

func (self *IfStatement) GetLn() int {
	return self.statement.GetLn()
}

func (self *IfStatement) GetCol() int {
	return self.statement.GetCol()
}
