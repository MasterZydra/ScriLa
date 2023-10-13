package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"ScriLa/cmd/scrila/lexer"
	"fmt"

	"golang.org/x/exp/slices"
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
		if ast.IdentIsBool(ast.ExprToIdent(binOp.GetLeft())) {
			lhs.SetTranspilat(boolIdentToBashComparison(ast.ExprToIdent(binOp.GetLeft())))
		} else {
			lhs.SetTranspilat(identNodeToBashVar(binOp.GetLeft()))
		}
	case ast.IntLiteralNode, ast.StrLiteralNode:
		lhs.SetTranspilat(lhs.ToString())
	default:
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Left side of binary expression with unsupported type '%s'", fileName, binOp.GetLeft().GetLn(), binOp.GetLeft().GetCol(), binOp.GetLeft().GetKind())
	}

	rhs, rhsError := transpile(binOp.GetRight(), env)
	if rhsError != nil {
		return NewNullVal(), rhsError
	}
	switch binOp.GetRight().GetKind() {
	case ast.BinaryExprNode:
		// Do nothing
	case ast.IdentifierNode:
		if ast.IdentIsBool(ast.ExprToIdent(binOp.GetRight())) {
			rhs.SetTranspilat(boolIdentToBashComparison(ast.ExprToIdent(binOp.GetRight())))
		} else {
			rhs.SetTranspilat(identNodeToBashVar(binOp.GetRight()))
		}
	case ast.IntLiteralNode, ast.StrLiteralNode:
		rhs.SetTranspilat(rhs.ToString())
	default:
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Right side of binary expression with unsupported type '%s'", fileName, binOp.GetRight().GetLn(), binOp.GetRight().GetCol(), binOp.GetRight().GetKind())
	}

	if lhs.GetType() == IntValueType && rhs.GetType() == IntValueType {
		result, err := evalIntBinaryExpr(runtimeToIntVal(lhs), runtimeToIntVal(rhs), binOp.GetOperator())
		if err != nil {
			return NewNullVal(), fmt.Errorf("%s:%d:%d: %s", fileName, binOp.GetLn(), binOp.GetCol(), err)
		}
		return result, nil
	}

	if lhs.GetType() == StrValueType && rhs.GetType() == StrValueType {
		result, err := evalStrBinaryExpr(runtimeToStrVal(lhs), runtimeToStrVal(rhs), binOp.GetOperator())
		if err != nil {
			return NewNullVal(), fmt.Errorf("%s:%d:%d: %s", fileName, binOp.GetLn(), binOp.GetCol(), err)
		}
		return result, nil
	}

	if lhs.GetType() == BoolValueType && rhs.GetType() == BoolValueType {
		result, err := evalBoolBinaryExpr(runtimeToBoolVal(lhs), runtimeToBoolVal(rhs), binOp.GetOperator())
		if err != nil {
			return NewNullVal(), fmt.Errorf("%s:%d:%d: %s", fileName, binOp.GetLn(), binOp.GetCol(), err)
		}
		return result, nil
	}

	return NewNullVal(), fmt.Errorf("%s:%d:%d: No support for binary expressions of type '%s' and '%s'", fileName, binOp.GetLn(), binOp.GetCol(), lhs.GetType(), rhs.GetType())
}

func evalBoolBinaryExpr(lhs IBoolVal, rhs IBoolVal, operator string) (IBoolVal, error) {
	var result bool
	transpilat := ""
	switch operator {
	case "&&":
		transpilat += lhs.GetTranspilat() + " && " + rhs.GetTranspilat()
		result = lhs.GetValue() && rhs.GetValue()
	case "||":
		transpilat += lhs.GetTranspilat() + " || " + rhs.GetTranspilat()
		result = lhs.GetValue() || rhs.GetValue()
	default:
		return NewBoolVal(false), fmt.Errorf("Binary bool expression with unsupported operator '%s'", operator)
	}

	boolVal := NewBoolVal(result)
	boolVal.SetTranspilat(transpilat)
	return boolVal, nil
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
		return NewIntVal(0), fmt.Errorf("Binary int expression with unsupported operator '%s'", operator)
	}
	transpilat += "))"

	intVal := NewIntVal(result)
	intVal.SetTranspilat(transpilat)
	return intVal, nil
}

func evalStrBinaryExpr(lhs IStrVal, rhs IStrVal, operator string) (IStrVal, error) {
	switch operator {
	case "+":
		strVal := NewStrVal(lhs.GetValue() + rhs.GetValue())
		strVal.SetTranspilat(lhs.GetTranspilat() + rhs.GetTranspilat())
		return strVal, nil
	default:
		return NewStrVal(""), fmt.Errorf("Binary string expression with unsupported operator '%s'", operator)
	}
}

func evalAssignment(assignment ast.IAssignmentExpr, env *Environment) (IRuntimeVal, error) {
	if assignment.GetAssigne().GetKind() == ast.MemberExprNode {
		return evalAssignmentObjMember(assignment, env)
	}

	if assignment.GetAssigne().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Left side of an assignment must be a variable. Got '%s'", fileName, assignment.GetAssigne().GetLn(), assignment.GetAssigne().GetCol(), assignment.GetAssigne().GetKind())
	}

	value, err := transpile(assignment.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}

	varName := identNodeGetSymbol(assignment.GetAssigne())
	varType, err := env.lookupVarType(varName)
	if err != nil {
		return NewNullVal(), fmt.Errorf("%s:%d:%d: %s", fileName, assignment.GetAssigne().GetLn(), assignment.GetAssigne().GetCol(), err)
	}

	writeToFile(varName + "=")

	switch assignment.GetValue().GetKind() {
	case ast.CallExprNode:
		returnVarName, err := getCallerResultVarName(ast.ExprToCallExpr(assignment.GetValue()), env)
		if err != nil {
			return NewNullVal(), err
		}
		switch varType {
		case lexer.StrType:
			writeLnToFile(strToBashStr(returnVarName))
			value = NewStrVal("")
		case lexer.IntType:
			writeLnToFile(returnVarName)
			value = NewIntVal(1)
		default:
			return NewNullVal(), fmt.Errorf("%s:%d:%d: Assigning return values is not implemented for variables of type '%s'", fileName, assignment.GetLn(), assignment.GetCol(), varType)
		}
	case ast.BinaryExprNode:
		varType, err := env.lookupVarType(varName)
		if err != nil {
			return NewNullVal(), err
		}
		switch varType {
		case lexer.StrType:
			writeLnToFile(strToBashStr(value.GetTranspilat()))
		case lexer.IntType:
			writeLnToFile(value.GetTranspilat())
		default:
			return NewNullVal(), fmt.Errorf("%s:%d:%d: Assigning binary expressions is not implemented for variables of type '%s'", fileName, assignment.GetLn(), assignment.GetCol(), varType)
		}
	case ast.IntLiteralNode:
		if varType != lexer.IntType {
			return NewNullVal(), fmt.Errorf("%s:%d:%d: Cannot assign a value of type '%s' to a var of type '%s'", fileName, assignment.GetValue().GetLn(), assignment.GetValue().GetCol(), lexer.IntType, varType)
		}
		writeLnToFile(value.ToString())
	case ast.StrLiteralNode:
		if varType != lexer.StrType {
			return NewNullVal(), fmt.Errorf("%s:%d:%d: Cannot assign a value of type '%s' to a var of type '%s'", fileName, assignment.GetValue().GetLn(), assignment.GetValue().GetCol(), lexer.StrType, varType)
		}
		writeLnToFile(strToBashStr(value.ToString()))
	case ast.IdentifierNode:
		symbol := identNodeGetSymbol(assignment.GetValue())
		if symbol == "null" || ast.IdentIsBool(ast.ExprToIdent(assignment.GetValue())) {
			writeLnToFile(strToBashStr(symbol))
		} else if slices.Contains(reservedIdentifiers, symbol) {
			writeLnToFile(symbol)
		} else {
			valueVarType, err := env.lookupVarType(identNodeGetSymbol(assignment.GetValue()))
			if err != nil {
				return NewNullVal(), err
			}
			if valueVarType != varType {
				return NewNullVal(), fmt.Errorf("%s:%d:%d: Cannot assign a value of type '%s' to a var of type '%s'", fileName, assignment.GetValue().GetLn(), assignment.GetValue().GetCol(), valueVarType, varType)
			}
			switch varType {
			case lexer.StrType:
				writeLnToFile(strToBashStr(identNodeToBashVar(assignment.GetValue())))
			case lexer.IntType:
				writeLnToFile(identNodeToBashVar(assignment.GetValue()))
			default:
				return NewNullVal(), fmt.Errorf("%s:%d:%d: Assigning variables is not implemented for variables of type '%s'", fileName, assignment.GetLn(), assignment.GetCol(), varType)
			}
		}
	default:
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Assigning variables is not implemented for variables of type '%s'", fileName, assignment.GetLn(), assignment.GetCol(), assignment.GetKind())
	}

	result, err := env.assignVar(varName, value)
	if err != nil {
		return NewNullVal(), fmt.Errorf("%s:%d:%d: %s", fileName, assignment.GetLn(), assignment.GetCol(), err)
	}
	return result, nil
}

func evalAssignmentObjMember(assignment ast.IAssignmentExpr, env *Environment) (IRuntimeVal, error) {
	if assignment.GetAssigne().GetKind() != ast.MemberExprNode {
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Left side of object member assignment is invalid type '%s'", fileName, assignment.GetAssigne().GetLn(), assignment.GetAssigne().GetCol(), assignment.GetAssigne().GetKind())
	}

	memberExpr := ast.ExprToMemberExpr(assignment.GetAssigne())

	if memberExpr.GetObject().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Object name is not the right type. Got '%s'", fileName, memberExpr.GetObject().GetLn(), memberExpr.GetObject().GetCol(), memberExpr.GetObject().GetKind())
	}

	if memberExpr.GetProperty().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Object property name is not the right type. Got '%s'", fileName, memberExpr.GetProperty().GetLn(), memberExpr.GetProperty().GetCol(), memberExpr.GetProperty().GetKind())
	}

	objName := identNodeGetSymbol(memberExpr.GetObject())
	obj, err := env.lookupVar(objName)
	if err != nil {
		return NewNullVal(), err
	}
	if obj.GetType() != ObjValueType {
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Variable '%s' is not of type 'object'", fileName, memberExpr.GetObject().GetLn(), memberExpr.GetObject().GetCol(), objName)
	}

	propName := identNodeGetSymbol(memberExpr.GetProperty())

	value, err := transpile(assignment.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}

	writeToFile(objName + "[" + strToBashStr(propName) + "]=")

	switch assignment.GetValue().GetKind() {
	case ast.IntLiteralNode:
		writeLnToFile(value.ToString())
	default:
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Object member value '%s' is not supported", fileName, assignment.GetValue().GetLn(), assignment.GetValue().GetCol(), assignment.GetValue().GetKind())
	}

	runtimeToObjVal(obj).GetProperties()[propName] = value
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
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Object name is not the right type. Got '%s'", fileName, memberExpr.GetObject().GetLn(), memberExpr.GetObject().GetCol(), memberExpr.GetObject().GetKind())
	}

	objName := identNodeGetSymbol(memberExpr.GetObject())
	obj, err := env.lookupVar(objName)
	if err != nil {
		return NewNullVal(), err
	}
	if obj.GetType() != ObjValueType {
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Variable '%s' is not of type 'object'", fileName, memberExpr.GetObject().GetLn(), memberExpr.GetObject().GetCol(), objName)
	}

	if memberExpr.GetProperty().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Object property name is not the right type. Got '%s'", fileName, memberExpr.GetProperty().GetLn(), memberExpr.GetProperty().GetCol(), memberExpr.GetProperty().GetKind())
	}

	propName := identNodeGetSymbol(memberExpr.GetProperty())
	propTranspilat := strToBashStr(propName)

	result := runtimeToObjVal(obj).GetProperties()[propName]
	result.SetTranspilat("${" + objName + "[" + propTranspilat + "]}")
	return result, nil
}

func getCallerResultVarName(call ast.ICallExpr, env *Environment) (string, error) {
	if call.GetCaller().GetKind() != ast.IdentifierNode {
		return "", fmt.Errorf("%s:%d:%d: Function name must be an identifier. Got: '%s'", fileName, call.GetCaller().GetLn(), call.GetCaller().GetCol(), call.GetCaller().GetKind())
	}

	funcName := identNodeGetSymbol(call.GetCaller())
	caller, err := env.lookupFunc(funcName)
	if err != nil {
		return "", err
	}
	if caller.GetType() == FunctionValueType {
		switch returnType := runtimeToFuncVal(caller).GetReturnType(); returnType {
		case lexer.IntType:
			return "${tmpInt}", nil
		case lexer.StrType:
			return "${tmpStr}", nil
		case lexer.VoidType:
			return "", fmt.Errorf("%s:%d:%d: Func '%s' does not have a return value", fileName, call.GetCaller().GetLn(), call.GetCol(), funcName)
		default:
			return "", fmt.Errorf("%s:%d:%d: Function return type '%s' is not supported", fileName, call.GetCaller().GetLn(), call.GetCaller().GetCol(), returnType)
		}
	} else if caller.GetType() == NativeFnType {
		// TODO Determine based on return type, if that is implemented
		switch funcName {
		case "print", "printLn", "sleep":
			return "", fmt.Errorf("%s:%d:%d: Function '%s' has no return value", fileName, call.GetCaller().GetLn(), call.GetCaller().GetCol(), funcName)
		case "input":
			return "${tmpStr}", nil
		default:
			return "", fmt.Errorf("%s:%d:%d: Return type for native func '%s' is unknown", fileName, call.GetLn(), call.GetCol(), funcName)
		}
	} else {
		return "", fmt.Errorf("%s:%d:%d: Cannot call value that is not a function: %s", fileName, call.GetLn(), call.GetCol(), caller.GetType())
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
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Function name must be an identifier. Got: '%s'", fileName, call.GetCaller().GetLn(), call.GetCaller().GetCol(), call.GetCaller().GetKind())
	}

	caller, err := env.lookupFunc(identNodeGetSymbol(call.GetCaller()))
	if err != nil {
		return NewNullVal(), fmt.Errorf("%s:%d:%d: %s", fileName, call.GetLn(), call.GetCol(), err)
	}

	switch caller.GetType() {
	case NativeFnType:
		result, err := runtimeToNativeFunc(caller).GetCall()(call.GetArgs(), env)
		if err != nil {
			return NewNullVal(), fmt.Errorf("%s:%d:%d: %s", fileName, call.GetLn(), call.GetCol(), err)
		}
		writeToFile(result.GetTranspilat())
		return result, nil

	case FunctionValueType:
		fn := runtimeToFuncVal(caller)

		if len(fn.GetParams()) != len(args) {
			return NewNullVal(), fmt.Errorf("%s:%d:%d: %s(): The amount of passed parameters does not match with the function declaration. Expected: %d, Got: %d", fileName, call.GetLn(), call.GetCol(), fn.GetName(), len(fn.GetParams()), len(args))
		}
		writeToFile(fn.GetName())
		for i, param := range fn.GetParams() {
			if !doTypesMatch(param.GetParamType(), args[i].GetType()) {
				return NewNullVal(), fmt.Errorf("%s:%d:%d: %s(): Parameter '%s' type does not match. Expected: %s, Got: %s", fileName, call.GetLn(), call.GetCol(), fn.GetName(), param.GetName(), param.GetParamType(), args[i].GetType())
			}
			switch call.GetArgs()[i].GetKind() {
			case ast.IntLiteralNode:
				writeToFile(" " + args[i].ToString())
			case ast.StrLiteralNode:
				writeToFile(" " + strToBashStr(args[i].ToString()))
			case ast.IdentifierNode:
				switch param.GetParamType() {
				case lexer.IntType:
					writeToFile(" " + identNodeToBashVar(call.GetArgs()[i]))
				case lexer.StrType:
					writeToFile(" " + strToBashStr(identNodeToBashVar(call.GetArgs()[i])))
				default:
					return NewNullVal(), fmt.Errorf("%s:%d:%d: %s(): Parameter of variable type '%s' is not supported", fileName, call.GetLn(), call.GetCol(), fn.GetName(), param.GetParamType())
				}
			default:
				return NewNullVal(), fmt.Errorf("%s:%d:%d: %s(): Parameter type '%s' is not supported", fileName, call.GetLn(), call.GetCol(), fn.GetName(), call.GetArgs()[i].GetKind())
			}
		}

		writeLnToFile("")
		return NewNullVal(), nil

	default:
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Cannot call value that is not a function: %s", fileName, call.GetLn(), call.GetCol(), caller)
	}
}

func evalReturnExpr(returnExpr ast.IReturnExpr, env *Environment) (IRuntimeVal, error) {
	if !funcContext || currentFunc == nil {
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Return is only allowed inside a function", fileName, returnExpr.GetLn(), returnExpr.GetCol())
	}

	value, err := transpile(returnExpr.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}

	switch currentFunc.GetReturnType() {
	case lexer.IntType:
		writeToFile("tmpInt=")
	case lexer.StrType:
		writeToFile("tmpStr=")
	case lexer.BoolType:
		writeToFile("tmpBool=")
	default:
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Return type '%s' is not supported", fileName, returnExpr.GetLn(), returnExpr.GetCol(), currentFunc.GetReturnType())
	}

	switch returnExpr.GetValue().GetKind() {
	case ast.BinaryExprNode:
		switch value.GetType() {
		case StrValueType:
			writeLnToFile(strToBashStr(value.GetTranspilat()))
		case IntValueType:
			writeLnToFile(value.GetTranspilat())
		default:
			return NewNullVal(), fmt.Errorf("%s:%d:%d: Returning binary expression of type '%s' is not supported", fileName, returnExpr.GetLn(), returnExpr.GetCol(), value.GetType())
		}
	case ast.IntLiteralNode:
		writeLnToFile(value.ToString())
	case ast.StrLiteralNode:
		writeLnToFile(strToBashStr(value.ToString()))
	case ast.IdentifierNode:
		writeLnToFile(identNodeToBashVar(returnExpr.GetValue()))
	default:
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Return type '%s' is not supported", fileName, returnExpr.GetLn(), returnExpr.GetCol(), returnExpr.GetValue().GetKind())
	}
	writeLnToFile("\treturn")
	return value, nil
}
