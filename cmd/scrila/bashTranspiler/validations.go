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

	// Check if the return type of a member expression matches with the wanted type
	if givenType == scrilaAst.MemberExprNode {
		bashStmt, ok := self.bashStmtStack[expr.GetId()]
		if !ok {
			return false, givenType, fmt.Errorf("exprIsType(): MemberExpr is not stored in stack")
		}
		memberExpr := bashAst.StmtToMemberExpr(bashStmt)
		varType, err := env.lookupVarType(memberExpr.GetVarname().GetValue())
		if err != nil {
			return false, givenType, err
		}

		varType, err = scrilaAst.ArrayTypeToDataType(varType)
		if err != nil {
			return false, givenType, err
		}

		return varType == wantedType, varType, nil
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

		givenType, err := bashNodeTypeToScrilaNodeType(binOpExpr.GetDataType())
		if err != nil {
			return false, givenType, err
		}

		return givenType == wantedType, givenType, nil
	}

	// Check array data type
	if givenType == scrilaAst.ArrayLiteralNode {
		array, ok := self.bashStmtStack[expr.GetId()]
		if !ok {
			return false, givenType, fmt.Errorf("exprIsType(): Array is not stored in stack")
		}
		bashArray := bashAst.StmtToArray(array)

		scrilaWantedDataType, err := scrilaAst.ArrayTypeToDataType(wantedType)
		if err != nil {
			return false, givenType, err
		}

		if bashArray.GetDataType() == bashAst.VoidNode {
			bashDataType, err := scrilaNodeTypeToBashNodeType(scrilaWantedDataType)
			if err != nil {
				return false, givenType, err
			}
			bashArray.SetDataType(bashDataType)
			return true, wantedType, nil
		}

		arrayDataType, err := bashNodeTypeToScrilaNodeType(bashArray.GetDataType())
		if err != nil {
			return false, givenType, err
		}
		return scrilaWantedDataType == arrayDataType, arrayDataType, nil
	}

	return false, givenType, nil
}
