package scrilaAst

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

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

var indentDepth int = 0

func indent() string {
	return strings.Repeat("  ", indentDepth+1)
}

func SprintAST(program IProgram) string {
	astString := ""
	for _, stmt := range program.GetBody() {
		astString += fmt.Sprintf("%s%s\n", indent(), stmt)
	}
	return astString
}
