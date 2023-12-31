package scrilaAst

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

func (self *Expr) String() string {
	return fmt.Sprintf("{%s - id: %d}", self.GetKind(), self.GetId())
}

func NewExpr(kind NodeType, ln int, col int) *Expr {
	return &Expr{statement: NewStatement(kind, ln, col)}
}

func NewEmptyExpr() *Expr {
	return NewExpr(ExprNode, 0, 0)
}

func (self *Expr) GetId() int {
	return self.statement.GetId()
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

func (self *Expr) GetResult() IRuntimeVal {
	return self.statement.GetResult()
}

func (self *Expr) SetResult(value IRuntimeVal) {
	self.statement.SetResult(value)
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
	indentDepth++
	str := fmt.Sprintf("{%s - id: %d, varName: '%s',\n%svalue: %s}", self.GetKind(), self.GetId(), self.GetAssigne(), indent(), self.GetValue())
	indentDepth--
	return str
}

func NewAssignmentExpr(assigne IExpr, value IExpr) *AssignmentExpr {
	return &AssignmentExpr{
		expr:    NewExpr(AssignmentExprNode, assigne.GetLn(), assigne.GetCol()),
		assigne: assigne,
		value:   value,
	}
}

func (self *AssignmentExpr) GetId() int {
	return self.expr.GetId()
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

func (self *AssignmentExpr) GetResult() IRuntimeVal {
	return self.expr.GetResult()
}

func (self *AssignmentExpr) SetResult(value IRuntimeVal) {
	self.expr.SetResult(value)
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
	indentDepth++
	str := fmt.Sprintf("{%s - id: %d,\n%sleft: %s,\n%soperator: '%s',\n%sright: %s}", self.GetKind(), self.GetId(), indent(), self.GetLeft(), indent(), self.GetOperator(), indent(), self.GetRight())
	indentDepth--
	return str
}

func NewBinaryExpr(left IExpr, right IExpr, operator string, ln int, col int) *BinaryExpr {
	return &BinaryExpr{
		expr:     NewExpr(BinaryExprNode, ln, col),
		left:     left,
		right:    right,
		operator: operator,
	}
}

func (self *BinaryExpr) GetId() int {
	return self.expr.GetId()
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

func (self *BinaryExpr) GetResult() IRuntimeVal {
	return self.expr.GetResult()
}

func (self *BinaryExpr) SetResult(value IRuntimeVal) {
	self.expr.SetResult(value)
}

// BreakExpr

type IBreakExpr interface {
	IExpr
}

func NewBreakExpr(ln int, col int) *Expr {
	return NewExpr(BreakExprNode, ln, col)
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
	indentDepth++
	str := fmt.Sprintf("{%s - id: %d,\n%sfuncName: %s,", self.GetKind(), self.GetId(), indent(), self.GetCaller())
	for i, arg := range self.GetArgs() {
		str += fmt.Sprintf("\n%sarg%d: %s", indent(), i, arg)
	}
	indentDepth--
	return str + "}"
}

func NewCallExpr(caller IExpr, args []IExpr) *CallExpr {
	return &CallExpr{
		expr:   NewExpr(CallExprNode, caller.GetLn(), caller.GetCol()),
		args:   args,
		caller: caller,
	}
}

func (self *CallExpr) GetId() int {
	return self.expr.GetId()
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

func (self *CallExpr) GetResult() IRuntimeVal {
	return self.expr.GetResult()
}

func (self *CallExpr) SetResult(value IRuntimeVal) {
	self.expr.SetResult(value)
}

// ContinueExpr

type IContinueExpr interface {
	IExpr
}

func NewContinueExpr(ln int, col int) *Expr {
	return NewExpr(ContinueExprNode, ln, col)
}

// MemberExpr

// foo.bar()
// foo["bar"]() <- Computed
// foo[getBar()]() <- Computed

type IMemberExpr interface {
	IExpr
	GetObject() IExpr
	GetProperty() IExpr
	IsEmpty() bool
}

type MemberExpr struct {
	expr     *Expr
	object   IExpr
	property IExpr
	isEmpty  bool
}

func (self *MemberExpr) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - id: %d, isEmpty: %t,\n%sobject: %s,\n%sproperty: %s}", self.GetKind(), self.GetId(), self.IsEmpty(), indent(), self.GetObject(), indent(), self.GetProperty())
	indentDepth--
	return str
}

func NewMemberExpr(object IExpr, property IExpr, isEmpty bool) *MemberExpr {
	return &MemberExpr{
		expr:     NewExpr(MemberExprNode, 0, 0),
		object:   object,
		property: property,
		isEmpty:  isEmpty,
	}
}

func (self *MemberExpr) GetId() int {
	return self.expr.GetId()
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

func (self *MemberExpr) IsEmpty() bool {
	return self.isEmpty
}

func (self *MemberExpr) GetLn() int {
	return self.expr.GetLn()
}

func (self *MemberExpr) GetCol() int {
	return self.expr.GetCol()
}

func (self *MemberExpr) GetResult() IRuntimeVal {
	return self.expr.GetResult()
}

func (self *MemberExpr) SetResult(value IRuntimeVal) {
	self.expr.SetResult(value)
}

// ReturnExpr

type IReturnExpr interface {
	IExpr
	GetValue() IExpr
	IsEmpty() bool
}

type ReturnExpr struct {
	expr    *Expr
	value   IExpr
	isEmpty bool
}

func (self *ReturnExpr) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - id: %d,", self.GetKind(), self.GetId())
	if !self.IsEmpty() {
		str += fmt.Sprintf("\n%svalue: %s", indent(), self.GetValue())
	}
	indentDepth--
	return str + "}"
}

func NewReturnExpr(value IExpr, isEmpty bool, ln int, col int) *ReturnExpr {
	return &ReturnExpr{
		expr:    NewExpr(ReturnExprNode, ln, col),
		value:   value,
		isEmpty: isEmpty,
	}
}

func (self *ReturnExpr) GetId() int {
	return self.expr.GetId()
}

func (self *ReturnExpr) GetKind() NodeType {
	return self.expr.GetKind()
}

func (self *ReturnExpr) GetValue() IExpr {
	return self.value
}

func (self *ReturnExpr) IsEmpty() bool {
	return self.isEmpty
}

func (self *ReturnExpr) GetLn() int {
	return self.expr.GetLn()
}

func (self *ReturnExpr) GetCol() int {
	return self.expr.GetCol()
}

func (self *ReturnExpr) GetResult() IRuntimeVal {
	return self.expr.GetResult()
}

func (self *ReturnExpr) SetResult(value IRuntimeVal) {
	self.expr.SetResult(value)
}
