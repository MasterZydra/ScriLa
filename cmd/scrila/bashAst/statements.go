package bashAst

import "fmt"

// Statement

type IStatement interface {
	GetKind() NodeType
}

type Statement struct {
	kind NodeType
}

func (self *Statement) String() string {
	return fmt.Sprintf("{%s}", self.GetKind())
}

func NewStatement(kind NodeType) *Statement {
	return &Statement{kind: kind}
}

func (self *Statement) GetKind() NodeType {
	return self.kind
}

// BashStmt

type IBashStmt interface {
	IStrStmt
}

func NewBashStmt(bashCode string) *StrStmt {
	return NewStrStmt(BashStmtNode, bashCode)
}

// Comment

type IComment interface {
	IStrStmt
}

func NewComment(comment string) *StrStmt {
	return NewStrStmt(CommentNode, comment)
}

// FuncDeclaration

type IFuncDeclaration interface {
	IAppendBody
	AppendParams(param IFuncParameter)
	GetName() string
	GetParams() []IFuncParameter
	GetBody() []IStatement
	GetReturnType() NodeType
}

type FuncDeclaration struct {
	stmt       *Statement
	name       string
	params     []IFuncParameter
	body       []IStatement
	returnType NodeType
}

func (self *FuncDeclaration) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - name: '%s', returnType: '%s'", self.GetKind(), self.GetName(), self.GetReturnType())
	for i, param := range self.GetParams() {
		str += fmt.Sprintf("\n%sparam%d: %s", indent(), i, param)
	}
	if len(self.GetBody()) > 0 {
		str += fmt.Sprintf("\n%sbody:", indent())
		indentDepth++
		for _, stmt := range self.GetBody() {
			str += fmt.Sprintf("\n%s%s", indent(), stmt)
		}
		indentDepth--
	}
	indentDepth--
	return str + "}"
}

func NewFuncDeclaration(name string, returnType NodeType) *FuncDeclaration {
	return &FuncDeclaration{
		stmt:       NewStatement(FuncDeclarationNode),
		name:       name,
		params:     make([]IFuncParameter, 0),
		body:       make([]IStatement, 0),
		returnType: returnType,
	}
}

func (self *FuncDeclaration) AppendBody(stmt IStatement) {
	self.body = append(self.body, stmt)
}

func (self *FuncDeclaration) AppendParams(param IFuncParameter) {
	self.params = append(self.params, param)
}

func (self *FuncDeclaration) GetKind() NodeType {
	return self.stmt.GetKind()
}

func (self *FuncDeclaration) GetName() string {
	return self.name
}

func (self *FuncDeclaration) GetParams() []IFuncParameter {
	return self.params
}

func (self *FuncDeclaration) GetBody() []IStatement {
	return self.body
}

func (self *FuncDeclaration) GetReturnType() NodeType {
	return self.returnType
}

// FuncDeclaration - FuncParameter

type IFuncParameter interface {
	GetName() string
	GetType() NodeType
}

type FuncParameter struct {
	name      string
	paramType NodeType
}

func NewFuncParameter(name string, paramType NodeType) *FuncParameter {
	return &FuncParameter{name: name, paramType: paramType}
}

func (self *FuncParameter) GetName() string {
	return self.name
}

func (self *FuncParameter) GetType() NodeType {
	return self.paramType
}

// IfStmt

type IIfStmt interface {
	IAppendBody
	GetCondition() IStatement
	GetBody() []IStatement
	GetElse() IIfStmt
	SetElse(elseBlock IIfStmt)
}

type IfStmt struct {
	stmt      *Statement
	condition IStatement
	body      []IStatement
	elseBlock IIfStmt
}

func (self *IfStmt) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s", self.GetKind())
	if self.GetCondition() != nil {
		str += fmt.Sprintf("\n%scondition: %s,", indent(), self.GetCondition())
	}
	if len(self.GetBody()) > 0 {
		str += fmt.Sprintf("\n%sbody:", indent())
		indentDepth++
		for _, stmt := range self.GetBody() {
			str += fmt.Sprintf("\n%s%s", indent(), stmt)
		}
		indentDepth--
	}
	if self.GetElse() != nil {
		str += fmt.Sprintf("\n%selse: %s", indent(), self.GetElse())
	}
	indentDepth--
	return str + "}"
}

func NewIfStmt(condition IStatement) *IfStmt {
	return &IfStmt{
		stmt:      NewStatement(IfStmtNode),
		condition: condition,
		body:      make([]IStatement, 0),
		elseBlock: nil,
	}
}

func (self *IfStmt) AppendBody(stmt IStatement) {
	self.body = append(self.body, stmt)
}

func (self *IfStmt) GetKind() NodeType {
	return self.stmt.GetKind()
}

func (self *IfStmt) GetCondition() IStatement {
	return self.condition
}

func (self *IfStmt) GetBody() []IStatement {
	return self.body
}

func (self *IfStmt) GetElse() IIfStmt {
	return self.elseBlock
}

func (self *IfStmt) SetElse(elseBlock IIfStmt) {
	self.elseBlock = elseBlock
}

// Program

type IProgram interface {
	IStatement
	AppendNativeBody(stmt IStatement)
	AppendUserBody(stmt IStatement)
	GetNativeBody() []IStatement
	GetUserBody() []IStatement
}

type Program struct {
	stmt       *Statement
	nativeBody []IStatement
	userBody   []IStatement
}

func NewProgram() *Program {
	return &Program{
		stmt:       NewStatement(ProgramNode),
		nativeBody: make([]IStatement, 0),
		userBody:   make([]IStatement, 0),
	}
}

func (self *Program) AppendNativeBody(stmt IStatement) {
	self.nativeBody = append(self.nativeBody, stmt)
}

func (self *Program) AppendUserBody(stmt IStatement) {
	self.userBody = append(self.userBody, stmt)
}

func (self *Program) GetKind() NodeType {
	return self.stmt.GetKind()
}

func (self *Program) GetNativeBody() []IStatement {
	return self.nativeBody
}

func (self *Program) GetUserBody() []IStatement {
	return self.userBody
}

// WhileStmt

type IWhileStmt interface {
	IAppendBody
	GetCondition() IStatement
	GetBody() []IStatement
}

func NewWhileStmt(condition IStatement) *IfStmt {
	whileStmt := NewIfStmt(condition)
	whileStmt.stmt.kind = WhileStmtNode
	return whileStmt
}
