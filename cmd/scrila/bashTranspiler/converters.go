package bashTranspiler

import (
	"ScriLa/cmd/scrila/bashAst"
	"ScriLa/cmd/scrila/scrilaAst"
	"fmt"
)

// Returns the symbol of the given expr of kind Identifier
func identNodeGetSymbol(expr scrilaAst.IExpr) string {
	return scrilaAst.ExprToIdent(expr).GetSymbol()
}

func (self *Transpiler) exprToRhsBashStmt(expr scrilaAst.IExpr, env *Environment) (bashAst.IStatement, error) {
	bashStmt, err := self.exprToBashStmt(expr, env)
	if err != nil {
		return nil, err
	}

	// A comparison must be converted into an if statement
	if bashStmt.GetKind() == bashAst.BinaryCompExprNode ||
		(bashStmt.GetKind() == bashAst.BinaryOpExprNode && bashAst.StmtToBinaryOpExpr(bashStmt).GetDataType() == bashAst.BoolLiteralNode) {
		varname := fmt.Sprintf("tmpBools[%d]", self.currentCallArgIndex())
		if self.contextContains(FunctionContext) {
			varname = "tmpBools[${tmpIndex}]"
		}
		ifStmt := bashAst.NewIfStmt(bashStmt)
		ifStmt.AppendBody(bashAst.NewBashStmt(fmt.Sprintf("%s=\"true\"", varname)))
		elseStmt := bashAst.NewIfStmt(nil)
		elseStmt.AppendBody(bashAst.NewBashStmt(fmt.Sprintf("%s=\"false\"", varname)))
		ifStmt.SetElse(elseStmt)
		self.appendUserBody(ifStmt)
		bashStmt = bashAst.NewVarLiteral(varname, bashAst.BoolLiteralNode)
		if len(self.callArgIndexStack) > 0 {
			self.incCallArgIndex()
			self.setCallArgIndex()
		}
	}

	return bashStmt, nil
}

func (self *Transpiler) exprToBashStmt(expr scrilaAst.IExpr, env *Environment) (bashAst.IStatement, error) {
	switch expr.GetKind() {
	case scrilaAst.ArrayLiteralNode, scrilaAst.BinaryExprNode, scrilaAst.MemberExprNode:
		bashArray, ok := self.bashStmtStack[expr.GetId()]
		if !ok {
			return nil, fmt.Errorf("exprToBashStmt(): %s is not stored in stack", expr.GetKind())
		}
		return bashArray, nil
	case scrilaAst.BoolLiteralNode:
		return bashAst.NewBoolLiteral(scrilaAst.ExprToBoolLit(expr).GetValue()), nil
	case scrilaAst.BreakExprNode:
		return bashAst.NewBreakExpr(), nil
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
	case scrilaAst.ContinueExprNode:
		return bashAst.NewContinueExpr(), nil
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

func runtimeToNativeFunc(runtimeVal scrilaAst.IRuntimeVal) INativeFunc {
	var i interface{} = runtimeVal
	return i.(INativeFunc)
}

func runtimeToFuncVal(runtimeVal scrilaAst.IRuntimeVal) IFunctionVal {
	var i interface{} = runtimeVal
	return i.(IFunctionVal)
}

var scrilaNodeTypeToBashNodeTypeMapping = map[scrilaAst.NodeType]bashAst.NodeType{
	scrilaAst.BoolArrayNode:   bashAst.ArrayLiteralNode,
	scrilaAst.BoolLiteralNode: bashAst.BoolLiteralNode,
	scrilaAst.IntArrayNode:    bashAst.ArrayLiteralNode,
	scrilaAst.IntLiteralNode:  bashAst.IntLiteralNode,
	scrilaAst.StrArrayNode:    bashAst.ArrayLiteralNode,
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
	scrilaAst.BoolArrayNode:   NewArrayVal(scrilaAst.BoolArrayValueType),
	scrilaAst.BoolLiteralNode: NewBoolVal(true),
	scrilaAst.IntArrayNode:    NewArrayVal(scrilaAst.IntArrayValueType),
	scrilaAst.IntLiteralNode:  NewIntVal(1),
	scrilaAst.VoidNode:        NewNullVal(),
	scrilaAst.StrArrayNode:    NewArrayVal(scrilaAst.StrArrayValueType),
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
	scrilaAst.BoolLiteralNode: "tmpBools",
	scrilaAst.IntLiteralNode:  "tmpInts",
	scrilaAst.StrLiteralNode:  "tmpStrs",
}

func (self *Transpiler) scrilaNodeTypeToTmpVarName(nodeType scrilaAst.NodeType) (string, error) {
	value, ok := scrilaNodeTypeToTmpVarNameMapping[nodeType]
	if !ok {
		return "", fmt.Errorf("scrilaNodeTypeToTmpVarName(): Node type '%s' is not in mapping", nodeType)
	}
	index := self.currentCallArgIndex() - 1
	if index < 0 {
		index = 0
	}
	return fmt.Sprintf("%s[%d]", value, index), nil
}

func (self *Transpiler) scrilaNodeTypeToDynTmpVarName(nodeType scrilaAst.NodeType) (string, error) {
	value, ok := scrilaNodeTypeToTmpVarNameMapping[nodeType]
	if !ok {
		return "", fmt.Errorf("scrilaNodeTypeToDynTmpVarName(): Node type '%s' is not in mapping", nodeType)
	}
	return fmt.Sprintf("%s[${tmpIndex}]", value), nil
}

func (self *Transpiler) getFuncReturnType(call scrilaAst.ICallExpr, env *Environment) (scrilaAst.NodeType, error) {
	self.printFuncName("")

	if call.GetCaller().GetKind() != scrilaAst.IdentifierNode {
		return "", fmt.Errorf("%s: Function name must be an identifier. Got: '%s'", self.getPos(call.GetCaller()), call.GetCaller().GetKind())
	}

	funcName := identNodeGetSymbol(call.GetCaller())
	caller, err := env.lookupFunc(funcName)
	if err != nil {
		return "", err
	}

	switch caller.GetType() {
	case scrilaAst.FunctionValueType:
		return runtimeToFuncVal(caller).GetReturnType(), nil
	case scrilaAst.NativeFnType:
		return runtimeToNativeFunc(caller).GetReturnType(), nil
	default:
		return "", fmt.Errorf("%s: Cannot call value that is not a function: %s", self.getPos(call), caller.GetType())
	}
}

func (self *Transpiler) getCallerResultVarName(call scrilaAst.ICallExpr, env *Environment) (string, error) {
	self.printFuncName("")

	returnType, err := self.getFuncReturnType(call, env)
	if err != nil {
		return "", err
	}

	if returnType == scrilaAst.VoidNode {
		return "", fmt.Errorf("%s: Func '%s' does not have a return value", self.getPos(call.GetCaller()), identNodeGetSymbol(call.GetCaller()))
	}

	resultVarName, err := self.scrilaNodeTypeToTmpVarName(returnType)
	if err != nil {
		return "", err
	}
	return resultVarName, nil
}
