package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"ScriLa/cmd/scrila/lexer"
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"
)

func declareNativeFunctions(env *Environment) {
	var nativeFunctions = map[string]FunctionCall{
		"input":   nativeInput,
		"print":   nativePrint,
		"printLn": nativePrintLn,
		"sleep":   nativeSleep,
	}

	for name, function := range nativeFunctions {
		env.declareFunc(name, NewNativeFunc(function))
	}
}

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
		case ast.CallExprNode:
			var i interface{} = arg
			callExpr, _ := i.(ast.ICallExpr)
			varName, err := getCallerResultVarName(callExpr, env)
			if err != nil {
				return err
			}
			writeToFile(varName)
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
		case ast.MemberExprNode:
			var i interface{} = arg
			memberExpr, _ := i.(ast.IMemberExpr)
			memberVal, err := evalMemberExpr(memberExpr, env)
			if err != nil {
				return err
			}
			writeToFile(memberVal.GetTranspilat())
		default:
			return fmt.Errorf("nativePrint: Arg kind '%s' not supported", arg.GetKind())
		}
	}
	return nil
}

func nativeInput(args []ast.IExpr, env *Environment) (IRuntimeVal, error) {
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: input(str prompt)")
	}
	value, err := transpile(args[0], env)
	if err != nil {
		return NewNullVal(), err
	}

	transpilat := "read -p "

	switch args[0].GetKind() {
	case ast.IdentifierNode:
		var i interface{} = args[0]
		identifier, _ := i.(ast.IIdentifier)
		varType, err := env.lookupVarType(identifier.GetSymbol())
		if err != nil {
			return NewNullVal(), err
		}
		if varType != lexer.StrType {
			return NewNullVal(), fmt.Errorf("input: parameter prompt has to be a string or a variable of type string. Got '%s'", varType)
		}

		transpilat += "\"$" + identifier.GetSymbol() + " \""
	case ast.StrLiteralNode:
		transpilat += "\"" + value.ToString() + " \""
	default:
		return NewNullVal(), fmt.Errorf("nativeInput: Arg kind '%s' not supported", args[0].GetKind())
	}

	transpilat += " tmpStr\n"

	result := NewNullVal()
	result.SetTranspilat(transpilat)
	return result, nil
}

func nativeSleep(args []ast.IExpr, env *Environment) (IRuntimeVal, error) {
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: sleep(int seconds)")
	}
	value, err := transpile(args[0], env)
	if err != nil {
		return NewNullVal(), err
	}

	transpilat := "sleep "
	switch args[0].GetKind() {
	case ast.IdentifierNode:
		var i interface{} = args[0]
		identifier, _ := i.(ast.IIdentifier)
		varType, err := env.lookupVarType(identifier.GetSymbol())
		if err != nil {
			return NewNullVal(), err
		}
		if varType != lexer.IntType {
			return NewNullVal(), fmt.Errorf("sleep: parameter has to be a int or a variable of type int. Got '%s'", varType)
		}

		transpilat += "$" + identifier.GetSymbol() + "\n"
	case ast.IntLiteralNode:
		transpilat += value.ToString() + "\n"
	default:
		return NewNullVal(), fmt.Errorf("nativeSleep: Arg kind '%s' not supported", args[0].GetKind())
	}
	result := NewNullVal()
	result.SetTranspilat(transpilat)
	return result, nil
}
