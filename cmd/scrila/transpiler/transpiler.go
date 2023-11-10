package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"fmt"
	"os"
	"runtime"
	"strings"

	"golang.org/x/exp/slices"
)

type Context string

const (
	NoContext        Context = "NoContext"
	FunctionContext  Context = "FunctionContext"
	WhileLoopContext Context = "WhileLoopContext"
	IfStmtContext    Context = "IfStmtContext"
)

type Transpiler struct {
	usedNativeFunctions  []string
	nativeFuncTranspilat string
	userScriptTranspilat string

	filename       string
	outputFilename string
	outputFile     *os.File

	testMode      bool
	testPrintMode bool
	contexts      []Context
	currentFunc   IFunctionVal

	showCallStack bool
}

func NewTranspiler(showCallStack bool) *Transpiler {
	return &Transpiler{
		usedNativeFunctions: []string{},
		contexts:            []Context{NoContext},
		showCallStack:       showCallStack,
	}
}

func (self *Transpiler) writeLnTranspilat(content string) {
	self.printFuncName("")

	self.writeTranspilat(content + "\n")
}

func (self *Transpiler) writeTranspilat(content string) {
	self.printFuncName("")

	self.userScriptTranspilat += content
}

func (self *Transpiler) writeLnToFile(content string) {
	self.printFuncName("")

	self.writeToFile(content + "\n")
}

func (self *Transpiler) writeToFile(content string) {
	self.printFuncName("")

	if self.testPrintMode {
		fmt.Print(content)
	} else {
		if self.outputFile != nil {
			self.outputFile.WriteString(content)
		}
	}
}

func (self *Transpiler) getPos(astNode ast.IStatement) string {
	return fmt.Sprintf("%s:%d:%d", self.filename, astNode.GetLn(), astNode.GetCol())
}

func (self *Transpiler) Transpile(astNode ast.IStatement, env *Environment, filename string) error {
	if !self.testMode {
		self.filename = filename
	}
	if filename != "" {
		self.outputFilename = filename + ".sh"
		f, err := os.Create(self.outputFilename)
		if err != nil {
			fmt.Println("Something went wrong creating the output file:", err)
		}
		defer f.Close()
		self.outputFile = f
	}

	_, err := self.transpile(astNode, env)
	if err != nil {
		return err
	}

	self.writeLnToFile("#!/bin/bash")
	self.writeLnToFile("# Created by Scrila Transpiler v0.0.1\n")

	if self.nativeFuncTranspilat != "" {
		self.writeLnToFile("# Native function implementations")
		self.writeToFile(self.nativeFuncTranspilat)
	}

	self.writeLnToFile("# User script")
	self.writeLnToFile(self.userScriptTranspilat)

	return nil
}

func (self *Transpiler) transpile(astNode ast.IStatement, env *Environment) (ast.IRuntimeVal, error) {
	switch astNode.GetKind() {
	// Handle Expressions
	case ast.IntLiteralNode:
		return NewIntVal(ast.ExprToIntLit(astNode).GetValue()), nil

	case ast.StrLiteralNode:
		return NewStrVal(ast.ExprToStrLit(astNode).GetValue()), nil

	case ast.IdentifierNode:
		return self.evalIdentifier(ast.ExprToIdent(astNode), env)

	case ast.ObjectLiteralNode:
		return self.evalObjectExpr(ast.ExprToObjLit(astNode), env)

	case ast.CallExprNode:
		return self.evalCallExpr(ast.ExprToCallExpr(astNode), env)

	case ast.AssignmentExprNode:
		return self.evalAssignment(ast.ExprToAssignmentExpr(astNode), env)

	case ast.BinaryExprNode:
		return self.evalBinaryExpr(ast.ExprToBinExpr(astNode), env)

	case ast.MemberExprNode:
		return self.evalMemberExpr(ast.ExprToMemberExpr(astNode), env)

	case ast.ReturnExprNode:
		return self.evalReturnExpr(ast.ExprToReturnExpr(astNode), env)

	// Handle Statements
	case ast.CommentNode:
		self.writeLnTranspilat("# " + ast.ExprToComment(astNode).GetComment())
		return NewNullVal(), nil

	case ast.ProgramNode:
		return self.evalProgram(ast.ExprToProgram(astNode), env)

	case ast.VarDeclarationNode:
		return self.evalVarDeclaration(ast.ExprToVarDecl(astNode), env)

	case ast.IfStatementNode:
		return self.evalIfStatement(ast.ExprToIfStmt(astNode), env)

	case ast.WhileStatementNode:
		return self.evalWhileStatement(ast.ExprToWhileStmt(astNode), env)

	case ast.FunctionDeclarationNode:
		return self.evalFunctionDeclaration(ast.ExprToFuncDecl(astNode), env)

	default:
		return NewNullVal(), fmt.Errorf("%s: This AST Node has not been setup for interpretion: %s", self.getPos(astNode), astNode)
	}
}

func (self *Transpiler) pushContext(context Context) {
	self.contexts = append(self.contexts, context)
}

func (self *Transpiler) popContext() {
	self.contexts = self.contexts[:len(self.contexts)-1]
}

func (self *Transpiler) currentContext() Context {
	return self.contexts[len(self.contexts)-1]
}

func (self *Transpiler) contextContains(context Context) bool {
	return slices.Contains(self.contexts, context)
}

func (self *Transpiler) indent(offset int) string {
	return strings.Repeat("\t", len(self.contexts)-1-offset)
}

func (self *Transpiler) printFuncName(msg string) {
	if self.showCallStack {
		pc, _, _, _ := runtime.Caller(1)
		funcName := runtime.FuncForPC(pc).Name()
		funcName = strings.Replace(funcName, "ScriLa/cmd/scrila/transpiler.(*Transpiler).", "", -1)

		if msg == "" {
			fmt.Printf("%s()\n", funcName)
		} else {
			fmt.Printf("%s(): %s\n", funcName, msg)
		}
	}
}
