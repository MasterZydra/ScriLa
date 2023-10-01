package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"fmt"
	"os"
)

var outputFileName string
var outputFile *os.File

var testMode bool
var funcContext bool

func writeLnToFile(content string) {
	writeToFile(content + "\n")
}

func writeToFile(content string) {
	if testMode {
		fmt.Print(content)
	} else {
		if outputFile != nil {
			outputFile.WriteString(content)
		}
	}
}

func Transpile(astNode ast.IStatement, env *Environment, fileName string) error {
	if fileName != "" {
		outputFileName = fileName + ".sh"
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
		var i interface{} = astNode
		intLiteral, ok := i.(ast.IIntLiteral)
		if !ok {
			return NewNullVal(), fmt.Errorf("Evaluate: Failed to convert Statement to IntLiteral")
		}
		return NewIntVal(intLiteral.GetValue()), nil

	case ast.StrLiteralNode:
		var i interface{} = astNode
		strLiteral, ok := i.(ast.IStrLiteral)
		if !ok {
			return NewNullVal(), fmt.Errorf("Evaluate: Failed to convert Statement to StrLiteral")
		}
		return NewStrVal(strLiteral.GetValue()), nil

	case ast.IdentifierNode:
		var i interface{} = astNode
		identifier, ok := i.(ast.IIdentifier)
		if !ok {
			return NewNullVal(), fmt.Errorf("Evaluate: Failed to convert Statement to Identifier")
		}
		return evalIdentifier(identifier, env)

	case ast.ObjectLiteralNode:
		var i interface{} = astNode
		object, ok := i.(ast.IObjectLiteral)
		if !ok {
			return NewNullVal(), fmt.Errorf("Evaluate: Failed to convert Statement to ObjectLiteral")
		}
		return evalObjectExpr(object, env)

	case ast.CallExprNode:
		var i interface{} = astNode
		call, ok := i.(ast.ICallExpr)
		if !ok {
			return NewNullVal(), fmt.Errorf("Evaluate: Failed to convert Statement to CallExpr")
		}
		return evalCallExpr(call, env)

	case ast.AssignmentExprNode:
		var i interface{} = astNode
		assignment, ok := i.(ast.IAssignmentExpr)
		if !ok {
			return NewNullVal(), fmt.Errorf("Evaluate: Failed to convert Statement to AssignmentExpr")
		}
		return evalAssignment(assignment, env)

	case ast.BinaryExprNode:
		var i interface{} = astNode
		binaryExpr, ok := i.(ast.IBinaryExpr)
		if !ok {
			return NewNullVal(), fmt.Errorf("Evaluate: Failed to convert Statement to BinaryExpr")
		}
		return evalBinaryExpr(binaryExpr, env)

	case ast.MemberExprNode:
		var i interface{} = astNode
		memberExpr, ok := i.(ast.IMemberExpr)
		if !ok {
			return NewNullVal(), fmt.Errorf("Evaluate: Failed to convert Statement to MemberExpr")
		}
		return evalMemberExpr(memberExpr, env)

	// Handle Statements
	case ast.CommentNode:
		var i interface{} = astNode
		comment, ok := i.(ast.IComment)
		if !ok {
			return NewNullVal(), fmt.Errorf("Evaluate: Failed to convert Statement to Comment")
		}
		writeLnToFile("# " + comment.GetComment())
		return NewNullVal(), nil

	case ast.ProgramNode:
		var i interface{} = astNode
		program, ok := i.(ast.IProgram)
		if !ok {
			return NewNullVal(), fmt.Errorf("Evaluate: Failed to convert Statement to Program")
		}
		return evalProgram(program, env)

	case ast.VarDeclarationNode:
		var i interface{} = astNode
		varDeclaration, ok := i.(ast.IVarDeclaration)
		if !ok {
			return NewNullVal(), fmt.Errorf("Evaluate: Failed to convert Statement to VarDeclaration")
		}
		return evalVarDeclaration(varDeclaration, env)

	case ast.FunctionDeclarationNode:
		var i interface{} = astNode
		funcDeclaration, ok := i.(ast.IFunctionDeclaration)
		if !ok {
			return NewNullVal(), fmt.Errorf("Evaluate: Failed to convert Statement to FunctionDeclaration")
		}
		return evalFunctionDeclaration(funcDeclaration, env)

	default:
		return NewNullVal(), fmt.Errorf("This AST Node has not been setup for interpretion: %s", astNode)
	}
}
