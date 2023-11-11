package scrilaAst

import (
	"golang.org/x/exp/slices"
)

var bools = []string{"true", "false"}

func IdentIsBool(ident IIdentifier) bool {
	return slices.Contains(bools, ident.GetSymbol())
}

var ComparisonOps = []string{"<", ">", "<=", ">=", "!=", "=="}

func BinExprIsComp(binOp IBinaryExpr) bool {
	return slices.Contains(ComparisonOps, binOp.GetOperator())
}

var BooleanOps = []string{"||", "&&"}

func BinExprIsBoolOp(binOp IBinaryExpr) bool {
	return slices.Contains(BooleanOps, binOp.GetOperator())
}

func BinExprReturnsBool(binOp IBinaryExpr) bool {
	return BinExprIsBoolOp(binOp) || BinExprIsComp(binOp)
}
