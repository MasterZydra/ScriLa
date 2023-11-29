package bashAst

import "fmt"

// Array

type IArray interface {
	IStatement
	AddValue(value IStatement)
	GetValues() []IStatement
	GetDataType() NodeType
	SetDataType(dataType NodeType)
}

type Array struct {
	stmt     *Statement
	values   []IStatement
	dataType NodeType
}

func (self *Array) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - dataType: %s", self.GetKind(), self.GetDataType())
	for i, val := range self.values {
		str += fmt.Sprintf("\n%s%d: %s", indent(), i, val)
	}
	indentDepth--
	return str + "}"
}

func NewArray() *Array {
	return &Array{stmt: NewStatement(ArrayLiteralNode)}
}

func (self *Array) AddValue(value IStatement) {
	self.values = append(self.values, value)
}

func (self *Array) GetKind() NodeType {
	return self.stmt.GetKind()
}

func (self *Array) GetValues() []IStatement {
	return self.values
}

func (self *Array) GetDataType() NodeType {
	return self.dataType
}

func (self *Array) SetDataType(dataType NodeType) {
	self.dataType = dataType
}

// BoolLiteral

type IBoolLiteral interface {
	IStatement
	GetValue() bool
}

type BoolLiteral struct {
	stmt  *Statement
	value bool
}

func (self *BoolLiteral) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - value: %t}", self.GetKind(), self.GetValue())
	indentDepth--
	return str
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
	GetDataType() NodeType
}

type VarLiteral struct {
	stmt    *Statement
	value   string
	varType NodeType
}

func (self *VarLiteral) String() string {
	indentDepth++
	str := fmt.Sprintf("{%s - varName: '%s', varType: '%s'}", self.GetKind(), self.GetValue(), self.GetDataType())
	indentDepth--
	return str
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

func (self *VarLiteral) GetDataType() NodeType {
	return self.varType
}
