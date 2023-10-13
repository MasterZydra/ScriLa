package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"fmt"
	"os"
)

var fileName string
var outputFileName string
var outputFile *os.File

var testMode bool
var testPrintMode bool
var funcContext bool
var currentFunc IFunctionVal

func writeLnToFile(content string) {
	writeToFile(content + "\n")
}

func writeToFile(content string) {
	if testPrintMode {
		fmt.Print(content)
	} else {
		if outputFile != nil {
			outputFile.WriteString(content)
		}
	}
}

func Transpile(astNode ast.IStatement, env *Environment, filename string) error {
	if !testMode {
		fileName = filename
	}
	if filename != "" {
		outputFileName = filename + ".sh"
		f, err := os.Create(outputFileName)
		if err != nil {
			fmt.Println("Something went wrong creating the output file:", err)
		}
		defer f.Close()
		outputFile = f
	}

	writeLnToFile("#!/bin/bash")
	writeLnToFile("# Created by Scrila Transpiler v0.0.1")
	_, err := transpile(astNode, env)
	return err
}

func transpile(astNode ast.IStatement, env *Environment) (IRuntimeVal, error) {
	switch astNode.GetKind() {
	// Handle Expressions
	case ast.IntLiteralNode:
		return NewIntVal(ast.ExprToIntLit(astNode).GetValue()), nil

	case ast.StrLiteralNode:
		return NewStrVal(ast.ExprToStrLit(astNode).GetValue()), nil

	case ast.IdentifierNode:
		return evalIdentifier(ast.ExprToIdent(astNode), env)

	case ast.ObjectLiteralNode:
		return evalObjectExpr(ast.ExprToObjLit(astNode), env)

	case ast.CallExprNode:
		return evalCallExpr(ast.ExprToCallExpr(astNode), env)

	case ast.AssignmentExprNode:
		return evalAssignment(ast.ExprToAssignmentExpr(astNode), env)

	case ast.BinaryExprNode:
		return evalBinaryExpr(ast.ExprToBinExpr(astNode), env)

	case ast.MemberExprNode:
		return evalMemberExpr(ast.ExprToMemberExpr(astNode), env)

	case ast.ReturnExprNode:
		return evalReturnExpr(ast.ExprToReturnExpr(astNode), env)

	// Handle Statements
	case ast.CommentNode:
		writeLnToFile("# " + ast.ExprToComment(astNode).GetComment())
		return NewNullVal(), nil

	case ast.ProgramNode:
		return evalProgram(ast.ExprToProgram(astNode), env)

	case ast.VarDeclarationNode:
		return evalVarDeclaration(ast.ExprToVarDecl(astNode), env)

	case ast.IfStatementNode:
		return evalIfStatement(ast.ExprToIfStmt(astNode), env)

	case ast.FunctionDeclarationNode:
		return evalFunctionDeclaration(ast.ExprToFuncDecl(astNode), env)

	default:
		return NewNullVal(), fmt.Errorf("%s:%d:%d: This AST Node has not been setup for interpretion: %s", fileName, astNode.GetLn(), astNode.GetCol(), astNode)
	}
}
