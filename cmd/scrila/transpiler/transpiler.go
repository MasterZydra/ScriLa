package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"fmt"
	"os"
)

var outputFileName string
var outputFile *os.File

var testMode bool

func writeLnToFile(content string) {
	writeToFile(content + "\n")
}

func writeToFile(content string) {
	if testMode {
		fmt.Print(content)
	} else {
		outputFile.WriteString(content)
	}
}

func Transpile(astNode ast.IStatement, env *Environment, fileName string) {
	outputFileName = fileName + ".sh"
	f, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Something went wrong creating the output file:", err)
	}
	defer f.Close()
	outputFile = f

	writeLnToFile("#!/bin/bash")
	writeLnToFile("# Created by Scrila Transpiler v0.0.1")
	transpile(astNode, env)
}

func transpile(astNode ast.IStatement, env *Environment) IRuntimeVal {
	switch astNode.GetKind() {
	// Handle Expressions

	case ast.IntLiteralNode:
		var i interface{} = astNode
		intLiteral, _ := i.(ast.IIntLiteral)
		// TODO error handling
		return NewIntVal(intLiteral.GetValue())

	case ast.StrLiteralNode:
		var i interface{} = astNode
		strLiteral, _ := i.(ast.IStrLiteral)
		// TODO error handling
		return NewStrVal(strLiteral.GetValue())

	case ast.IdentifierNode:
		var i interface{} = astNode
		identifier, _ := i.(ast.IIdentifier)
		// TODO error handling
		return evalIdentifier(identifier, env)

	case ast.ObjectLiteralNode:
		var i interface{} = astNode
		object, _ := i.(ast.IObjectLiteral)
		// TODO error handling
		return evalObjectExpr(object, env)

	case ast.CallExprNode:
		var i interface{} = astNode
		call, _ := i.(ast.ICallExpr)
		// TODO error handling
		return evalCallExpr(call, env)

	case ast.AssignmentExprNode:
		var i interface{} = astNode
		assignment, _ := i.(ast.IAssignmentExpr)
		// TODO error handling
		return evalAssignment(assignment, env)

	case ast.BinaryExprNode:
		var i interface{} = astNode
		binaryExpr, _ := i.(ast.IBinaryExpr)
		// TODO error handling
		return evalBinaryExpr(binaryExpr, env)

	case ast.MemberExprNode:
		var i interface{} = astNode
		memberExpr, _ := i.(ast.IMemberExpr)
		// TODO error handling
		return evalMemberExpr(memberExpr, env)

	// Handle Statements

	case ast.ProgramNode:
		var i interface{} = astNode
		program, _ := i.(ast.IProgram)
		// TODO error handling
		return evalProgram(program, env)

	case ast.VarDeclarationNode:
		var i interface{} = astNode
		varDeclaration, _ := i.(ast.IVarDeclaration)
		// TODO error handling
		return evalVarDeclaration(varDeclaration, env)

	case ast.FunctionDeclarationNode:
		var i interface{} = astNode
		funcDeclaration, _ := i.(ast.IFunctionDeclaration)
		// TODO error handling
		return evalFunctionDeclaration(funcDeclaration, env)

	default:
		fmt.Println("This AST Node has not been setup for interpretion:", astNode)
		os.Exit(1)
		return NewNullVal()
	}
}
