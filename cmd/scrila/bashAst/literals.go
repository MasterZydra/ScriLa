package bashAst

// BoolLiteral

type IBoolLiteral interface {
	IStatement
	GetValue() bool
}

type BoolLiteral struct {
	stmt  *Statement
	value bool
}

func NewBoolLiteral(value bool) *BoolLiteral {
	return &BoolLiteral{
		stmt:  NewStatement(BoolLiteralNode),
		value: value,
	}
}

func (self *BoolLiteral) GetKind() NodeType {
	return self.stmt.GetKind()
}

func (self *BoolLiteral) GetValue() bool {
	return self.value
}

// IntLiteral

type IIntLiteral interface {
	IIntStmt
}

func NewIntLiteral(value int64) *IntStmt {
	return NewIntStmt(IntLiteralNode, value)
}

// StrLiteral

type IStrLiteral interface {
	IStrStmt
}

func NewStrLiteral(value string) *StrStmt {
	return NewStrStmt(StrLiteralNode, value)
}

// VarLiteral

type IVarLiteral interface {
	IStrStmt
	GetVarType() NodeType
}

type VarLiteral struct {
	stmt    *Statement
	value   string
	varType NodeType
}

func NewVarLiteral(name string, varType NodeType) *VarLiteral {
	return &VarLiteral{
		stmt:    NewStatement(VarLiteralNode),
		value:   name,
		varType: varType,
	}
}

func (self *VarLiteral) GetKind() NodeType {
	return self.stmt.GetKind()
}

func (self *VarLiteral) GetValue() string {
	return self.value
}

func (self *VarLiteral) GetVarType() NodeType {
	return self.varType
}
