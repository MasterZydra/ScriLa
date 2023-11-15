package bashAst

import "fmt"

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

// ReturnExpr

func NewReturnExpr() *Statement {
	return NewStatement(ReturnExprNode)
}
