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
		lhs.SetTranspilat(identNodeToBashVar(binOp.GetLeft()))
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
		rhs.SetTranspilat(identNodeToBashVar(binOp.GetRight()))
	case ast.IntLiteralNode, ast.StrLiteralNode:
		rhs.SetTranspilat(rhs.ToString())
	default:
		return NewNullVal(), fmt.Errorf("evalBinaryExpr: right kind '%s' not supported", binOp.GetLeft())
	}

	if lhs.GetType() == IntValueType && rhs.GetType() == IntValueType {
		return evalIntBinaryExpr(runtimetoIntVal(lhs), runtimetoIntVal(rhs), binOp.GetOperator())
	}

	if lhs.GetType() == StrValueType && rhs.GetType() == StrValueType {
		return evalStrBinaryExpr(runtimetoStrVal(lhs), runtimetoStrVal(rhs), binOp.GetOperator())
	}

	return NewNullVal(), fmt.Errorf("evalBinaryExpr: Give types not supported (lhs: %s, rhs: %s)", lhs, rhs)
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

	value, err := transpile(assignment.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}

	varName := identNodeGetSymbol(assignment.GetAssigne())

	writeToFile(varName + "=")

	switch assignment.GetValue().GetKind() {
	case ast.CallExprNode:
		varName, err := getCallerResultVarName(ast.ExprToCallExpr(assignment.GetValue()), env)
		if err != nil {
			return NewNullVal(), err
		}

		varType, err := env.lookupVarType(varName)
		if err != nil {
			return NewNullVal(), err
		}
		switch varType {
		case lexer.StrType:
			writeLnToFile("\"" + varName + "\"")
			value = NewStrVal("")
		case lexer.IntType:
			writeLnToFile(varName)
			value = NewIntVal(1)
		default:
			return NewNullVal(), fmt.Errorf("evalAssignment - CallExpr: Unsupported varType '%s'", varType)
		}
	case ast.BinaryExprNode:
		varType, err := env.lookupVarType(varName)
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
		writeLnToFile(identNodeGetSymbol(assignment.GetValue()))
	default:
		return NewNullVal(), fmt.Errorf("evalAssignment: value kind '%s' not supported", assignment.GetValue().GetKind())
	}

	result, err := env.assignVar(varName, value)
	if err != nil {
		return NewNullVal(), err
	}
	return result, nil
}

func evalAssignmentObjMember(assignment ast.IAssignmentExpr, env *Environment) (IRuntimeVal, error) {
	if assignment.GetAssigne().GetKind() != ast.MemberExprNode {
		return NewNullVal(), fmt.Errorf("evalAssignmentObjMember: Invalid LHS inside assignment expr: %s", assignment.GetAssigne())
	}

	memberExpr := ast.ExprToMemberExpr(assignment.GetAssigne())

	if memberExpr.GetObject().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("evalAssignmentObjMember: Object - Node kind '%s' not supported", memberExpr.GetObject().GetKind())
	}

	if memberExpr.GetProperty().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("evalAssignmentObjMember: Property - Node kind '%s' not supported", memberExpr.GetProperty().GetKind())
	}

	objName := identNodeGetSymbol(memberExpr.GetObject())
	obj, err := env.lookupVar(objName)
	if err != nil {
		return NewNullVal(), err
	}
	if obj.GetType() != ObjValueType {
		return NewNullVal(), fmt.Errorf("evalAssignmentObjMember: variable '%s' is not of type 'object'", objName)
	}

	propName := identNodeGetSymbol(memberExpr.GetProperty())

	value, err := transpile(assignment.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}

	writeToFile(objName + "[\"" + propName + "\"]=")

	switch assignment.GetValue().GetKind() {
	case ast.IntLiteralNode:
		writeLnToFile(value.ToString())
	default:
		return NewNullVal(), fmt.Errorf("evalAssignmentObjMember: value kind '%s' not supported", assignment.GetValue().GetKind())
	}

	runtimetoObjVal(obj).GetProperties()[propName] = value
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

	objName := identNodeGetSymbol(memberExpr.GetObject())
	obj, err := env.lookupVar(objName)
	if err != nil {
		return NewNullVal(), err
	}
	if obj.GetType() != ObjValueType {
		return NewNullVal(), fmt.Errorf("evalMemberExpr: variable '%s' is not of type 'object'", objName)
	}

	propName := identNodeGetSymbol(memberExpr.GetProperty())

	propTranspilat := ""
	switch memberExpr.GetProperty().GetKind() {
	case ast.IdentifierNode:
		propTranspilat = "\"" + propName + "\""
	default:
		return NewNullVal(), fmt.Errorf("evalMemberExpr: property kind '%s' not supported", memberExpr.GetProperty().GetKind())
	}

	result := runtimetoObjVal(obj).GetProperties()[propName]
	result.SetTranspilat("${" + objName + "[" + propTranspilat + "]}")
	return result, nil
}

func getCallerResultVarName(call ast.ICallExpr, env *Environment) (string, error) {
	if call.GetCaller().GetKind() != ast.IdentifierNode {
		return "", fmt.Errorf("getCallerResultVarName: Function caller has to be an identifier. Got: %s", call.GetCaller())
	}

	funcName := identNodeGetSymbol(call.GetCaller())
	caller, err := env.lookupFunc(funcName)
	if err != nil {
		return "", err
	}
	if caller.GetType() == FunctionValueType {
		return "$?", nil
	} else if caller.GetType() == NativeFnType {
		// TODO Determine based on return type, if that is implemented
		switch funcName {
		case "print", "printLn":
			return "", fmt.Errorf("'" + funcName + "' has no return value")
		case "input":
			return "${tmpStr}", nil
		default:
			return "", fmt.Errorf("getCallerResultVarName: Return type for func '%s' is unknown", funcName)
		}
	} else {
		return "", fmt.Errorf("getCallerResultVarName: Function type '%s' not supported", caller.GetType())

	}
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

	caller, err := env.lookupFunc(identNodeGetSymbol(call.GetCaller()))
	if err != nil {
		return NewNullVal(), err
	}

	switch caller.GetType() {
	case NativeFnType:
		result, err := runtimetoNativeFunc(caller).GetCall()(call.GetArgs(), env)
		if err != nil {
			return NewNullVal(), err
		}
		writeToFile(result.GetTranspilat())
		return result, nil

	case FunctionValueType:
		fn := runtimetoFuncVal(caller)
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
				switch param.GetParamType() {
				case lexer.IntType:
					writeToFile(" " + identNodeToBashVar(call.GetArgs()[i]))
				case lexer.StrType:
					writeToFile(" \"" + identNodeToBashVar(call.GetArgs()[i]) + "\"")
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

func evalReturnExpr(returnExpr ast.IReturnExpr, env *Environment) (IRuntimeVal, error) {
	if !funcContext {
		return NewNullVal(), fmt.Errorf("Return is only allowed inside of a function")
	}

	value, err := transpile(returnExpr.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}

	writeToFile("return ")
	switch returnExpr.GetValue().GetKind() {
	case ast.BinaryExprNode:
		switch value.GetType() {
		case StrValueType:
			writeLnToFile("\"" + value.GetTranspilat() + "\"")
		case IntValueType:
			writeLnToFile(value.GetTranspilat())
		default:
			return NewNullVal(), fmt.Errorf("evalReturnExpr - BinaryExpr: Unsupported varType '%s'", value.GetType())
		}
	case ast.IntLiteralNode:
		writeLnToFile(value.ToString())
	case ast.IdentifierNode:
		writeLnToFile(identNodeToBashVar(returnExpr.GetValue()))
	default:
		return NewNullVal(), fmt.Errorf("evalReturnExpr: Unsupported value kind '%s'", returnExpr.GetValue().GetKind())
	}
	return value, nil
}
