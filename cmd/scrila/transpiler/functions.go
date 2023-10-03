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
	argStr, err := printArgs(args, env)
	if err != nil {
		return NewNullVal(), err
	}
	writeLnToFile("echo " + strToBashStr(argStr))
	return NewNullVal(), nil
}

func nativePrint(args []ast.IExpr, env *Environment) (IRuntimeVal, error) {
	argStr, err := printArgs(args, env)
	if err != nil {
		return NewNullVal(), err
	}
	writeLnToFile("echo -n " + strToBashStr(argStr))
	return NewNullVal(), nil
}

func printArgs(args []ast.IExpr, env *Environment) (string, error) {
	argStr := ""
	var isFirst bool = true
	for _, arg := range args {
		if !isFirst {
			argStr += " "
		}
		if isFirst {
			isFirst = false
		}
		switch arg.GetKind() {
		case ast.CallExprNode:
			varName, err := getCallerResultVarName(ast.ExprToCallExpr(arg), env)
			if err != nil {
				return "", err
			}
			argStr += varName
		case ast.IdentifierNode:
			if symbol := identNodeGetSymbol(arg); slices.Contains(reservedIdentifiers, symbol) {
				argStr += symbol
			} else {
				argStr += identNodeToBashVar(arg)
			}
		case ast.IntLiteralNode:
			argStr += strconv.Itoa(int(ast.ExprToIntLit(arg).GetValue()))
		case ast.StrLiteralNode:
			argStr += ast.ExprToStrLit(arg).GetValue()
		case ast.BinaryExprNode:
			value, err := transpile(arg, env)
			if err != nil {
				return "", err
			}
			argStr += value.GetTranspilat()
		case ast.MemberExprNode:
			memberVal, err := evalMemberExpr(ast.ExprToMemberExpr(arg), env)
			if err != nil {
				return "", err
			}
			argStr += memberVal.GetTranspilat()
		default:
			return "", fmt.Errorf("nativePrint: Arg kind '%s' not supported", arg.GetKind())
		}
	}
	return argStr, nil
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
		varType, err := env.lookupVarType(identNodeGetSymbol(args[0]))
		if err != nil {
			return NewNullVal(), err
		}
		if varType != lexer.StrType {
			return NewNullVal(), fmt.Errorf("input: parameter prompt has to be a string or a variable of type string. Got '%s'", varType)
		}

		transpilat += strToBashStr(identNodeToBashVar(args[0]) + " ")
	case ast.StrLiteralNode:
		transpilat += strToBashStr(value.ToString() + " ")
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
		symbol := identNodeGetSymbol(args[0])
		varType, err := env.lookupVarType(symbol)
		if err != nil {
			return NewNullVal(), err
		}
		if varType != lexer.IntType {
			return NewNullVal(), fmt.Errorf("sleep: parameter has to be a int or a variable of type int. Got '%s'", varType)
		}

		transpilat += identNodeToBashVar(args[0]) + "\n"
	case ast.IntLiteralNode:
		transpilat += value.ToString() + "\n"
	default:
		return NewNullVal(), fmt.Errorf("nativeSleep: Arg kind '%s' not supported", args[0].GetKind())
	}
	result := NewNullVal()
	result.SetTranspilat(transpilat)
	return result, nil
}
