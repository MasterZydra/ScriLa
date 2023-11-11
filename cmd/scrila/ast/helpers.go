package ast

import (
	"ScriLa/cmd/scrila/lexer"

	"golang.org/x/exp/slices"
)

var bools = []string{"true", "false"}

func IdentIsBool(ident IIdentifier) bool {
	return slices.Contains(bools, ident.GetSymbol())
}

func BinExprIsComp(binOp IBinaryExpr) bool {
	return slices.Contains(lexer.ComparisonOps, binOp.GetOperator())
}

func BinExprIsBoolOp(binOp IBinaryExpr) bool {
	return slices.Contains(lexer.BooleanOps, binOp.GetOperator())
}

func BinExprReturnsBool(binOp IBinaryExpr) bool {
	return BinExprIsBoolOp(binOp) || BinExprIsComp(binOp)
}
