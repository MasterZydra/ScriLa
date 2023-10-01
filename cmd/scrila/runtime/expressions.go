package runtime

import (
	"ScriLa/cmd/scrila/ast"
	"fmt"
)

func evalIdentifier(identifier ast.IIdentifier, env *Environment) (IRuntimeVal, error) {
	return env.lookupVar(identifier.GetSymbol())
}

func evalBinaryExpr(binOp ast.IBinaryExpr, env *Environment) (IRuntimeVal, error) {
	lhs, lhsError := Evaluate(binOp.GetLeft(), env)
	if lhsError != nil {
		return NewNullVal(), lhsError
	}
	rhs, rhsError := Evaluate(binOp.GetRight(), env)
	if rhsError != nil {
		return NewNullVal(), rhsError
	}

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

	return NewNullVal(), fmt.Errorf("evalBinaryExpr: Give types not supported (lhs: %s, rhs: %s)", lhs, rhs)
}

func evalIntBinaryExpr(lhs IIntVal, rhs IIntVal, operator string) (IIntVal, error) {
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
		return NewIntVal(0), fmt.Errorf("evalIntBinaryExpr: Unsupported binary operator: %s", operator)
	}

	return NewIntVal(result), nil
}

func evalStrBinaryExpr(lhs IStrVal, rhs IStrVal, operator string) (IStrVal, error) {
	var result string

	switch operator {
	case "+":
		result = lhs.GetValue() + rhs.GetValue()
	default:
		return NewStrVal(""), fmt.Errorf("evalStrBinaryExpr: Unsupported binary operator: %s", operator)
	}

	return NewStrVal(result), nil
}

func evalAssignment(assignment ast.IAssignmentExpr, env *Environment) (IRuntimeVal, error) {
	if assignment.GetAssigne().GetKind() == ast.MemberExprNode {
		return evalAssignmentObjMember(assignment, env)
	}

	if assignment.GetAssigne().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("evalAssignment: Invalid LHS inside assignment expr %s", assignment.GetAssigne())
	}

	var i interface{} = assignment.GetAssigne()
	assigne, _ := i.(ast.IIdentifier)
	value, err := Evaluate(assignment.GetValue(), env)
	if err != nil {
		return NewNullVal(), nil
	}
	return env.assignVar(assigne.GetSymbol(), value)
}

func evalAssignmentObjMember(assignment ast.IAssignmentExpr, env *Environment) (IRuntimeVal, error) {
	if assignment.GetAssigne().GetKind() != ast.MemberExprNode {
		return NewNullVal(), fmt.Errorf("evalAssignmentObjMember: Invalid LHS inside assignment expr %s", assignment.GetAssigne())
	}

	var i interface{} = assignment.GetAssigne()
	memberExpr, _ := i.(ast.IMemberExpr)

	if memberExpr.GetObject().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("evalMemberExpr: Object - Node kind '%s' not supported", memberExpr.GetObject().GetKind())
	}

	if memberExpr.GetProperty().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("evalMemberExpr: Property - Node kind '%s' not supported", memberExpr.GetProperty().GetKind())
	}

	i = memberExpr.GetObject()
	identifier, _ := i.(ast.IIdentifier)
	obj, err := env.lookupVar(identifier.GetSymbol())
	if err != nil {
		return NewNullVal(), err
	}
	if obj.GetType() != ObjValueType {
		return NewNullVal(), fmt.Errorf("evalMemberExpr: variable '%s' is not of type 'object'", identifier.GetSymbol())
	}

	i = obj
	objVal, _ := i.(IObjVal)

	i = memberExpr.GetProperty()
	property, _ := i.(ast.IIdentifier)

	value, err := Evaluate(assignment.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}
	objVal.GetProperties()[property.GetSymbol()] = value
	return value, nil
}

func evalObjectExpr(object ast.IObjectLiteral, env *Environment) (IRuntimeVal, error) {
	obj := NewObjVal()

	for _, property := range object.GetProperties() {
		value, err := Evaluate(property.GetValue(), env)
		if err != nil {
			return NewNullVal(), err
		}
		obj.properties[property.GetKey()] = value
	}

	return obj, nil
}

func evalMemberExpr(memberExpr ast.IMemberExpr, env *Environment) (IRuntimeVal, error) {
	if memberExpr.GetObject().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("evalMemberExpr: Object - Node kind '%s' not supported", memberExpr.GetObject().GetKind())
	}

	if memberExpr.GetProperty().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("evalMemberExpr: Property - Node kind '%s' not supported", memberExpr.GetProperty().GetKind())
	}

	var i interface{} = memberExpr.GetObject()
	identifier, _ := i.(ast.IIdentifier)
	obj, err := env.lookupVar(identifier.GetSymbol())
	if err != nil {
		return NewNullVal(), err
	}
	if obj.GetType() != ObjValueType {
		return NewNullVal(), fmt.Errorf("evalMemberExpr: variable '%s' is not of type 'object'", identifier.GetSymbol())
	}

	i = obj
	objVal, _ := i.(IObjVal)

	i = memberExpr.GetProperty()
	property, _ := i.(ast.IIdentifier)

	return objVal.GetProperties()[property.GetSymbol()], nil
}

func evalCallExpr(call ast.ICallExpr, env *Environment) (IRuntimeVal, error) {
	// TODO add helpers? https://zetcode.com/golang/filter-map/
	var args []IRuntimeVal
	for _, arg := range call.GetArgs() {
		evalArg, err := Evaluate(arg, env)
		if err != nil {
			return NewNullVal(), err
		}
		args = append(args, evalArg)
	}

	if call.GetCaller().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("Function caller has to be an identifier. Got: %s", call.GetCaller())
	}
	var j interface{} = call.GetCaller()
	identifier, _ := j.(ast.IIdentifier)
	caller, err := env.lookupFunc(identifier.GetSymbol())
	if err != nil {
		return NewNullVal(), err
	}

	switch caller.GetType() {
	case NativeFnType:
		var i interface{} = caller
		fn, _ := i.(INativeFunc)
		return fn.GetCall()(args, env), nil

	case FunctionValueType:
		var i interface{} = caller
		fn, _ := i.(IFunctionVal)
		scope := NewEnvironment(fn.GetDeclarationEnv())

		// Create variables for the parameters list
		for i := 0; i < len(fn.GetParams()); i++ {
			// TODO Check the bounds here. Verify arity of function.
			// Which means: len(fn.GetParams()) == len(args)
			// TODO var type - Get from function declaration and validate type against given type
			scope.declareVar(fn.GetParams()[i].GetName(), args[i], false, fn.GetParams()[i].GetParamType())
		}

		var result IRuntimeVal
		result = NewNullVal()
		// Evaluate the function body line by line
		for _, stmt := range fn.GetBody() {
			var err error
			result, err = Evaluate(stmt, scope)
			if err != nil {
				return NewNullVal(), err
			}
		}
		return result, nil

	default:
		return NewNullVal(), fmt.Errorf("Cannot call value that is not a function: %s", caller)
	}
}
