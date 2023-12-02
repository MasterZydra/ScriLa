package bashTranspiler

import (
	"ScriLa/cmd/scrila/bashAst"
	"ScriLa/cmd/scrila/scrilaAst"
	"fmt"

	"golang.org/x/exp/slices"
)

func (self *Transpiler) declareNativeFunctions(env *Environment) {
	env.declareFunc("exec", NewNativeFunc(self.nativeExec, scrilaAst.VoidNode))
	env.declareFunc("exit", NewNativeFunc(self.nativeExit, scrilaAst.VoidNode))
	env.declareFunc("input", NewNativeFunc(self.nativeInput, scrilaAst.StrLiteralNode))
	env.declareFunc("print", NewNativeFunc(self.nativePrintLn, scrilaAst.VoidNode))
	env.declareFunc("printLn", NewNativeFunc(self.nativePrintLn, scrilaAst.VoidNode))
	env.declareFunc("sleep", NewNativeFunc(self.nativeSleep, scrilaAst.VoidNode))
	env.declareFunc("strIsBool", NewNativeFunc(self.nativeStrIsBool, scrilaAst.BoolLiteralNode))
	env.declareFunc("strIsInt", NewNativeFunc(self.nativeStrIsInt, scrilaAst.BoolLiteralNode))
	env.declareFunc("strToBool", NewNativeFunc(self.nativeStrToBool, scrilaAst.BoolLiteralNode))
	env.declareFunc("strToInt", NewNativeFunc(self.nativeStrToInt, scrilaAst.IntLiteralNode))
}

func (self *Transpiler) nativePrintLn(args []scrilaAst.IExpr, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")
	return NewNullVal(), nil
}

func (self *Transpiler) nativeInput(args []scrilaAst.IExpr, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	// Validate args
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: input(str prompt)")
	}

	doMatch, givenType, err := self.exprIsType(args[0], scrilaAst.StrLiteralNode, env)
	if err != nil {
		return NewNullVal(), err
	}
	if !doMatch {
		return NewNullVal(), fmt.Errorf("input() - Parameter prompt must be a string or a variable of type string. Got '%s'", givenType)
	}

	// Add bash code for input to "usedNativeFunctions"
	if !slices.Contains(self.usedNativeFunctions, "input") {
		self.usedNativeFunctions = append(self.usedNativeFunctions, "input")
		funcDecl := bashAst.NewFuncDeclaration("input", bashAst.StrLiteralNode)
		funcDecl.AppendParams(bashAst.NewFuncParameter("prompt", bashAst.StrLiteralNode))
		funcDecl.AppendBody(bashAst.NewBashStmt("read -p \"${prompt} \" tmpStrs[${tmpInts[0]}]"))
		self.bashProgram.AppendNativeBody(funcDecl)
	}

	return NewStrVal("str"), nil
}

func (self *Transpiler) nativeSleep(args []scrilaAst.IExpr, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	// Validate args
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: sleep(int seconds)")
	}

	doMatch, givenType, err := self.exprIsType(args[0], scrilaAst.IntLiteralNode, env)
	if err != nil {
		return NewNullVal(), err
	}
	if !doMatch {
		return NewNullVal(), fmt.Errorf("sleep() - Parameter seconds must be an int or a variable of type int. Got '%s'", givenType)
	}

	return NewNullVal(), nil
}

func (self *Transpiler) nativeStrIsBool(args []scrilaAst.IExpr, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	// Validate args
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: strIsBool(str value)")
	}

	doMatch, givenType, err := self.exprIsType(args[0], scrilaAst.StrLiteralNode, env)
	if err != nil {
		return NewNullVal(), err
	}
	if !doMatch {
		return NewNullVal(), fmt.Errorf("strIsBool() - Parameter value must be a string or a variable of type string. Got '%s'", givenType)
	}

	// Add bash code for strIsBool to "usedNativeFunctions"
	if !slices.Contains(self.usedNativeFunctions, "strIsBool") {
		self.usedNativeFunctions = append(self.usedNativeFunctions, "strIsBool")
		funcDecl := bashAst.NewFuncDeclaration("strIsBool", bashAst.BoolLiteralNode)
		funcDecl.AppendParams(bashAst.NewFuncParameter("value", bashAst.StrLiteralNode))
		cond := bashAst.NewBinaryOpExpr(
			bashAst.BoolLiteralNode,
			bashAst.NewBinaryCompExpr(bashAst.BoolLiteralNode, bashAst.NewVarLiteral("value", bashAst.StrLiteralNode), bashAst.NewStrLiteral("true"), "=="),
			bashAst.NewBinaryCompExpr(bashAst.BoolLiteralNode, bashAst.NewVarLiteral("value", bashAst.StrLiteralNode), bashAst.NewStrLiteral("false"), "=="),
			"||")
		ifStmt := bashAst.NewIfStmt(cond)
		ifStmt.AppendBody(bashAst.NewBashStmt("tmpBools[${tmpInts[0]}]=\"true\""))
		elseStmt := bashAst.NewIfStmt(nil)
		elseStmt.AppendBody(bashAst.NewBashStmt("tmpBools[${tmpInts[0]}]=\"false\""))
		ifStmt.SetElse(elseStmt)
		funcDecl.AppendBody(ifStmt)
		self.bashProgram.AppendNativeBody(funcDecl)
	}
	return NewBoolVal(true), nil
}

func (self *Transpiler) nativeStrIsInt(args []scrilaAst.IExpr, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	// Validate args
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: strIsInt(str value)")
	}

	doMatch, givenType, err := self.exprIsType(args[0], scrilaAst.StrLiteralNode, env)
	if err != nil {
		return NewNullVal(), err
	}
	if !doMatch {
		return NewNullVal(), fmt.Errorf("strIsInt() - Parameter value must be a string or a variable of type string. Got '%s'", givenType)
	}

	// Add bash code for strIsInt to "usedNativeFunctions"
	if !slices.Contains(self.usedNativeFunctions, "strIsInt") {
		self.usedNativeFunctions = append(self.usedNativeFunctions, "strIsInt")
		// https://stackoverflow.com/questions/806906/how-do-i-test-if-a-variable-is-a-number-in-bash/3951175#3951175
		funcDecl := bashAst.NewFuncDeclaration("strIsInt", bashAst.BoolLiteralNode)
		funcDecl.AppendParams(bashAst.NewFuncParameter("value", bashAst.StrLiteralNode))
		funcDecl.AppendBody(bashAst.NewBashStmt("case ${value} in"))
		funcDecl.AppendBody(bashAst.NewBashStmt("\t''|*[!0-9]*) tmpBools[${tmpInts[0]}]=\"false\" ;;"))
		funcDecl.AppendBody(bashAst.NewBashStmt("\t*) tmpBools[${tmpInts[0]}]=\"true\" ;;"))
		funcDecl.AppendBody(bashAst.NewBashStmt("esac"))
		self.bashProgram.AppendNativeBody(funcDecl)
	}
	return NewBoolVal(true), nil
}

func (self *Transpiler) nativeStrToBool(args []scrilaAst.IExpr, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	// Validate args
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: strToBool(str value)")
	}

	doMatch, givenType, err := self.exprIsType(args[0], scrilaAst.StrLiteralNode, env)
	if err != nil {
		return NewNullVal(), err
	}
	if !doMatch {
		return NewNullVal(), fmt.Errorf("strToBool() - Parameter value must be a string or a variable of type string. Got '%s'", givenType)
	}

	// Add bash code for strToBool to "usedNativeFunctions"
	if !slices.Contains(self.usedNativeFunctions, "strToBool") {
		self.usedNativeFunctions = append(self.usedNativeFunctions, "strToBool")
		funcDecl := bashAst.NewFuncDeclaration("strToBool", bashAst.BoolLiteralNode)
		funcDecl.AppendParams(bashAst.NewFuncParameter("value", bashAst.StrLiteralNode))
		cond := bashAst.NewBinaryCompExpr(bashAst.BoolLiteralNode, bashAst.NewVarLiteral("value", bashAst.StrLiteralNode), bashAst.NewStrLiteral("true"), "==")
		ifStmt := bashAst.NewIfStmt(cond)
		ifStmt.AppendBody(bashAst.NewBashStmt("tmpBools[${tmpInts[0]}]=\"true\""))
		elseStmt := bashAst.NewIfStmt(nil)
		elseStmt.AppendBody(bashAst.NewBashStmt("tmpBools[${tmpInts[0]}]=\"false\""))
		ifStmt.SetElse(elseStmt)
		funcDecl.AppendBody(ifStmt)
		self.bashProgram.AppendNativeBody(funcDecl)
	}
	return NewBoolVal(true), nil
}

func (self *Transpiler) nativeStrToInt(args []scrilaAst.IExpr, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	// Validate args
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: strToInt(str value)")
	}

	doMatch, givenType, err := self.exprIsType(args[0], scrilaAst.StrLiteralNode, env)
	if err != nil {
		return NewNullVal(), err
	}
	if !doMatch {
		return NewNullVal(), fmt.Errorf("strToInt() - Parameter value must be a string or a variable of type string. Got '%s'", givenType)
	}

	// TODO After error handling in ScriLa is thought-out: Add error handling for the case that the value is not an int.

	// Add bash code for strToInt to "usedNativeFunctions"
	if !slices.Contains(self.usedNativeFunctions, "strToInt") {
		self.usedNativeFunctions = append(self.usedNativeFunctions, "strToInt")
		funcDecl := bashAst.NewFuncDeclaration("strToInt", bashAst.IntLiteralNode)
		funcDecl.AppendParams(bashAst.NewFuncParameter("value", bashAst.StrLiteralNode))
		funcDecl.AppendBody(bashAst.NewBashStmt("tmpInts[${tmpInts[0]}]=${value}"))
		self.bashProgram.AppendNativeBody(funcDecl)
	}

	return NewIntVal(1), nil
}

func (self *Transpiler) nativeExec(args []scrilaAst.IExpr, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	// Validate args
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: exec(str command)")
	}

	doMatch, givenType, err := self.exprIsType(args[0], scrilaAst.StrLiteralNode, env)
	if err != nil {
		return NewNullVal(), err
	}
	if !doMatch {
		return NewNullVal(), fmt.Errorf("exec() - Parameter value must be a string or a variable of type string. Got '%s'", givenType)
	}

	// Add bash code for exec to "usedNativeFunctions"
	if !slices.Contains(self.usedNativeFunctions, "exec") {
		self.usedNativeFunctions = append(self.usedNativeFunctions, "exec")
		funcDecl := bashAst.NewFuncDeclaration("exec", bashAst.VoidNode)
		funcDecl.AppendParams(bashAst.NewFuncParameter("command", bashAst.StrLiteralNode))
		funcDecl.AppendBody(bashAst.NewBashStmt("${command}"))
		self.bashProgram.AppendNativeBody(funcDecl)
	}

	return NewNullVal(), nil
}

func (self *Transpiler) nativeExit(args []scrilaAst.IExpr, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	// Validate args
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: exit(int code)")
	}

	doMatch, givenType, err := self.exprIsType(args[0], scrilaAst.IntLiteralNode, env)
	if err != nil {
		return NewNullVal(), err
	}
	if !doMatch {
		return NewNullVal(), fmt.Errorf("exit() - Parameter value must be a int or a variable of type int. Got '%s'", givenType)
	}

	return NewNullVal(), nil
}
