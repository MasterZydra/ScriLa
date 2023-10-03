package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"fmt"
)

func strToBashStr(str string) string {
	return fmt.Sprintf("\"%s\"", str)
}

func identNodeGetSymbol(expr ast.IExpr) string {
	return ast.ExprToIdent(expr).GetSymbol()
}

// Return the given Identifier as Bash variable.
func identNodeToBashVar(expr ast.IExpr) string {
	return fmt.Sprintf("${%s}", identNodeGetSymbol(expr))
}
