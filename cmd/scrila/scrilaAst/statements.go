package scrilaAst

import (
	"fmt"
)

var currElemId int = 0

func getNextElemId() int {
	currElemId++
	return currElemId
}

// Statement

type IStatement interface {
	GetId() int
	GetKind() NodeType
	GetLn() int
	GetCol() int
	GetResult() IRuntimeVal
	SetResult(value IRuntimeVal)
}

type Statement struct {
	id     int
	kind   NodeType
	ln     int
	col    int
	result IRuntimeVal
}

func NewStatement(kind NodeType, ln int, col int) *Statement {
	return &Statement{id: getNextElemId(), kind: kind, ln: ln, col: col}
}

func NewEmptyStatement() *Statement {
	return NewStatement(StatementNode, 0, 0)
}

func (self *Statement) GetId() int {
	return self.id
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

func (self *Statement) GetResult() IRuntimeVal {
	return self.result
}

func (self *Statement) SetResult(value IRuntimeVal) {
	self.result = value
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

func (self *Program) GetId() int {
	return self.statement.GetId()
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

func (self *Program) GetResult() IRuntimeVal {
	return self.statement.GetResult()
}

func (self *Program) SetResult(value IRuntimeVal) {
	self.statement.SetResult(value)
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
	return fmt.Sprintf("{%s - id: %d, comment: '%s'}", self.GetKind(), self.GetId(), self.GetComment())
}

func NewComment(comment string, ln int, col int) *Comment {
	return &Comment{
		statement: NewStatement(CommentNode, ln, col),
		comment:   comment,
	}
}

func (self *Comment) GetId() int {
	return self.statement.GetId()
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

func (self *Comment) GetResult() IRuntimeVal {
	return self.statement.GetResult()
}

func (self *Comment) SetResult(value IRuntimeVal) {
	self.statement.SetResult(value)
}

// VarDeclaration

type IVarDeclaration interface {
	IStatement
	GetDataType() NodeType
	IsConstant() bool
	GetIdentifier() string
	GetValue() IExpr
}

type VarDeclaration struct {
	statement  *Statement
	varType    NodeType
	constant   bool
	identifier string
	value      IExpr
}

func (self *VarDeclaration) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - id: %d, varName: '%s', varType: '%s', isConstant: %t,\n%svalue: %s}", self.GetKind(), self.GetId(), self.GetIdentifier(), self.GetDataType(), self.IsConstant(), indent(), self.GetValue())
	indentDepth--
	return str
}

func NewVarDeclaration(varType NodeType, constant bool, identifier string, value IExpr, ln int, col int) *VarDeclaration {
	return &VarDeclaration{
		statement:  NewStatement(VarDeclarationNode, ln, col),
		varType:    varType,
		constant:   constant,
		identifier: identifier,
		value:      value,
	}
}

func (self *VarDeclaration) GetId() int {
	return self.statement.GetId()
}

func (self *VarDeclaration) GetKind() NodeType {
	return self.statement.GetKind()
}

func (self *VarDeclaration) GetDataType() NodeType {
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

func (self *VarDeclaration) GetResult() IRuntimeVal {
	return self.statement.GetResult()
}

func (self *VarDeclaration) SetResult(value IRuntimeVal) {
	self.statement.SetResult(value)
}

// ForStatement

type IForStatement interface {
	IStatement
	GetIndexVarType() NodeType
	GetIndex() IIdentifier
	GetArray() IExpr
	GetBody() []IStatement
}

type ForStatement struct {
	statement    *Statement
	indexVarType NodeType
	index        IIdentifier
	array        IExpr
	body         []IStatement
}

func (self *ForStatement) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - id: %d, index var type: %s,\n%sindex: %s\n%sarray: %s", self.GetKind(), self.GetId(), self.GetIndexVarType(), indent(), self.GetIndex(), indent(), self.GetArray())
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

func NewForStatement(indexVarType NodeType, index IIdentifier, array IExpr, body []IStatement, ln int, col int) *ForStatement {
	return &ForStatement{
		statement:    NewStatement(ForStatementNode, ln, col),
		indexVarType: indexVarType,
		index:        index,
		array:        array,
		body:         body,
	}
}

func (self *ForStatement) GetId() int {
	return self.statement.GetId()
}

func (self *ForStatement) GetKind() NodeType {
	return self.statement.GetKind()
}

func (self *ForStatement) GetIndexVarType() NodeType {
	return self.indexVarType
}

func (self *ForStatement) GetIndex() IIdentifier {
	return self.index
}

func (self *ForStatement) GetArray() IExpr {
	return self.array
}

func (self *ForStatement) GetBody() []IStatement {
	return self.body
}

func (self *ForStatement) GetLn() int {
	return self.statement.GetLn()
}

func (self *ForStatement) GetCol() int {
	return self.statement.GetCol()
}

func (self *ForStatement) GetResult() IRuntimeVal {
	return self.statement.GetResult()
}

func (self *ForStatement) SetResult(value IRuntimeVal) {
	self.statement.SetResult(value)
}

// FunctionDeclaration

type Parameter struct {
	name      string
	paramType NodeType
}

func NewParameter(name string, paramType NodeType) *Parameter {
	return &Parameter{
		name:      name,
		paramType: paramType,
	}
}

func (self *Parameter) GetName() string {
	return self.name
}

func (self *Parameter) GetParamType() NodeType {
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
	GetReturnType() NodeType
}

type FunctionDeclaration struct {
	statement  *Statement
	parameters []*Parameter
	name       string
	body       []IStatement
	returnType NodeType
}

func (self *FunctionDeclaration) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - id: %d, name: '%s',", self.GetKind(), self.GetId(), self.GetName())
	for i, param := range self.GetParameters() {
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

func NewFunctionDeclaration(name string, parameters []*Parameter, body []IStatement, returnType NodeType, ln int, col int) *FunctionDeclaration {
	return &FunctionDeclaration{
		statement:  NewStatement(FunctionDeclarationNode, ln, col),
		name:       name,
		parameters: parameters,
		body:       body,
		returnType: returnType,
	}
}

func (self *FunctionDeclaration) GetId() int {
	return self.statement.GetId()
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

func (self *FunctionDeclaration) GetReturnType() NodeType {
	return self.returnType
}

func (self *FunctionDeclaration) GetLn() int {
	return self.statement.GetLn()
}

func (self *FunctionDeclaration) GetCol() int {
	return self.statement.GetCol()
}

func (self *FunctionDeclaration) GetResult() IRuntimeVal {
	return self.statement.GetResult()
}

func (self *FunctionDeclaration) SetResult(value IRuntimeVal) {
	self.statement.SetResult(value)
}

// IfStatement

type IIfStatement interface {
	IStatement
	GetCondition() IExpr
	GetBody() []IStatement
	GetElse() IIfStatement
}

type IfStatement struct {
	statement *Statement
	condition IExpr
	body      []IStatement
	elseBlock IIfStatement
}

func (self *IfStatement) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - id: %d,", self.GetKind(), self.GetId())
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

func NewIfStatement(condition IExpr, body []IStatement, elseBlock IIfStatement, ln int, col int) *IfStatement {
	return &IfStatement{
		statement: NewStatement(IfStatementNode, ln, col),
		condition: condition,
		body:      body,
		elseBlock: elseBlock,
	}
}

func (self *IfStatement) GetId() int {
	return self.statement.GetId()
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

func (self *IfStatement) GetElse() IIfStatement {
	return self.elseBlock
}

func (self *IfStatement) GetLn() int {
	return self.statement.GetLn()
}

func (self *IfStatement) GetCol() int {
	return self.statement.GetCol()
}

func (self *IfStatement) GetResult() IRuntimeVal {
	return self.statement.GetResult()
}

func (self *IfStatement) SetResult(value IRuntimeVal) {
	self.statement.SetResult(value)
}

// WhileStatement

type IWhileStatement interface {
	IStatement
	GetCondition() IExpr
	GetBody() []IStatement
}

type WhileStatement struct {
	statement *Statement
	condition IExpr
	body      []IStatement
}

func (self *WhileStatement) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - id: %d,\n%scondition: %s,", self.GetKind(), self.GetId(), indent(), self.GetCondition())
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

func NewWhileStatement(condition IExpr, body []IStatement, ln int, col int) *WhileStatement {
	return &WhileStatement{
		statement: NewStatement(WhileStatementNode, ln, col),
		condition: condition,
		body:      body,
	}
}

func (self *WhileStatement) GetId() int {
	return self.statement.GetId()
}

func (self *WhileStatement) GetKind() NodeType {
	return self.statement.GetKind()
}

func (self *WhileStatement) GetCondition() IExpr {
	return self.condition
}

func (self *WhileStatement) GetBody() []IStatement {
	return self.body
}

func (self *WhileStatement) GetLn() int {
	return self.statement.GetLn()
}

func (self *WhileStatement) GetCol() int {
	return self.statement.GetCol()
}

func (self *WhileStatement) GetResult() IRuntimeVal {
	return self.statement.GetResult()
}

func (self *WhileStatement) SetResult(value IRuntimeVal) {
	self.statement.SetResult(value)
}
