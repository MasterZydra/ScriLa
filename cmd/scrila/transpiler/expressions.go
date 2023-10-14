package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"ScriLa/cmd/scrila/lexer"
	"fmt"

	"golang.org/x/exp/slices"
)

func (self *Transpiler) evalIdentifier(identifier ast.IIdentifier, env *Environment) (IRuntimeVal, error) {
	return env.lookupVar(identifier.GetSymbol())
}

func (self *Transpiler) evalBinaryExpr(binOp ast.IBinaryExpr, env *Environment) (IRuntimeVal, error) {
	lhs, lhsError := self.transpile(binOp.GetLeft(), env)
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
		return NewNullVal(), fmt.Errorf("%s: Left side of binary expression with unsupported type '%s'", self.getPos(binOp.GetLeft()), binOp.GetLeft().GetKind())
	}

	rhs, rhsError := self.transpile(binOp.GetRight(), env)
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
		return NewNullVal(), fmt.Errorf("%s: Right side of binary expression with unsupported type '%s'", self.getPos(binOp.GetRight()), binOp.GetRight().GetKind())
	}

	if lhs.GetType() == IntValueType && rhs.GetType() == IntValueType {
		result, err := self.evalIntBinaryExpr(runtimeToIntVal(lhs), runtimeToIntVal(rhs), binOp.GetOperator())
		if err != nil {
			return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(binOp), err)
		}
		return result, nil
	}

	if lhs.GetType() == StrValueType && rhs.GetType() == StrValueType {
		result, err := self.evalStrBinaryExpr(runtimeToStrVal(lhs), runtimeToStrVal(rhs), binOp.GetOperator())
		if err != nil {
			return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(binOp), err)
		}
		return result, nil
	}

	if lhs.GetType() == BoolValueType && rhs.GetType() == BoolValueType {
		result, err := self.evalBoolBinaryExpr(runtimeToBoolVal(lhs), runtimeToBoolVal(rhs), binOp.GetOperator())
		if err != nil {
			return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(binOp), err)
		}
		return result, nil
	}

	return NewNullVal(), fmt.Errorf("%s: No support for binary expressions of type '%s' and '%s'", self.getPos(binOp), lhs.GetType(), rhs.GetType())
}

func (self *Transpiler) evalBoolBinaryExpr(lhs IBoolVal, rhs IBoolVal, operator string) (IBoolVal, error) {
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

func (self *Transpiler) evalIntBinaryExpr(lhs IIntVal, rhs IIntVal, operator string) (IIntVal, error) {
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

func (self *Transpiler) evalStrBinaryExpr(lhs IStrVal, rhs IStrVal, operator string) (IStrVal, error) {
	switch operator {
	case "+":
		strVal := NewStrVal(lhs.GetValue() + rhs.GetValue())
		strVal.SetTranspilat(lhs.GetTranspilat() + rhs.GetTranspilat())
		return strVal, nil
	default:
		return NewStrVal(""), fmt.Errorf("Binary string expression with unsupported operator '%s'", operator)
	}
}

func (self *Transpiler) evalAssignment(assignment ast.IAssignmentExpr, env *Environment) (IRuntimeVal, error) {
	if assignment.GetAssigne().GetKind() == ast.MemberExprNode {
		return self.evalAssignmentObjMember(assignment, env)
	}

	if assignment.GetAssigne().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("%s: Left side of an assignment must be a variable. Got '%s'", self.getPos(assignment.GetAssigne()), assignment.GetAssigne().GetKind())
	}

	value, err := self.transpile(assignment.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}

	varName := identNodeGetSymbol(assignment.GetAssigne())
	varType, err := env.lookupVarType(varName)
	if err != nil {
		return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(assignment.GetAssigne()), err)
	}

	self.writeTranspilat(varName + "=")

	switch assignment.GetValue().GetKind() {
	case ast.CallExprNode:
		returnVarName, err := self.getCallerResultVarName(ast.ExprToCallExpr(assignment.GetValue()), env)
		if err != nil {
			return NewNullVal(), err
		}
		switch varType {
		case lexer.StrType:
			self.writeLnTranspilat(strToBashStr(returnVarName))
			value = NewStrVal("")
		case lexer.IntType:
			self.writeLnTranspilat(returnVarName)
			value = NewIntVal(1)
		case lexer.BoolType:
			self.writeLnTranspilat(strToBashStr(returnVarName))
			value = NewBoolVal(true)
		default:
			return NewNullVal(), fmt.Errorf("%s: Assigning return values is not implemented for variables of type '%s'", self.getPos(assignment), varType)
		}
	case ast.BinaryExprNode:
		varType, err := env.lookupVarType(varName)
		if err != nil {
			return NewNullVal(), err
		}
		switch varType {
		case lexer.StrType:
			self.writeLnTranspilat(strToBashStr(value.GetTranspilat()))
		case lexer.IntType:
			self.writeLnTranspilat(value.GetTranspilat())
		default:
			return NewNullVal(), fmt.Errorf("%s: Assigning binary expressions is not implemented for variables of type '%s'", self.getPos(assignment), varType)
		}
	case ast.IntLiteralNode:
		if varType != lexer.IntType {
			return NewNullVal(), fmt.Errorf("%s: Cannot assign a value of type '%s' to a var of type '%s'", self.getPos(assignment.GetValue()), lexer.IntType, varType)
		}
		self.writeLnTranspilat(value.ToString())
	case ast.StrLiteralNode:
		if varType != lexer.StrType {
			return NewNullVal(), fmt.Errorf("%s: Cannot assign a value of type '%s' to a var of type '%s'", self.getPos(assignment.GetValue()), lexer.StrType, varType)
		}
		self.writeLnTranspilat(strToBashStr(value.ToString()))
	case ast.IdentifierNode:
		symbol := identNodeGetSymbol(assignment.GetValue())
		if symbol == "null" || ast.IdentIsBool(ast.ExprToIdent(assignment.GetValue())) {
			self.writeLnTranspilat(strToBashStr(symbol))
		} else if slices.Contains(reservedIdentifiers, symbol) {
			self.writeLnTranspilat(symbol)
		} else {
			valueVarType, err := env.lookupVarType(identNodeGetSymbol(assignment.GetValue()))
			if err != nil {
				return NewNullVal(), err
			}
			if valueVarType != varType {
				return NewNullVal(), fmt.Errorf("%s: Cannot assign a value of type '%s' to a var of type '%s'", self.getPos(assignment.GetValue()), valueVarType, varType)
			}
			switch varType {
			case lexer.StrType:
				self.writeLnTranspilat(strToBashStr(identNodeToBashVar(assignment.GetValue())))
			case lexer.IntType:
				self.writeLnTranspilat(identNodeToBashVar(assignment.GetValue()))
			default:
				return NewNullVal(), fmt.Errorf("%s: Assigning variables is not implemented for variables of type '%s'", self.getPos(assignment), varType)
			}
		}
	default:
		return NewNullVal(), fmt.Errorf("%s: Assigning variables is not implemented for variables of type '%s'", self.getPos(assignment), assignment.GetKind())
	}

	result, err := env.assignVar(varName, value)
	if err != nil {
		return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(assignment), err)
	}
	return result, nil
}

func (self *Transpiler) evalAssignmentObjMember(assignment ast.IAssignmentExpr, env *Environment) (IRuntimeVal, error) {
	if assignment.GetAssigne().GetKind() != ast.MemberExprNode {
		return NewNullVal(), fmt.Errorf("%s: Left side of object member assignment is invalid type '%s'", self.getPos(assignment.GetAssigne()), assignment.GetAssigne().GetKind())
	}

	memberExpr := ast.ExprToMemberExpr(assignment.GetAssigne())

	if memberExpr.GetObject().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("%s: Object name is not the right type. Got '%s'", self.getPos(memberExpr.GetObject()), memberExpr.GetObject().GetKind())
	}

	if memberExpr.GetProperty().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("%s: Object property name is not the right type. Got '%s'", self.getPos(memberExpr.GetProperty()), memberExpr.GetProperty().GetKind())
	}

	objName := identNodeGetSymbol(memberExpr.GetObject())
	obj, err := env.lookupVar(objName)
	if err != nil {
		return NewNullVal(), err
	}
	if obj.GetType() != ObjValueType {
		return NewNullVal(), fmt.Errorf("%s: Variable '%s' is not of type 'object'", self.getPos(memberExpr.GetObject()), objName)
	}

	propName := identNodeGetSymbol(memberExpr.GetProperty())

	value, err := self.transpile(assignment.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}

	self.writeTranspilat(objName + "[" + strToBashStr(propName) + "]=")

	switch assignment.GetValue().GetKind() {
	case ast.IntLiteralNode:
		self.writeLnTranspilat(value.ToString())
	default:
		return NewNullVal(), fmt.Errorf("%s: Object member value '%s' is not supported", self.getPos(assignment.GetValue()), assignment.GetValue().GetKind())
	}

	runtimeToObjVal(obj).GetProperties()[propName] = value
	return value, nil
}

func (self *Transpiler) evalObjectExpr(object ast.IObjectLiteral, env *Environment) (IRuntimeVal, error) {
	obj := NewObjVal()

	for _, property := range object.GetProperties() {
		value, err := self.transpile(property.GetValue(), env)
		if err != nil {
			return NewNullVal(), err
		}
		obj.properties[property.GetKey()] = value
	}

	return obj, nil
}

func (self *Transpiler) evalMemberExpr(memberExpr ast.IMemberExpr, env *Environment) (IRuntimeVal, error) {
	if memberExpr.GetObject().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("%s: Object name is not the right type. Got '%s'", self.getPos(memberExpr.GetObject()), memberExpr.GetObject().GetKind())
	}

	objName := identNodeGetSymbol(memberExpr.GetObject())
	obj, err := env.lookupVar(objName)
	if err != nil {
		return NewNullVal(), err
	}
	if obj.GetType() != ObjValueType {
		return NewNullVal(), fmt.Errorf("%s: Variable '%s' is not of type 'object'", self.getPos(memberExpr.GetObject()), objName)
	}

	if memberExpr.GetProperty().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("%s: Object property name is not the right type. Got '%s'", self.getPos(memberExpr.GetProperty()), memberExpr.GetProperty().GetKind())
	}

	propName := identNodeGetSymbol(memberExpr.GetProperty())
	propTranspilat := strToBashStr(propName)

	result := runtimeToObjVal(obj).GetProperties()[propName]
	result.SetTranspilat("${" + objName + "[" + propTranspilat + "]}")
	return result, nil
}

func (self *Transpiler) getCallerResultVarName(call ast.ICallExpr, env *Environment) (string, error) {
	if call.GetCaller().GetKind() != ast.IdentifierNode {
		return "", fmt.Errorf("%s: Function name must be an identifier. Got: '%s'", self.getPos(call.GetCaller()), call.GetCaller().GetKind())
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
		case lexer.BoolType:
			return "${tmpBool}", nil
		case lexer.VoidType:
			return "", fmt.Errorf("%s: Func '%s' does not have a return value", self.getPos(call.GetCaller()), funcName)
		default:
			return "", fmt.Errorf("%s: Function return type '%s' is not supported", self.getPos(call.GetCaller()), returnType)
		}
	} else if caller.GetType() == NativeFnType {
		// TODO Determine based on return type, if that is implemented
		switch funcName {
		case "print", "printLn", "sleep":
			return "", fmt.Errorf("%s: Function '%s' has no return value", self.getPos(call.GetCaller()), funcName)
		case "input":
			return "${tmpStr}", nil
		case "isInt":
			return "${tmpBool}", nil
		default:
			return "", fmt.Errorf("%s: Return type for native func '%s' is unknown", self.getPos(call), funcName)
		}
	} else {
		return "", fmt.Errorf("%s: Cannot call value that is not a function: %s", self.getPos(call), caller.GetType())
	}
}

func (self *Transpiler) evalCallExpr(call ast.ICallExpr, env *Environment) (IRuntimeVal, error) {
	// TODO add helpers? https://zetcode.com/golang/filter-map/
	var args []IRuntimeVal
	for _, arg := range call.GetArgs() {
		evalArg, err := self.transpile(arg, env)
		if err != nil {
			return NewNullVal(), err
		}
		args = append(args, evalArg)
	}

	if call.GetCaller().GetKind() != ast.IdentifierNode {
		return NewNullVal(), fmt.Errorf("%s: Function name must be an identifier. Got: '%s'", self.getPos(call.GetCaller()), call.GetCaller().GetKind())
	}

	caller, err := env.lookupFunc(identNodeGetSymbol(call.GetCaller()))
	if err != nil {
		return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(call), err)
	}

	switch caller.GetType() {
	case NativeFnType:
		result, err := runtimeToNativeFunc(caller).GetCall()(call.GetArgs(), env)
		if err != nil {
			return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(call), err)
		}
		self.writeTranspilat(result.GetTranspilat())
		return result, nil

	case FunctionValueType:
		fn := runtimeToFuncVal(caller)

		if len(fn.GetParams()) != len(args) {
			return NewNullVal(), fmt.Errorf("%s: %s(): The amount of passed parameters does not match with the function declaration. Expected: %d, Got: %d", self.getPos(call), fn.GetName(), len(fn.GetParams()), len(args))
		}
		self.writeTranspilat(fn.GetName())
		for i, param := range fn.GetParams() {
			if !doTypesMatch(param.GetParamType(), args[i].GetType()) {
				return NewNullVal(), fmt.Errorf("%s: %s(): Parameter '%s' type does not match. Expected: %s, Got: %s", self.getPos(call), fn.GetName(), param.GetName(), param.GetParamType(), args[i].GetType())
			}
			switch call.GetArgs()[i].GetKind() {
			case ast.IntLiteralNode:
				self.writeTranspilat(" " + args[i].ToString())
			case ast.StrLiteralNode:
				self.writeTranspilat(" " + strToBashStr(args[i].ToString()))
			case ast.IdentifierNode:
				switch param.GetParamType() {
				case lexer.IntType:
					self.writeTranspilat(" " + identNodeToBashVar(call.GetArgs()[i]))
				case lexer.StrType:
					self.writeTranspilat(" " + strToBashStr(identNodeToBashVar(call.GetArgs()[i])))
				default:
					return NewNullVal(), fmt.Errorf("%s: %s(): Parameter of variable type '%s' is not supported", self.getPos(call), fn.GetName(), param.GetParamType())
				}
			default:
				return NewNullVal(), fmt.Errorf("%s: %s(): Parameter type '%s' is not supported", self.getPos(call), fn.GetName(), call.GetArgs()[i].GetKind())
			}
		}

		self.writeLnTranspilat("")
		return NewNullVal(), nil

	default:
		return NewNullVal(), fmt.Errorf("%s: Cannot call value that is not a function: %s", self.getPos(call), caller)
	}
}

func (self *Transpiler) evalReturnExpr(returnExpr ast.IReturnExpr, env *Environment) (IRuntimeVal, error) {
	if !self.funcContext || self.currentFunc == nil {
		return NewNullVal(), fmt.Errorf("%s: Return is only allowed inside a function", self.getPos(returnExpr))
	}

	value, err := self.transpile(returnExpr.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}

	switch self.currentFunc.GetReturnType() {
	case lexer.IntType:
		self.writeTranspilat("tmpInt=")
	case lexer.StrType:
		self.writeTranspilat("tmpStr=")
	case lexer.BoolType:
		self.writeTranspilat("tmpBool=")
	default:
		return NewNullVal(), fmt.Errorf("%s: Return type '%s' is not supported", self.getPos(returnExpr), self.currentFunc.GetReturnType())
	}

	switch returnExpr.GetValue().GetKind() {
	case ast.BinaryExprNode:
		switch value.GetType() {
		case StrValueType:
			self.writeLnTranspilat(strToBashStr(value.GetTranspilat()))
		case IntValueType:
			self.writeLnTranspilat(value.GetTranspilat())
		default:
			return NewNullVal(), fmt.Errorf("%s: Returning binary expression of type '%s' is not supported", self.getPos(returnExpr), value.GetType())
		}
	case ast.IntLiteralNode:
		self.writeLnTranspilat(value.ToString())
	case ast.StrLiteralNode:
		self.writeLnTranspilat(strToBashStr(value.ToString()))
	case ast.IdentifierNode:
		self.writeLnTranspilat(identNodeToBashVar(returnExpr.GetValue()))
	default:
		return NewNullVal(), fmt.Errorf("%s: Return type '%s' is not supported", self.getPos(returnExpr), returnExpr.GetValue().GetKind())
	}
	self.writeLnTranspilat("\treturn")
	return value, nil
}
