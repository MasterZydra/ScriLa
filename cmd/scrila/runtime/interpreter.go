package runtime

import (
	"ScriLa/cmd/scrila/ast"
	"fmt"
	"os"
)

func Evaluate(astNode ast.IStatement, env *Environment) IRuntimeVal {
	switch astNode.GetKind() {
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

	case ast.BinaryExprNode:
		var i interface{} = astNode
		binaryExpr, _ := i.(ast.IBinaryExpr)
		// TODO error handling
		return evalBinaryExpr(binaryExpr, env)

	case ast.ProgramNode:
		var i interface{} = astNode
		program, _ := i.(ast.IProgram)
		// TODO error handling
		return evalProgram(program, env)

	default:
		fmt.Println("This AST Node has not benn setup for interpretion:", astNode)
		os.Exit(1)
		return NewNullVal()
	}
}

func evalIdentifier(identifier ast.IIdentifier, env *Environment) IRuntimeVal {
	return env.lookupVar(identifier.GetSymbol())
}

func evalBinaryExpr(binOp ast.IBinaryExpr, env *Environment) IRuntimeVal {
	lhs := Evaluate(binOp.GetLeft(), env)
	rhs := Evaluate(binOp.GetRight(), env)

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

func evalProgram(program ast.IProgram, env *Environment) IRuntimeVal {
	var lastEvaluated IRuntimeVal = NewNullVal()

	for _, statement := range program.GetBody() {
		lastEvaluated = Evaluate(statement, env)
	}

	return lastEvaluated
}
