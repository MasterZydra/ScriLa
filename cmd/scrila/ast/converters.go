package ast

func ExprToProgram(expr IExpr) IProgram {
	var i interface{} = expr
	return i.(IProgram)
}

func ExprToIdent(expr IExpr) IIdentifier {
	var i interface{} = expr
	return i.(IIdentifier)
}

func ExprToComment(expr IExpr) IComment {
	var i interface{} = expr
	return i.(IComment)
}

func ExprToVarDecl(expr IExpr) IVarDeclaration {
	var i interface{} = expr
	return i.(IVarDeclaration)
}

func ExprToIfStmt(expr IExpr) IIfStatement {
	var i interface{} = expr
	return i.(IIfStatement)
}

func ExprToFuncDecl(expr IExpr) IFunctionDeclaration {
	var i interface{} = expr
	return i.(IFunctionDeclaration)
}

func ExprToAssignmentExpr(expr IExpr) IAssignmentExpr {
	var i interface{} = expr
	return i.(IAssignmentExpr)
}

func ExprToReturnExpr(expr IExpr) IReturnExpr {
	var i interface{} = expr
	return i.(IReturnExpr)
}

func ExprToBinExpr(expr IExpr) IBinaryExpr {
	var i interface{} = expr
	return i.(IBinaryExpr)
}

func ExprToCallExpr(expr IExpr) ICallExpr {
	var i interface{} = expr
	return i.(ICallExpr)
}

func ExprToMemberExpr(expr IExpr) IMemberExpr {
	var i interface{} = expr
	memberExpr, _ := i.(IMemberExpr)
	return memberExpr
}

func ExprToIntLit(expr IExpr) IIntLiteral {
	var i interface{} = expr
	return i.(IIntLiteral)
}

func ExprToStrLit(expr IExpr) IStrLiteral {
	var i interface{} = expr
	return i.(IStrLiteral)
}

func ExprToObjLit(expr IExpr) IObjectLiteral {
	var i interface{} = expr
	return i.(IObjectLiteral)
}
