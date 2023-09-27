package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"time"
)

func nativeTime(args []ast.IExpr, env *Environment) (IRuntimeVal, error) {
	return NewIntVal(time.Now().UnixMilli()), nil
}

func nativePrint(args []ast.IExpr, env *Environment) (IRuntimeVal, error) {
	writeToFile("echo \"")
	for _, arg := range args {
		if arg.GetKind() == ast.IdentifierNode {
			var i interface{} = arg
			identifier, _ := i.(ast.IIdentifier)
			writeToFile("$" + identifier.GetSymbol())
		}
	}
	writeLnToFile("\"")
	return NewNullVal(), nil
}
