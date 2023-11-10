package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"ScriLa/cmd/scrila/lexer"
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"
)

func (self *Transpiler) declareNativeFunctions(env *Environment) {
	env.declareFunc("input", NewNativeFunc(self.nativeInput, lexer.StrType))
	env.declareFunc("print", NewNativeFunc(self.nativePrint, lexer.VoidType))
	env.declareFunc("printLn", NewNativeFunc(self.nativePrintLn, lexer.VoidType))
	env.declareFunc("sleep", NewNativeFunc(self.nativeSleep, lexer.VoidType))
	env.declareFunc("strIsInt", NewNativeFunc(self.nativeStrIsInt, lexer.BoolType))
	env.declareFunc("strToInt", NewNativeFunc(self.nativeStrToInt, lexer.IntType))
	env.declareFunc("exec", NewNativeFunc(self.nativeExec, lexer.VoidType))
}

func (self *Transpiler) nativePrintLn(args []ast.IExpr, env *Environment) (ast.IRuntimeVal, error) {
	self.printFuncName("")

	argStr, err := self.printArgs(args, env)
	if err != nil {
		return NewNullVal(), err
	}
	self.writeLnTranspilat("echo " + strToBashStr(argStr))
	return NewNullVal(), nil
}

func (self *Transpiler) nativePrint(args []ast.IExpr, env *Environment) (ast.IRuntimeVal, error) {
	self.printFuncName("")

	argStr, err := self.printArgs(args, env)
	if err != nil {
		return NewNullVal(), err
	}
	self.writeLnTranspilat("echo -n " + strToBashStr(argStr))
	return NewNullVal(), nil
}

func (self *Transpiler) printArgs(args []ast.IExpr, env *Environment) (string, error) {
	self.printFuncName("")

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

func (self *Transpiler) nativeInput(args []ast.IExpr, env *Environment) (ast.IRuntimeVal, error) {
	self.printFuncName("")

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

	result := NewStrVal("str")
	result.SetTranspilat(transpilat)
	return result, nil
}

func (self *Transpiler) nativeSleep(args []ast.IExpr, env *Environment) (ast.IRuntimeVal, error) {
	self.printFuncName("")

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

func (self *Transpiler) nativeStrIsInt(args []ast.IExpr, env *Environment) (ast.IRuntimeVal, error) {
	self.printFuncName("")

	// Validate args
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: strIsInt(str value)")
	}
	value, err := self.transpile(args[0], env)
	if err != nil {
		return NewNullVal(), err
	}

	// Add bash code for strIsInt to "usedNativeFunctions"
	if !slices.Contains(self.usedNativeFunctions, "strIsInt") {
		self.usedNativeFunctions = append(self.usedNativeFunctions, "strIsInt")
		// https://stackoverflow.com/questions/806906/how-do-i-test-if-a-variable-is-a-number-in-bash/3951175#3951175
		self.nativeFuncTranspilat += "strIsInt () {\n"
		self.nativeFuncTranspilat += "\tcase $1 in\n"
		self.nativeFuncTranspilat += "\t\t''|*[!0-9]*) tmpBool=\"false\" ;;\n"
		self.nativeFuncTranspilat += "\t\t*) tmpBool=\"true\" ;;\n"
		self.nativeFuncTranspilat += "\tesac\n"
		self.nativeFuncTranspilat += "}\n\n"
	}

	transpilat := "strIsInt "
	switch args[0].GetKind() {
	case ast.IdentifierNode:
		varType, err := env.lookupVarType(identNodeGetSymbol(args[0]))
		if err != nil {
			return NewNullVal(), err
		}
		if varType != lexer.StrType {
			return NewNullVal(), fmt.Errorf("strIsInt() - Parameter value must be a string or a variable of type string. Got '%s'", varType)
		}

		transpilat += strToBashStr(identNodeToBashVar(args[0]))
	case ast.StrLiteralNode:
		transpilat += strToBashStr(value.ToString())
	default:
		return NewNullVal(), fmt.Errorf("strIsInt() - Parameter value must be a string or a variable of type string. Got '%s'", args[0].GetKind())
	}
	transpilat += "\n"

	result := NewBoolVal(true)
	result.SetTranspilat(transpilat)
	return result, nil
}

func (self *Transpiler) nativeStrToInt(args []ast.IExpr, env *Environment) (ast.IRuntimeVal, error) {
	self.printFuncName("")

	// Validate args
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: strToInt(str value)")
	}
	value, err := self.transpile(args[0], env)
	if err != nil {
		return NewNullVal(), err
	}
	// TODO After error handling in ScriLa is thought-out: Add error handling for the case that the value is not an int.
	transpilat := "tmpInt="
	switch args[0].GetKind() {
	case ast.IdentifierNode:
		varType, err := env.lookupVarType(identNodeGetSymbol(args[0]))
		if err != nil {
			return NewNullVal(), err
		}
		if varType != lexer.StrType {
			return NewNullVal(), fmt.Errorf("strToInt() - Parameter value must be a string or a variable of type string. Got '%s'", varType)
		}

		transpilat += strToBashStr(identNodeToBashVar(args[0]))
	case ast.StrLiteralNode:
		transpilat += strToBashStr(value.ToString())
	default:
		return NewNullVal(), fmt.Errorf("strToInt() - Parameter value must be a string or a variable of type string. Got '%s'", args[0].GetKind())
	}
	transpilat += "\n"

	result := NewIntVal(1)
	result.SetTranspilat(transpilat)
	return result, nil
}

func (self *Transpiler) nativeExec(args []ast.IExpr, env *Environment) (ast.IRuntimeVal, error) {
	self.printFuncName("")

	// Validate args
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: exec(str command)")
	}
	value, err := self.transpile(args[0], env)
	if err != nil {
		return NewNullVal(), err
	}
	transpilat := ""
	switch args[0].GetKind() {
	case ast.IdentifierNode:
		varType, err := env.lookupVarType(identNodeGetSymbol(args[0]))
		if err != nil {
			return NewNullVal(), err
		}
		if varType != lexer.StrType {
			return NewNullVal(), fmt.Errorf("exec() - Parameter value must be a string or a variable of type string. Got '%s'", varType)
		}

		transpilat += identNodeToBashVar(args[0])
	case ast.StrLiteralNode:
		transpilat += value.ToString()
	default:
		return NewNullVal(), fmt.Errorf("exec() - Parameter value must be a string or a variable of type string. Got '%s'", args[0].GetKind())
	}
	transpilat += "\n"

	result := NewNullVal()
	result.SetTranspilat(transpilat)
	return result, nil
}
