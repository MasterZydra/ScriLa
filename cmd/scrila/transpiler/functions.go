package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"
)

func nativePrintLn(args []ast.IExpr, env *Environment) (IRuntimeVal, error) {
	writeToFile("echo \"")
	if err := printArgs(args, env); err != nil {
		return NewNullVal(), err
	}
	writeLnToFile("\"")
	return NewNullVal(), nil
}

func nativePrint(args []ast.IExpr, env *Environment) (IRuntimeVal, error) {
	writeToFile("echo -n \"")
	if err := printArgs(args, env); err != nil {
		return NewNullVal(), err
	}
	writeLnToFile("\"")
	return NewNullVal(), nil
}

func printArgs(args []ast.IExpr, env *Environment) error {
	var isFirst bool = true
	for _, arg := range args {
		if !isFirst {
			writeToFile(" ")
		}
		if isFirst {
			isFirst = false
		}
		switch arg.GetKind() {
		case ast.IdentifierNode:
			var i interface{} = arg
			identifier, _ := i.(ast.IIdentifier)
			if slices.Contains(reservedIdentifiers, identifier.GetSymbol()) {
				writeToFile(identifier.GetSymbol())
			} else {
				writeToFile("${" + identifier.GetSymbol() + "}")
			}
		case ast.IntLiteralNode:
			var i interface{} = arg
			intLiteral, _ := i.(ast.IIntLiteral)
			writeToFile(strconv.Itoa(int(intLiteral.GetValue())))
		case ast.StrLiteralNode:
			var i interface{} = arg
			strLiteral, _ := i.(ast.IStrLiteral)
			writeToFile(strLiteral.GetValue())
		case ast.BinaryExprNode:
			value, err := transpile(arg, env)
			if err != nil {
				return err
			}
			writeToFile(value.GetTranspilat())
		default:
			return fmt.Errorf("nativePrint: Arg kind '%s' not supported", arg)
		}
	}
	return nil
}
