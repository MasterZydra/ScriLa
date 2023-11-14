package bashTranspiler

import (
	"ScriLa/cmd/scrila/bashAst"
	"ScriLa/cmd/scrila/scrilaAst"
	"fmt"
)

func (self *Transpiler) exprIsType(expr scrilaAst.IExpr, wantedType scrilaAst.NodeType, env *Environment) (bool, scrilaAst.NodeType, error) {
	givenType := expr.GetKind()
	// Check types directly
	if givenType == wantedType {
		return true, givenType, nil
	}

	// If the wanted type is an Identifier the following checks make no sens
	if wantedType == scrilaAst.IdentifierNode {
		return false, givenType, nil
	}

	// Check if identifier is variable and if that variable type matches with the wanted type
	if givenType == scrilaAst.IdentifierNode {
		givenType, err := env.lookupVarType(identNodeGetSymbol(expr))
		if err != nil {
			return false, givenType, err
		}

		return givenType == wantedType, givenType, nil
	}

	// Check if the return type of the function call matches with the wanted type
	if givenType == scrilaAst.CallExprNode {
		givenType, err := self.getFuncReturnType(scrilaAst.ExprToCallExpr(expr), env)
		if err != nil {
			return false, givenType, err
		}
		return givenType == wantedType, givenType, nil
	}

	// Check if the return type of a binary expression matches with the wanted type
	if givenType == scrilaAst.BinaryExprNode {
		bashStmt, ok := self.bashStmtStack[expr.GetId()]
		if !ok {
			return false, givenType, fmt.Errorf("exprIsType(): BinaryExpr is not stored in stack")
		}

		binOpExpr := bashAst.StmtToBinaryOpExpr(bashStmt)
		if binOpExpr.GetKind() == bashAst.BinaryCompExprNode {
			return true, wantedType, nil
		}

		givenType, err := bashNodeTypeToScrilaNodeType(binOpExpr.GetOpType())
		if err != nil {
			return false, givenType, err
		}

		return givenType == wantedType, givenType, nil
	}

	return false, givenType, nil
}