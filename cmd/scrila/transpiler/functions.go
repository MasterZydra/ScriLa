package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"time"
)

func nativeTime(args []ast.IExpr, env *Environment) IRuntimeVal {
	return NewIntVal(time.Now().UnixMilli())
}

func nativePrint(args []ast.IExpr, env *Environment) IRuntimeVal {
	writeToFile("echo \"")
	for _, arg := range args {
		if arg.GetKind() == ast.IdentifierNode {
			var i interface{} = arg
			identifier, _ := i.(ast.IIdentifier)
			writeToFile("$" + identifier.GetSymbol())
		}
	}
	writeLnToFile("\"")
	return NewNullVal()
}
