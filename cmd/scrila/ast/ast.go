package ast

type NodeType string

var nodeTypes = []string{
	"Program",
	"NumericLiteral",
	"Identifier",
	"BinaryExpr",
	"CallExpr",
	"UnaryExpr",
	"FunctionDeclaration",
}

// Statement:
//  Do not return a value
//  E.g. "int x = 42;"

// Expression:
//   An assignment is not a statement
//   E.g. "x = 42;"

type Statement struct {
	kind NodeType
}

type Program struct {
	kind NodeType
	body []Statement
}

func NewProgram() *Program {
	return &Program{
		kind: "Program",
		body: make([]Statement, 0),
	}
}

type Expr struct {
	kind NodeType
}

type BinaryExpr struct {
	kind     NodeType
	left     Expr
	right    Expr
	operator string
}

func NewBinaryExpr(left Expr, right Expr, operator string) *BinaryExpr {
	return &BinaryExpr{
		kind:     "BinaryExpr",
		left:     left,
		right:    right,
		operator: operator,
	}
}

type Identifier struct {
	kind   NodeType
	symbol string
}

func NewIdentifier(symbol string) *Identifier {
	return &Identifier{
		kind:   "Identifier",
		symbol: symbol,
	}
}

type NumericLiteral struct {
	kind  NodeType
	value int64
}

func NewNumericLiteral(value int64) *NumericLiteral {
	return &NumericLiteral{
		kind:  "NumericLiteral",
		value: value,
	}
}
