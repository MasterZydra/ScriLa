package bashTranspiler

import (
	"ScriLa/cmd/scrila/bashAst"
	"ScriLa/cmd/scrila/scrilaAst"
	"fmt"
)

func (self *Transpiler) exprToBashStmt(expr scrilaAst.IExpr, env *Environment) (bashAst.IStatement, error) {
	switch expr.GetKind() {
	case scrilaAst.BinaryExprNode:
		binaryBashStmt, ok := self.bashStmtStack[expr.GetId()]
		if !ok {
			return nil, fmt.Errorf("exprToBashStmt(): BinaryExpr is not stored in stack")
		}
		return binaryBashStmt, nil
	case scrilaAst.BoolLiteralNode:
		return bashAst.NewBoolLiteral(scrilaAst.ExprToBoolLit(expr).GetValue()), nil
	case scrilaAst.CallExprNode:
		returnVarName, err := self.getCallerResultVarName(scrilaAst.ExprToCallExpr(expr), env)
		if err != nil {
			return nil, err
		}
		scrilaReturnType, err := self.getFuncReturnType(scrilaAst.ExprToCallExpr(expr), env)
		if err != nil {
			return nil, err
		}
		bashReturnType, err := scrilaNodeTypeToBashNodeType(scrilaReturnType)
		if err != nil {
			return nil, err
		}
		return bashAst.NewVarLiteral(returnVarName, bashReturnType), nil
	case scrilaAst.IdentifierNode:
		varName := identNodeGetSymbol(expr)
		scrilaVarType, err := env.lookupVarType(varName)
		if err != nil {
			return nil, err
		}
		bashVarType, err := scrilaNodeTypeToBashNodeType(scrilaVarType)
		if err != nil {
			return nil, err
		}
		return bashAst.NewVarLiteral(varName, bashVarType), nil
	case scrilaAst.IntLiteralNode:
		return bashAst.NewIntLiteral(scrilaAst.ExprToIntLit(expr).GetValue()), nil
	case scrilaAst.StrLiteralNode:
		return bashAst.NewStrLiteral(scrilaAst.ExprToStrLit(expr).GetValue()), nil
	default:
		return nil, fmt.Errorf("exprToBashStmt: Expr of kind '%s' not implemented", expr.GetKind())
	}
}

func runtimeToBoolVal(runtimeVal scrilaAst.IRuntimeVal) IBoolVal {
	var i interface{} = runtimeVal
	return i.(IBoolVal)
}

func runtimeToIntVal(runtimeVal scrilaAst.IRuntimeVal) IIntVal {
	var i interface{} = runtimeVal
	return i.(IIntVal)
}

func runtimeToStrVal(runtimeVal scrilaAst.IRuntimeVal) IStrVal {
	var i interface{} = runtimeVal
	return i.(IStrVal)
}

func runtimeToObjVal(runtimeVal scrilaAst.IRuntimeVal) IObjVal {
	var i interface{} = runtimeVal
	return i.(IObjVal)
}

func runtimeToNativeFunc(runtimeVal scrilaAst.IRuntimeVal) INativeFunc {
	var i interface{} = runtimeVal
	return i.(INativeFunc)
}

func runtimeToFuncVal(runtimeVal scrilaAst.IRuntimeVal) IFunctionVal {
	var i interface{} = runtimeVal
	return i.(IFunctionVal)
}

var scrilaNodeTypeToBashNodeTypeMapping = map[scrilaAst.NodeType]bashAst.NodeType{
	scrilaAst.BoolLiteralNode: bashAst.BoolLiteralNode,
	scrilaAst.IntLiteralNode:  bashAst.IntLiteralNode,
	scrilaAst.StrLiteralNode:  bashAst.StrLiteralNode,
	scrilaAst.VoidNode:        bashAst.VoidNode,
}

func scrilaNodeTypeToBashNodeType(nodeType scrilaAst.NodeType) (bashAst.NodeType, error) {
	value, ok := scrilaNodeTypeToBashNodeTypeMapping[nodeType]
	if !ok {
		return "", fmt.Errorf("scrilaNodeTypeToBashNodeType(): Type '%s' is not in mapping", nodeType)
	}
	return value, nil
}

func bashNodeTypeToScrilaNodeType(nodeType bashAst.NodeType) (scrilaAst.NodeType, error) {
	for k, v := range scrilaNodeTypeToBashNodeTypeMapping {
		if v == nodeType {
			return k, nil
		}
	}
	return "", fmt.Errorf("scrilaNodeTypeToBashNodeType(): Type '%s' is not in mapping", nodeType)
}

var scrilaNodeTypeToRuntimeValMapping = map[scrilaAst.NodeType]scrilaAst.IRuntimeVal{
	scrilaAst.BoolLiteralNode: NewBoolVal(true),
	scrilaAst.IntLiteralNode:  NewIntVal(1),
	scrilaAst.StrLiteralNode:  NewStrVal("str"),
}

func scrilaNodeTypeToRuntimeVal(nodeType scrilaAst.NodeType) (scrilaAst.IRuntimeVal, error) {
	value, ok := scrilaNodeTypeToRuntimeValMapping[nodeType]
	if !ok {
		return NewNullVal(), fmt.Errorf("scrilaNodeTypeToRuntimeVal(): Type '%s' is not in mapping", nodeType)
	}
	return value, nil
}

var scrilaNodeTypeToTmpVarNameMapping = map[scrilaAst.NodeType]string{
	scrilaAst.BoolLiteralNode: "tmpBool",
	scrilaAst.IntLiteralNode:  "tmpInt",
	scrilaAst.StrLiteralNode:  "tmpStr",
}

func scrilaNodeTypeToTmpVarName(nodeType scrilaAst.NodeType) (string, error) {
	value, ok := scrilaNodeTypeToTmpVarNameMapping[nodeType]
	if !ok {
		return "", fmt.Errorf("scrilaNodeTypeToTmpVarName(): Node type '%s' is not in mapping", nodeType)
	}
	return value, nil
}
