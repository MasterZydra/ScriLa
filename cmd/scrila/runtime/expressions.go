package runtime

import (
	"ScriLa/cmd/scrila/ast"
	"fmt"
	"os"
)

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

	if lhs.GetType() == StrValueType && rhs.GetType() == StrValueType {
		var i interface{} = lhs
		left, _ := i.(IStrVal)
		i = rhs
		right, _ := i.(IStrVal)
		return evalStrBinaryExpr(left, right, binOp.GetOperator())
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

func evalStrBinaryExpr(lhs IStrVal, rhs IStrVal, operator string) IStrVal {
	var result string

	switch operator {
	case "+":
		result = lhs.GetValue() + rhs.GetValue()
	default:
		fmt.Println("Unsupported binary operator: ", operator)
		os.Exit(1)
	}

	return NewStrVal(result)
}

func evalAssignment(assignment ast.IAssignmentExpr, env *Environment) IRuntimeVal {
	if assignment.GetAssigne().GetKind() == ast.MemberExprNode {
		return evalAssignmentObjMember(assignment, env)
	}

	if assignment.GetAssigne().GetKind() != ast.IdentifierNode {
		fmt.Println("evalAssignment: Invalid LHS inside assignment expr", assignment.GetAssigne())
		os.Exit(1)
	}

	var i interface{} = assignment.GetAssigne()
	assigne, _ := i.(ast.IIdentifier)
	return env.assignVar(assigne.GetSymbol(), Evaluate(assignment.GetValue(), env))
}

func evalAssignmentObjMember(assignment ast.IAssignmentExpr, env *Environment) IRuntimeVal {
	if assignment.GetAssigne().GetKind() != ast.MemberExprNode {
		fmt.Println("evalAssignmentObjMember: Invalid LHS inside assignment expr", assignment.GetAssigne())
		os.Exit(1)
	}

	var i interface{} = assignment.GetAssigne()
	memberExpr, _ := i.(ast.IMemberExpr)

	if memberExpr.GetObject().GetKind() != ast.IdentifierNode {
		fmt.Println("evalMemberExpr: Object - Node kind '" + memberExpr.GetObject().GetKind() + "' not supported")
		os.Exit(1)
	}

	if memberExpr.GetProperty().GetKind() != ast.IdentifierNode {
		fmt.Println("evalMemberExpr: Property - Node kind '" + memberExpr.GetProperty().GetKind() + "' not supported")
		os.Exit(1)
	}

	i = memberExpr.GetObject()
	identifier, _ := i.(ast.IIdentifier)
	obj := env.lookupVar(identifier.GetSymbol())
	if obj.GetType() != ObjValueType {
		fmt.Println("evalMemberExpr: variable '" + identifier.GetSymbol() + "' is not of type 'object'")
		os.Exit(1)
	}

	i = obj
	objVal, _ := i.(IObjVal)

	i = memberExpr.GetProperty()
	property, _ := i.(ast.IIdentifier)

	value := Evaluate(assignment.GetValue(), env)
	objVal.GetProperties()[property.GetSymbol()] = value
	return value
}

func evalObjectExpr(object ast.IObjectLiteral, env *Environment) IRuntimeVal {
	obj := NewObjVal()

	for _, property := range object.GetProperties() {
		obj.properties[property.GetKey()] = Evaluate(property.GetValue(), env)
	}

	return obj
}

func evalMemberExpr(memberExpr ast.IMemberExpr, env *Environment) IRuntimeVal {
	if memberExpr.GetObject().GetKind() != ast.IdentifierNode {
		fmt.Println("evalMemberExpr: Object - Node kind '" + memberExpr.GetObject().GetKind() + "' not supported")
		os.Exit(1)
	}

	if memberExpr.GetProperty().GetKind() != ast.IdentifierNode {
		fmt.Println("evalMemberExpr: Property - Node kind '" + memberExpr.GetProperty().GetKind() + "' not supported")
		os.Exit(1)
	}

	var i interface{} = memberExpr.GetObject()
	identifier, _ := i.(ast.IIdentifier)
	obj := env.lookupVar(identifier.GetSymbol())
	if obj.GetType() != ObjValueType {
		fmt.Println("evalMemberExpr: variable '" + identifier.GetSymbol() + "' is not of type 'object'")
		os.Exit(1)
	}

	i = obj
	objVal, _ := i.(IObjVal)

	i = memberExpr.GetProperty()
	property, _ := i.(ast.IIdentifier)

	return objVal.GetProperties()[property.GetSymbol()]
}

func evalCallExpr(call ast.ICallExpr, env *Environment) IRuntimeVal {
	// TODO add helpers? https://zetcode.com/golang/filter-map/
	var args []IRuntimeVal
	for _, arg := range call.GetArgs() {
		args = append(args, Evaluate(arg, env))
	}

	caller := Evaluate(call.GetCaller(), env)

	switch caller.GetType() {
	case NativeFnType:
		var i interface{} = caller
		fn, _ := i.(INativeFunc)
		return fn.GetCall()(args, env)

	case FunctionValueType:
		var i interface{} = caller
		fn, _ := i.(IFunctionVal)
		scope := NewEnvironment(fn.GetDeclarationEnv())

		// Create variables for the parameters list
		for i := 0; i < len(fn.GetParams()); i++ {
			// TODO Check the bounds here. Verify arity of function.
			// Which means: len(fn.GetParams()) == len(args)
			scope.declareVar(fn.GetParams()[i], args[i], false)
		}

		var result IRuntimeVal
		result = NewNullVal()
		// Evaluate the function body line by line
		for _, stmt := range fn.GetBody() {
			result = Evaluate(stmt, scope)
		}
		return result

	default:
		fmt.Println("Cannot call value that is not a function:", caller)
		os.Exit(1)
		return NewNullVal()
	}
}
