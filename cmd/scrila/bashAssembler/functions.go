package bashAssembler

import (
	"ScriLa/cmd/scrila/bashAst"
	"fmt"
)

type nativeScrilaFunc func(args []bashAst.IStatement) error

func (self *Assembler) registerNativeScrilaFuncs() {
	self.nativeScrilaFuncs = map[string]nativeScrilaFunc{
		"exit":    self.nativeFnExit,
		"print":   self.nativeFnPrint,
		"printLn": self.nativeFnPrintLn,
		"sleep":   self.nativeFnSleep,
	}
}

func (self *Assembler) nativeFnExit(args []bashAst.IStatement) error {
	bash, err := stmtToBashStr(args[0])
	if err != nil {
		return err
	}
	self.writeLnWithTabsToFile(fmt.Sprintf("exit %s", bash))
	return nil
}

func (self *Assembler) nativeFnPrint(args []bashAst.IStatement) error {
	self.writeWithTabsToFile("echo -n ")
	argStr, err := printArgsToBashStr(args)
	if err != nil {
		return err
	}
	self.writeLnToFile(strToBashStr(argStr))
	return nil
}

func (self *Assembler) nativeFnPrintLn(args []bashAst.IStatement) error {
	self.writeWithTabsToFile("echo ")
	argStr, err := printArgsToBashStr(args)
	if err != nil {
		return err
	}
	self.writeLnToFile(strToBashStr(argStr))
	return nil
}

func (self *Assembler) nativeFnSleep(args []bashAst.IStatement) error {
	bash, err := stmtToBashStr(args[0])
	if err != nil {
		return err
	}
	self.writeLnWithTabsToFile(fmt.Sprintf("sleep %s", bash))
	return nil
}

func printArgsToBashStr(args []bashAst.IStatement) (string, error) {
	argStr := ""
	for i, arg := range args {
		bash, err := stmtToBashStr(arg)
		if err != nil {
			return "", err
		}
		if i > 0 {
			argStr += " "
		}
		argStr += bash
	}
	return argStr, nil
}

func (self *Assembler) getFuncSignature(funcDecl bashAst.IFuncDeclaration) (string, error) {
	params := ""
	isFirstParam := true
	for _, param := range funcDecl.GetParams() {
		varType, err := nodeTypeToVarTypeKeyword(param.GetType())
		if err != nil {
			return "", err
		}

		if !isFirstParam {
			params += ", "
		}
		params += varType + " " + param.GetName()

		isFirstParam = false
	}

	returnType, err := nodeTypeToVarTypeKeyword(funcDecl.GetReturnType())
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s(%s) %s", funcDecl.GetName(), params, returnType), nil
}
