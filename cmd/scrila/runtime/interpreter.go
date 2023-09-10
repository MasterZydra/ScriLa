package runtime

import (
	"ScriLa/cmd/scrila/ast"
	"fmt"
	"os"
)

func Evaluate(astNode ast.IStatement, env *Environment) IRuntimeVal {
	switch astNode.GetKind() {
	// Handle Expressions

	case ast.IntLiteralNode:
		var i interface{} = astNode
		intLiteral, _ := i.(ast.IIntLiteral)
		// TODO error handling
		return NewIntVal(intLiteral.GetValue())

	case ast.IdentifierNode:
		var i interface{} = astNode
		identifier, _ := i.(ast.IIdentifier)
		// TODO error handling
		return evalIdentifier(identifier, env)

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

	default:
		fmt.Println("This AST Node has not been setup for interpretion:", astNode)
		os.Exit(1)
		return NewNullVal()
	}
}
