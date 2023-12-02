package bashAst

import "fmt"

// ArrayAssignmentExpr

type IArrayAssignmentExpr interface {
	IStatement
	GetVarname() IVarLiteral
	GetIndex() IStatement
	GetValue() IStatement
	IsDeclaration() bool
}

type ArrayAssignmentExpr struct {
	stmt          *Statement
	varname       IVarLiteral
	index         IStatement
	value         IStatement
	isDeclaration bool
}

func (self *ArrayAssignmentExpr) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - isDeclaration: %t,\n%svarName: %s\n%sindex: %s,\n%svalue: %s}", self.GetKind(), self.IsDeclaration(), indent(), self.GetVarname(), indent(), self.GetIndex(), indent(), self.GetValue())
	indentDepth--
	return str
}

func NewArrayAssignmentExpr(varname IVarLiteral, index IStatement, value IStatement, isDeclaration bool) *ArrayAssignmentExpr {
	return &ArrayAssignmentExpr{
		stmt:          NewStatement(ArrayAssignmentExprNode),
		varname:       varname,
		index:         index,
		value:         value,
		isDeclaration: isDeclaration,
	}
}

func (self *ArrayAssignmentExpr) GetKind() NodeType {
	return self.stmt.GetKind()
}

func (self *ArrayAssignmentExpr) GetVarname() IVarLiteral {
	return self.varname
}

func (self *ArrayAssignmentExpr) GetIndex() IStatement {
	return self.index
}

func (self *ArrayAssignmentExpr) GetValue() IStatement {
	return self.value
}

func (self *ArrayAssignmentExpr) IsDeclaration() bool {
	return self.isDeclaration
}

// AssignmentExpr

type IAssignmentExpr interface {
	IStatement
	GetVarname() IVarLiteral
	GetValue() IStatement
	IsDeclaration() bool
}

type AssignmentExpr struct {
	stmt          *Statement
	varname       IVarLiteral
	value         IStatement
	isDeclaration bool
}

func (self *AssignmentExpr) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - isDeclaration: %t,\n%svarName: %s\n%svalue: %s}", self.GetKind(), self.IsDeclaration(), indent(), self.GetVarname(), indent(), self.GetValue())
	indentDepth--
	return str
}

func NewAssignmentExpr(varname IVarLiteral, value IStatement, isDeclaration bool) *AssignmentExpr {
	return &AssignmentExpr{
		stmt:          NewStatement(AssignmentExprNode),
		varname:       varname,
		value:         value,
		isDeclaration: isDeclaration,
	}
}

func (self *AssignmentExpr) GetKind() NodeType {
	return self.stmt.GetKind()
}

func (self *AssignmentExpr) GetVarname() IVarLiteral {
	return self.varname
}

func (self *AssignmentExpr) GetValue() IStatement {
	return self.value
}

func (self *AssignmentExpr) IsDeclaration() bool {
	return self.isDeclaration
}

// BinaryOpExpr

type IBinaryOpExpr interface {
	IStatement
	GetDataType() NodeType
	GetLeft() IStatement
	GetRight() IStatement
	GetOperator() string
}

type BinaryOpExpr struct {
	stmt     *Statement
	opType   NodeType
	left     IStatement
	right    IStatement
	operator string
}

func (self *BinaryOpExpr) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - opType: '%s'\n%sleft: %s,\n%soperator: '%s',\n%sright: %s}", self.GetKind(), self.GetDataType(), indent(), self.GetLeft(), indent(), self.GetOperator(), indent(), self.GetRight())
	indentDepth--
	return str
}

func NewBinaryOpExpr(opType NodeType, left IStatement, right IStatement, operator string) *BinaryOpExpr {
	return &BinaryOpExpr{
		stmt:     NewStatement(BinaryOpExprNode),
		opType:   opType,
		left:     left,
		right:    right,
		operator: operator,
	}
}

func NewBinaryCompExpr(opType NodeType, left IStatement, right IStatement, operator string) *BinaryOpExpr {
	binCmpExpr := NewBinaryOpExpr(opType, left, right, operator)
	binCmpExpr.stmt.kind = BinaryCompExprNode
	return binCmpExpr
}

func (self *BinaryOpExpr) GetKind() NodeType {
	return self.stmt.GetKind()
}

func (self *BinaryOpExpr) GetDataType() NodeType {
	return self.opType
}

func (self *BinaryOpExpr) GetLeft() IStatement {
	return self.left
}

func (self *BinaryOpExpr) GetRight() IStatement {
	return self.right
}

func (self *BinaryOpExpr) GetOperator() string {
	return self.operator
}

// BreakExpr

func NewBreakExpr() *Statement {
	return NewStatement(BreakExprNode)
}

// CallExpr

type ICallExpr interface {
	IStatement
	GetFuncName() string
	GetArgs() []IStatement
}

type CallExpr struct {
	stmt     *Statement
	funcName string
	args     []IStatement
}

func (self *CallExpr) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - funcName: %s,", self.GetKind(), self.GetFuncName())
	for i, arg := range self.GetArgs() {
		str += fmt.Sprintf("\n%sarg%d: %s", indent(), i, arg)
	}
	indentDepth--
	return str + "}"
}

func NewCallExpr(funcName string, args []IStatement) *CallExpr {
	return &CallExpr{
		stmt:     NewStatement(CallExprNode),
		funcName: funcName,
		args:     args,
	}
}

func (self *CallExpr) GetKind() NodeType {
	return self.stmt.GetKind()
}

func (self *CallExpr) GetFuncName() string {
	return self.funcName
}

func (self *CallExpr) GetArgs() []IStatement {
	return self.args
}

// ContinueExpr

func NewContinueExpr() *Statement {
	return NewStatement(ContinueExprNode)
}

// MemberExpr

type IMemberExpr interface {
	IStatement
	GetVarname() IVarLiteral
	GetIndex() IStatement
}

type MemberExpr struct {
	stmt    *Statement
	varname IVarLiteral
	index   IStatement
}

func (self *MemberExpr) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - varname: %s,\n%sindex: %s}", self.GetKind(), self.GetVarname(), indent(), self.GetIndex())
	indentDepth--
	return str
}

func NewMemberExpr(varname IVarLiteral, index IStatement) *MemberExpr {
	return &MemberExpr{
		stmt:    NewStatement(MemberExprNode),
		varname: varname,
		index:   index,
	}
}

func (self *MemberExpr) GetKind() NodeType {
	return self.stmt.GetKind()
}

func (self *MemberExpr) GetVarname() IVarLiteral {
	return self.varname
}

func (self *MemberExpr) GetIndex() IStatement {
	return self.index
}

// ReturnExpr

func NewReturnExpr() *Statement {
	return NewStatement(ReturnExprNode)
}
