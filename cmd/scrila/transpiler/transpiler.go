package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"fmt"
	"os"
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
	funcContext   bool
	currentFunc   IFunctionVal
}

func NewTranspiler() *Transpiler {
	return &Transpiler{
		usedNativeFunctions: []string{},
	}
}

func (self *Transpiler) writeLnToFile(content string) {
	self.writeToFile(content + "\n")
}

func (self *Transpiler) writeToFile(content string) {
	if self.testPrintMode {
		fmt.Print(content)
	} else {
		if self.outputFile != nil {
			self.outputFile.WriteString(content)
		}
	}
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

	self.writeLnToFile("#!/bin/bash")
	self.writeLnToFile("# Created by Scrila Transpiler v0.0.1")
	_, err := self.transpile(astNode, env)
	return err
}

func (self *Transpiler) transpile(astNode ast.IStatement, env *Environment) (IRuntimeVal, error) {
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
		self.writeLnToFile("# " + ast.ExprToComment(astNode).GetComment())
		return NewNullVal(), nil

	case ast.ProgramNode:
		return self.evalProgram(ast.ExprToProgram(astNode), env)

	case ast.VarDeclarationNode:
		return self.evalVarDeclaration(ast.ExprToVarDecl(astNode), env)

	case ast.IfStatementNode:
		return self.evalIfStatement(ast.ExprToIfStmt(astNode), env)

	case ast.FunctionDeclarationNode:
		return self.evalFunctionDeclaration(ast.ExprToFuncDecl(astNode), env)

	default:
		return NewNullVal(), fmt.Errorf("%s:%d:%d: This AST Node has not been setup for interpretion: %s", self.filename, astNode.GetLn(), astNode.GetCol(), astNode)
	}
}
