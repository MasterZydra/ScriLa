package bashAst

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
	GetOpType() NodeType
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

func (self *BinaryOpExpr) GetOpType() NodeType {
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
