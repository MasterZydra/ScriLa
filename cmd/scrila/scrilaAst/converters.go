package scrilaAst

import (
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

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

func ExprToWhileStmt(expr IExpr) IWhileStatement {
	var i interface{} = expr
	return i.(IWhileStatement)
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

func ExprToArray(expr IExpr) IArray {
	var i interface{} = expr
	return i.(IArray)
}

func ExprToBoolLit(expr IExpr) IBoolLiteral {
	var i interface{} = expr
	return i.(IBoolLiteral)
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

var valueTypeToArrayMapping = map[ValueType]ValueType{
	BoolValueType: BoolArrayValueType,
	IntValueType:  IntArrayValueType,
	StrValueType:  StrArrayValueType,
}

func ValueTypeToArrayType(valueType ValueType) (ValueType, error) {
	value, ok := valueTypeToArrayMapping[valueType]
	if !ok {
		return NullValueType, fmt.Errorf("ValueTypeToArray(): Type '%s' is not in mapping", valueType)
	}
	return value, nil
}

func ArrayTypeToValueType(arrayType ValueType) (ValueType, error) {
	for k, v := range valueTypeToArrayMapping {
		if v == arrayType {
			return k, nil
		}
	}
	return "", fmt.Errorf("ArrayTypeToValueType(): Type '%s' is not in mapping", arrayType)
}

var dataTypeToArrayMapping = map[NodeType]NodeType{
	BoolLiteralNode: BoolArrayNode,
	IntLiteralNode:  IntArrayNode,
	StrLiteralNode:  StrArrayNode,
}

func DataTypeToArrayType(dataType NodeType) (NodeType, error) {
	value, ok := dataTypeToArrayMapping[dataType]
	if !ok {
		return ProgramNode, fmt.Errorf("DataTypeToArray(): Type '%s' is not in mapping", dataType)
	}
	return value, nil
}

func ArrayTypeToDataType(arrayType NodeType) (NodeType, error) {
	for k, v := range dataTypeToArrayMapping {
		if v == arrayType {
			return k, nil
		}
	}
	return "", fmt.Errorf("ArrayTypeToDataType(): Type '%s' is not in mapping", arrayType)
}

var ComparisonOps = []string{"<", ">", "<=", ">=", "!=", "=="}

func BinExprIsComp(binOp IBinaryExpr) bool {
	return slices.Contains(ComparisonOps, binOp.GetOperator())
}

var BooleanOps = []string{"||", "&&"}

func BinExprIsBoolOp(binOp IBinaryExpr) bool {
	return slices.Contains(BooleanOps, binOp.GetOperator())
}

func BinExprReturnsBool(binOp IBinaryExpr) bool {
	return BinExprIsBoolOp(binOp) || BinExprIsComp(binOp)
}

var indentDepth int = 0

func indent() string {
	return strings.Repeat("  ", indentDepth+1)
}

func SprintAST(program IProgram) string {
	astString := ""
	for _, stmt := range program.GetBody() {
		astString += fmt.Sprintf("%s%s\n", indent(), stmt)
	}
	return astString
}
