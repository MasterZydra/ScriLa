package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"ScriLa/cmd/scrila/lexer"
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"
)

func (self *Transpiler) declareNativeFunctions(env *Environment) {
	var nativeFunctions = map[string]FunctionCall{
		"input":   self.nativeInput,
		"print":   self.nativePrint,
		"printLn": self.nativePrintLn,
		"sleep":   self.nativeSleep,
		"isInt":   self.nativeIsInt,
	}

	for name, function := range nativeFunctions {
		env.declareFunc(name, NewNativeFunc(function))
	}
}

func (self *Transpiler) nativePrintLn(args []ast.IExpr, env *Environment) (IRuntimeVal, error) {
	argStr, err := self.printArgs(args, env)
	if err != nil {
		return NewNullVal(), err
	}
	self.writeLnTranspilat("echo " + strToBashStr(argStr))
	return NewNullVal(), nil
}

func (self *Transpiler) nativePrint(args []ast.IExpr, env *Environment) (IRuntimeVal, error) {
	argStr, err := self.printArgs(args, env)
	if err != nil {
		return NewNullVal(), err
	}
	self.writeLnTranspilat("echo -n " + strToBashStr(argStr))
	return NewNullVal(), nil
}

func (self *Transpiler) printArgs(args []ast.IExpr, env *Environment) (string, error) {
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
			varName, err := self.getCallerResultVarName(ast.ExprToCallExpr(arg), env)
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
			value, err := self.transpile(arg, env)
			if err != nil {
				return "", err
			}
			argStr += value.GetTranspilat()
		case ast.MemberExprNode:
			memberVal, err := self.evalMemberExpr(ast.ExprToMemberExpr(arg), env)
			if err != nil {
				return "", err
			}
			argStr += memberVal.GetTranspilat()
		default:
			return "", fmt.Errorf("print() - Unexpected %s expression at %s:%d:%d", arg.GetKind(), self.filename, arg.GetLn(), arg.GetCol())
		}
	}
	return argStr, nil
}

func (self *Transpiler) nativeInput(args []ast.IExpr, env *Environment) (IRuntimeVal, error) {
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: input(str prompt)")
	}
	value, err := self.transpile(args[0], env)
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
			return NewNullVal(), fmt.Errorf("input() - Parameter prompt must be a string or a variable of type string. Got '%s'", varType)
		}

		transpilat += strToBashStr(identNodeToBashVar(args[0]) + " ")
	case ast.StrLiteralNode:
		transpilat += strToBashStr(value.ToString() + " ")
	default:
		return NewNullVal(), fmt.Errorf("input() - Parameter prompt must be a string or a variable of type string. Got '%s'", args[0].GetKind())
	}

	transpilat += " tmpStr\n"

	result := NewNullVal()
	result.SetTranspilat(transpilat)
	return result, nil
}

func (self *Transpiler) nativeSleep(args []ast.IExpr, env *Environment) (IRuntimeVal, error) {
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: sleep(int seconds)")
	}
	value, err := self.transpile(args[0], env)
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
			return NewNullVal(), fmt.Errorf("sleep() - Parameter seconds must be an int or a variable of type int. Got '%s'", varType)
		}

		transpilat += identNodeToBashVar(args[0]) + "\n"
	case ast.IntLiteralNode:
		transpilat += value.ToString() + "\n"
	default:
		return NewNullVal(), fmt.Errorf("sleep() - Parameter seconds must be an int or a variable of type int. Got '%s'", args[0].GetKind())
	}
	result := NewNullVal()
	result.SetTranspilat(transpilat)
	return result, nil
}

func (self *Transpiler) nativeIsInt(args []ast.IExpr, env *Environment) (IRuntimeVal, error) {
	// Validate args
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: isInt(mixed value)")
	}
	value, err := self.transpile(args[0], env)
	if err != nil {
		return NewNullVal(), err
	}

	// Add bash code for isInt to "usedNativeFunctions"
	if !slices.Contains(self.usedNativeFunctions, "isInt") {
		self.usedNativeFunctions = append(self.usedNativeFunctions, "isInt")
		self.nativeFuncTranspilat += "isInt () {\n"
		self.nativeFuncTranspilat += "\tcase $1 in\n"
		self.nativeFuncTranspilat += "\t\t''|*[!0-9]*) tmpBool=\"false\" ;;\n"
		self.nativeFuncTranspilat += "\t\t*) tmpBool=\"true\" ;;\n"
		self.nativeFuncTranspilat += "\tesac\n"
		self.nativeFuncTranspilat += "}\n\n"
	}

	transpilat := "isInt "
	switch args[0].GetKind() {
	case ast.IdentifierNode:
		if ast.IdentIsBool(ast.ExprToIdent(args[0])) {
			transpilat += strToBashStr(identNodeGetSymbol(args[0]))
		} else {
			transpilat += strToBashStr(identNodeToBashVar(args[0]))
		}
	case ast.IntLiteralNode:
		transpilat += value.ToString()
	case ast.StrLiteralNode:
		transpilat += strToBashStr(value.ToString())
	default:
		return NewNullVal(), fmt.Errorf("isInt() - Support for type %s not implemented", args[0].GetKind())
	}
	transpilat += "\n"

	result := NewBoolVal(true)
	result.SetTranspilat(transpilat)
	return result, nil
}
