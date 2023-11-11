package bashAst

type IAppendBody interface {
	IStatement
	AppendBody(stmt IStatement)
}

// Base value types

type IIntValue interface {
	GetValue() int64
}

type IStrValue interface {
	GetValue() string
}

// IntStmt

type IIntStmt interface {
	IStatement
	IIntValue
}

type IntStmt struct {
	stmt  *Statement
	value int64
}

func NewIntStmt(kind NodeType, value int64) *IntStmt {
	return &IntStmt{
		stmt:  NewStatement(kind),
		value: value,
	}
}

func (self *IntStmt) GetKind() NodeType {
	return self.stmt.GetKind()
}

func (self *IntStmt) GetValue() int64 {
	return self.value
}

// IStrStmt

type IStrStmt interface {
	IStatement
	IStrValue
}

type StrStmt struct {
	stmt  *Statement
	value string
}

func NewStrStmt(kind NodeType, value string) *StrStmt {
	return &StrStmt{
		stmt:  NewStatement(kind),
		value: value,
	}
}

func (self *StrStmt) GetKind() NodeType {
	return self.stmt.GetKind()
}

func (self *StrStmt) GetValue() string {
	return self.value
}
