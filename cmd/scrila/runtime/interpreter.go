package runtime

import (
	"ScriLa/cmd/scrila/ast"
	"fmt"
	"os"
)

func Evaluate(astNode ast.IStatement) IRuntimeVal {
	switch astNode.GetKind() {
	case ast.IntLiteralNode:
		var i interface{} = astNode
		intLiteral, _ := i.(ast.IIntLiteral)
		// TODO error handling
		return NewIntVal(intLiteral.GetValue())

	case ast.NullLiteralNode:
		return NewNullVal()

	case ast.BinaryExprNode:
		var i interface{} = astNode
		binaryExpr, _ := i.(ast.IBinaryExpr)
		// TODO error handling
		return evalBinaryExpr(binaryExpr)

	case ast.ProgramNode:
		var i interface{} = astNode
		program, _ := i.(ast.IProgram)
		// TODO error handling
		return evalProgram(program)

	default:
		fmt.Println("This AST Node has not benn setup for interpretion:", astNode)
		os.Exit(1)
		return NewNullVal()
	}
}

func evalBinaryExpr(binOp ast.IBinaryExpr) IRuntimeVal {
	lhs := Evaluate(binOp.GetLeft())
	rhs := Evaluate(binOp.GetRight())

	if lhs.GetType() == IntValueType && rhs.GetType() == IntValueType {
		var i interface{} = lhs
		left, _ := i.(IIntVal)
		i = rhs
		right, _ := i.(IIntVal)
		return evalIntBinaryExpr(left, right, binOp.GetOperator())
	}

	// One or both are NULL, or another not yet supported type
	return NewNullVal()
}

func evalIntBinaryExpr(lhs IIntVal, rhs IIntVal, operator string) IIntVal {
	var result int64

	switch operator {
	case "+":
		result = lhs.GetValue() + rhs.GetValue()
	case "-":
		result = lhs.GetValue() - rhs.GetValue()
	case "*":
		result = lhs.GetValue() * rhs.GetValue()
	case "/":
		// TODO Division by zero
		result = lhs.GetValue() / rhs.GetValue()
	default:
		fmt.Println("Unsupported binary operator: ", operator)
		os.Exit(1)
	}

	return NewIntVal(result)
}

func evalProgram(program ast.IProgram) IRuntimeVal {
	var lastEvaluated IRuntimeVal = NewNullVal()

	for _, statement := range program.GetBody() {
		lastEvaluated = Evaluate(statement)
	}

	return lastEvaluated
}
