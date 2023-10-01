package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"ScriLa/cmd/scrila/lexer"
	"fmt"
)

func evalIdentifier(identifier ast.IIdentifier, env *Environment) (IRuntimeVal, error) {
	return env.lookupVar(identifier.GetSymbol())
}

func evalBinaryExpr(binOp ast.IBinaryExpr, env *Environment) (IRuntimeVal, error) {
	lhs, lhsError := transpile(binOp.GetLeft(), env)
	if lhsError != nil {
		return NewNullVal(), lhsError
	}
	switch binOp.GetLeft().GetKind() {
	case ast.BinaryExprNode:
		// Do nothing
	case ast.IdentifierNode:
		var i interface{} = binOp.GetLeft()
		identifier, _ := i.(ast.IIdentifier)
		lhs.SetTranspilat("${" + identifier.GetSymbol() + "}")
	case ast.IntLiteralNode, ast.StrLiteralNode:
		lhs.SetTranspilat(lhs.ToString())
	default:
		return NewNullVal(), fmt.Errorf("evalBinaryExpr: left kind '%s' not supported", binOp.GetLeft())
	}

	rhs, rhsError := transpile(binOp.GetRight(), env)
	if rhsError != nil {
		return NewNullVal(), rhsError
	}
	switch binOp.GetRight().GetKind() {
	case ast.BinaryExprNode:
		// Do nothing
	case ast.IdentifierNode:
		var i interface{} = binOp.GetRight()
		identifier, _ := i.(ast.IIdentifier)
		rhs.SetTranspilat("${" + identifier.GetSymbol() + "}")
	case ast.IntLiteralNode, ast.StrLiteralNode:
		rhs.SetTranspilat(rhs.ToString())
	default:
		return NewNullVal(), fmt.Errorf("evalBinaryExpr: right kind '%s' not supported", binOp.GetLeft())
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

	// TODO Return Error One or both are NULL, or another not yet supported type
	return NewNullVal(), nil
}

func evalIntBinaryExpr(lhs IIntVal, rhs IIntVal, operator string) (IIntVal, error) {
	var result int64
	transpilat := "$(("
	switch operator {
	case "+":
		transpilat += lhs.GetTranspilat() + " + " + rhs.GetTranspilat()
		result = lhs.GetValue() + rhs.GetValue()
	case "-":
		transpilat += lhs.GetTranspilat() + " - " + rhs.GetTranspilat()
		result = lhs.GetValue() - rhs.GetValue()
	case "*":
		transpilat += lhs.GetTranspilat() + " * " + rhs.GetTranspilat()
		result = lhs.GetValue() * rhs.GetValue()
	case "/":
		transpilat += lhs.GetTranspilat() + " / " + rhs.GetTranspilat()
		// TODO Division by zero
		result = lhs.GetValue() / rhs.GetValue()
	default:
		return NewIntVal(0), fmt.Errorf("evalIntBinaryExpr: Unsupported binary operator: %s", operator)
	}
	transpilat += "))"

	intVal := NewIntVal(result)
	intVal.SetTranspilat(transpilat)
	return intVal, nil
}

func evalStrBinaryExpr(lhs IStrVal, rhs IStrVal, operator string) (IStrVal, error) {
	switch operator {
	case "+":
		transpilat := lhs.GetTranspilat() + rhs.GetTranspilat()
		result := lhs.GetValue() + rhs.GetValue()
		strVal := NewStrVal(result)
		strVal.SetTranspilat(transpilat)
		return strVal, nil
	default:
		return NewStrVal(""), fmt.Errorf("evalStrBinaryExpr: Unsupported binary operator: %s", operator)
	}
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
	writeToFile(assigne.GetSymbol() + "=")

	value, err := transpile(assignment.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}

	switch assignment.GetValue().GetKind() {
	case ast.BinaryExprNode:
		varType, err := env.lookupVarType(assigne.GetSymbol())
		if err != nil {
			return NewNullVal(), err
		}
		switch varType {
		case lexer.StrType:
			writeLnToFile("\"" + value.GetTranspilat() + "\"")
		case lexer.IntType:
			writeLnToFile(value.GetTranspilat())
		default:
			return NewNullVal(), fmt.Errorf("evalAssignment - BinaryExpr: Unsupported varType '%s'", varType)
		}
	case ast.IntLiteralNode:
		writeLnToFile(value.ToString())
	case ast.IdentifierNode:
		var i interface{} = assignment.GetValue()
		identifier, _ := i.(ast.IIdentifier)
		writeLnToFile("$" + identifier.GetSymbol())
	default:
		return NewNullVal(), fmt.Errorf("evalAssignment: value kind '%s' not supported", assignment.GetValue().GetKind())
	}

	result, err := env.assignVar(assigne.GetSymbol(), value)
	if err != nil {
		return NewNullVal(), err
	}
	return result, nil
}

func evalAssignmentObjMember(assignment ast.IAssignmentExpr, env *Environment) (IRuntimeVal, error) {
	if assignment.GetAssigne().GetKind() != ast.MemberExprNode {
		return NewNullVal(), fmt.Errorf("evalAssignmentObjMember: Invalid LHS inside assignment expr: %s", assignment.GetAssigne())
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

	value, err := transpile(assignment.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}
	objVal.GetProperties()[property.GetSymbol()] = value
	return value, nil
}

func evalObjectExpr(object ast.IObjectLiteral, env *Environment) (IRuntimeVal, error) {
	obj := NewObjVal()

	for _, property := range object.GetProperties() {
		value, err := transpile(property.GetValue(), env)
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
		evalArg, err := transpile(arg, env)
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
		return fn.GetCall()(call.GetArgs(), env)

	case FunctionValueType:
		var i interface{} = caller
		fn, _ := i.(IFunctionVal)

		writeToFile(fn.GetName())
		for i, param := range fn.GetParams() {
			// TODO var type - Get from function declaration and validate type against given type
			// args[i] param.GetParamType()
			switch call.GetArgs()[i].GetKind() {
			case ast.IntLiteralNode:
				writeToFile(" " + args[i].ToString())
			case ast.StrLiteralNode:
				writeToFile(" \"" + args[i].ToString() + "\"")
			case ast.IdentifierNode:
				var iIdent interface{} = call.GetArgs()[i]
				ident, _ := iIdent.(ast.IIdentifier)
				switch param.GetParamType() {
				case lexer.IntType:
					writeToFile(" $" + ident.GetSymbol())
				case lexer.StrType:
					writeToFile(" \"$" + ident.GetSymbol() + "\"")
				default:
					return NewNullVal(), fmt.Errorf("evalCallExpr - Identifier: Param type '%s' not supported", param.GetParamType())
				}
			default:
				return NewNullVal(), fmt.Errorf("evalCallExpr: Arg type '%s' not supported", call.GetArgs()[i].GetKind())
			}
		}

		writeLnToFile("")
		return NewNullVal(), nil

	default:
		return NewNullVal(), fmt.Errorf("Cannot call value that is not a function: %s", caller)
	}
}
