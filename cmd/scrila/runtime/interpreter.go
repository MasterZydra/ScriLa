package runtime

import (
	"ScriLa/cmd/scrila/ast"
	"fmt"
)

func Evaluate(astNode ast.IStatement, env *Environment) (IRuntimeVal, error) {
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
		return NewNullVal(), nil

	case ast.ProgramNode:
		return evalProgram(ast.ExprToProgram(astNode), env)

	case ast.VarDeclarationNode:
		return evalVarDeclaration(ast.ExprToVarDecl(astNode), env)

	case ast.FunctionDeclarationNode:
		return evalFunctionDeclaration(ast.ExprToFuncDecl(astNode), env)

	default:
		return NewNullVal(), fmt.Errorf("Evaluate: This AST Node has not been setup for interpretion: %s", astNode)
	}
}
