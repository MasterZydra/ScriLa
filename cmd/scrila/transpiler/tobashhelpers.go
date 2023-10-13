package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"fmt"
)

// Returns the given string wrapped in double quotes
func strToBashStr(str string) string {
	return fmt.Sprintf("\"%s\"", str)
}

// Returns the symbol of the given expr of kind Identifier
func identNodeGetSymbol(expr ast.IExpr) string {
	return ast.ExprToIdent(expr).GetSymbol()
}

// Return the given Identifier as Bash variable.
func identNodeToBashVar(expr ast.IExpr) string {
	return fmt.Sprintf("${%s}", identNodeGetSymbol(expr))
}

// Return a bash comparision to represent a bool (true|false)
func boolIdentToBashComparison(ident ast.IIdentifier) string {
	return fmt.Sprintf("[[ %s == \"true\" ]]", strToBashStr(ident.GetSymbol()))
}

// Return a bash comparison for a given bool variable
func varIdentToBashComparison(ident ast.IIdentifier) string {
	return fmt.Sprintf("[[ \"%s\" == \"true\" ]]", identNodeToBashVar(ident))
}
