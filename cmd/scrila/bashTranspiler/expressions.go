package bashTranspiler

import (
	"ScriLa/cmd/scrila/bashAst"
	"ScriLa/cmd/scrila/scrilaAst"
	"fmt"

	"golang.org/x/exp/slices"
)

func (self *Transpiler) evalArray(array scrilaAst.IArray, env *Environment) (scrilaAst.IRuntimeVal, error) {
	bashArray := bashAst.NewArray()

	if len(array.GetValues()) == 0 {
		bashArray.SetDataType(bashAst.VoidNode)
	}

	for i, value := range array.GetValues() {
		if i == 0 {
			// Get the data type of the first element and set it as array data type
			_, givenType, err := self.exprIsType(value, scrilaAst.NodeType(scrilaAst.NullValueType), env)
			if err != nil {
				return NewNullVal(), err
			}

			bashDataType, err := scrilaNodeTypeToBashNodeType(givenType)
			if err != nil {
				return NewNullVal(), err
			}
			bashArray.SetDataType(bashDataType)
		}

		// Check if the data types of the elements match with the array data type
		scrilaDataType, err := bashNodeTypeToScrilaNodeType(bashArray.GetDataType())
		if err != nil {
			return NewNullVal(), err
		}
		doMatch, givenType, err := self.exprIsType(value, scrilaDataType, env)
		if !doMatch {
			return NewNullVal(), fmt.Errorf("%s: An array can only keep one data type. Wanted '%s'. Got '%s'", self.getPos(value), scrilaDataType, givenType)
		}

		bashStmt, err := self.exprToBashStmt(value, env)
		if err != nil {
			return NewNullVal(), err
		}
		bashArray.AddValue(bashStmt)
	}

	self.bashStmtStack[array.GetId()] = bashArray

	return NewNullVal(), nil
}

func (self *Transpiler) evalIdentifier(identifier scrilaAst.IIdentifier, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName(identifier.GetSymbol())

	return env.lookupVar(identifier.GetSymbol())
}

func (self *Transpiler) evalBinaryExpr(binOp scrilaAst.IBinaryExpr, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	// TODO Test if error persists with new change -> binary expr in function call param
	if binOp.GetResult() != nil {
		return binOp.GetResult(), nil
	}

	lhs, err := self.transpile(binOp.GetLeft(), env)
	if err != nil {
		return NewNullVal(), err
	}
	bashLhs, err := self.exprToBashStmt(binOp.GetLeft(), env)
	if err != nil {
		return NewNullVal(), err
	}

	rhs, err := self.transpile(binOp.GetRight(), env)
	if err != nil {
		return NewNullVal(), err
	}
	bashRhs, err := self.exprToBashStmt(binOp.GetRight(), env)
	if err != nil {
		return NewNullVal(), err
	}

	var result scrilaAst.IRuntimeVal = nil
	var opType bashAst.NodeType
	isComparison := false

	if scrilaAst.BinExprIsComp(binOp) {
		isComparison = true
		opType, err = self.evalComparisonBinaryExpr(lhs, rhs, binOp.GetOperator())
		if err != nil {
			return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(binOp), err)
		}
		result = NewBoolVal(true)
		binOp.SetResult(result)
	}

	if !isComparison {
		if lhs.GetType() == scrilaAst.IntValueType && rhs.GetType() == scrilaAst.IntValueType {
			opType = bashAst.IntLiteralNode
			if !slices.Contains([]string{"+", "-", "*", "/"}, binOp.GetOperator()) {
				return NewNullVal(), fmt.Errorf("%s: Binary int expression with unsupported operator '%s'", self.getPos(binOp), binOp.GetOperator())
			}
			result = NewIntVal(1)
			binOp.SetResult(result)
		}

		if lhs.GetType() == scrilaAst.StrValueType && rhs.GetType() == scrilaAst.StrValueType {
			opType = bashAst.StrLiteralNode
			if binOp.GetOperator() != "+" {
				return NewNullVal(), fmt.Errorf("%s: Binary string expression with unsupported operator '%s'", self.getPos(binOp), binOp.GetOperator())
			}
			result = NewStrVal("str")
			binOp.SetResult(result)
		}

		if lhs.GetType() == scrilaAst.BoolValueType && rhs.GetType() == scrilaAst.BoolValueType {
			opType = bashAst.BoolLiteralNode
			if !slices.Contains([]string{"&&", "||"}, binOp.GetOperator()) {
				return NewBoolVal(false), fmt.Errorf("%s: Binary bool expression with unsupported operator '%s'", self.getPos(binOp), binOp.GetOperator())
			}
			result = NewBoolVal(true)
			binOp.SetResult(result)
		}
	}

	if result == nil {
		return NewNullVal(), fmt.Errorf("%s: No support for binary expressions of type '%s' and '%s'", self.getPos(binOp), lhs.GetType(), rhs.GetType())
	} else {
		if bashLhs == nil {
			return NewNullVal(), fmt.Errorf("evalBinaryExpr(): LHS is nil")
		}
		if bashRhs == nil {
			return NewNullVal(), fmt.Errorf("evalBinaryExpr(): RHS is nil")
		}
		if isComparison {
			self.bashStmtStack[binOp.GetId()] = bashAst.NewBinaryCompExpr(opType, bashLhs, bashRhs, binOp.GetOperator())
		} else {
			self.bashStmtStack[binOp.GetId()] = bashAst.NewBinaryOpExpr(opType, bashLhs, bashRhs, binOp.GetOperator())
		}
		return result, nil
	}
}

func (self *Transpiler) evalComparisonBinaryExpr(lhs scrilaAst.IRuntimeVal, rhs scrilaAst.IRuntimeVal, operator string) (bashAst.NodeType, error) {
	self.printFuncName("")

	if lhs.GetType() != rhs.GetType() {
		return "", fmt.Errorf("Cannot compare type '%s' and '%s'", lhs.GetType(), rhs.GetType())
	}

	switch lhs.GetType() {
	case scrilaAst.BoolValueType:
		if !slices.Contains([]string{"==", "!="}, operator) {
			return "", fmt.Errorf("Bool comparison does not support operator '%s'", operator)
		}
		return bashAst.BoolLiteralNode, nil
	case scrilaAst.IntValueType:
		if !slices.Contains([]string{">", "<", ">=", "<=", "==", "!="}, operator) {
			return "", fmt.Errorf("Int comparison does not support operator '%s'", operator)
		}
		return bashAst.IntLiteralNode, nil
	case scrilaAst.StrValueType:
		if !slices.Contains([]string{">", "<", "==", "!="}, operator) {
			return "", fmt.Errorf("String comparison does not support operator '%s'", operator)
		}
		return bashAst.StrLiteralNode, nil
	default:
		return "", fmt.Errorf("Comparisons for type '%s' not implemented", lhs.GetType())
	}
}

func (self *Transpiler) evalAssignment(assignment scrilaAst.IAssignmentExpr, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	if assignment.GetAssigne().GetKind() != scrilaAst.IdentifierNode {
		return NewNullVal(), fmt.Errorf("%s: Left side of an assignment must be a variable. Got '%s'", self.getPos(assignment.GetAssigne()), assignment.GetAssigne().GetKind())
	}

	_, err := self.transpile(assignment.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}

	varName := identNodeGetSymbol(assignment.GetAssigne())
	varType, err := env.lookupVarType(varName)
	if err != nil {
		return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(assignment.GetAssigne()), err)
	}

	doMatch, givenType, err := self.exprIsType(assignment.GetValue(), varType, env)
	if err != nil {
		return NewNullVal(), err
	}
	if !doMatch {
		return NewNullVal(), fmt.Errorf("%s: Cannot assign a value of type '%s' to a var of type '%s'", self.getPos(assignment.GetValue()), givenType, varType)
	}

	bashVarType, err := scrilaNodeTypeToBashNodeType(varType)
	if err != nil {
		return NewNullVal(), err
	}
	bashStmt, err := self.exprToRhsBashStmt(assignment.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}
	self.appendUserBody(bashAst.NewAssignmentExpr(
		bashAst.NewVarLiteral(varName, bashVarType),
		bashStmt,
		false,
	))

	result, err := env.assignVar(varName)
	if err != nil {
		return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(assignment), err)
	}
	return result, nil
}

// func (self *Transpiler) evalAssignmentObjMember(assignment scrilaAst.IAssignmentExpr, env *Environment) (scrilaAst.IRuntimeVal, error) {
// 	self.printFuncName("")

// 	if assignment.GetAssigne().GetKind() != scrilaAst.MemberExprNode {
// 		return NewNullVal(), fmt.Errorf("%s: Left side of object member assignment is invalid type '%s'", self.getPos(assignment.GetAssigne()), assignment.GetAssigne().GetKind())
// 	}

// 	memberExpr := scrilaAst.ExprToMemberExpr(assignment.GetAssigne())

// 	if memberExpr.GetObject().GetKind() != scrilaAst.IdentifierNode {
// 		return NewNullVal(), fmt.Errorf("%s: Object name is not the right type. Got '%s'", self.getPos(memberExpr.GetObject()), memberExpr.GetObject().GetKind())
// 	}

// 	if memberExpr.GetProperty().GetKind() != scrilaAst.IdentifierNode {
// 		return NewNullVal(), fmt.Errorf("%s: Object property name is not the right type. Got '%s'", self.getPos(memberExpr.GetProperty()), memberExpr.GetProperty().GetKind())
// 	}

// 	objName := identNodeGetSymbol(memberExpr.GetObject())
// 	obj, err := env.lookupVar(objName)
// 	if err != nil {
// 		return NewNullVal(), err
// 	}
// 	if obj.GetType() != scrilaAst.ObjValueType {
// 		return NewNullVal(), fmt.Errorf("%s: Variable '%s' is not of type 'object'", self.getPos(memberExpr.GetObject()), objName)
// 	}

// 	propName := identNodeGetSymbol(memberExpr.GetProperty())

// 	value, err := self.transpile(assignment.GetValue(), env)
// 	if err != nil {
// 		return NewNullVal(), err
// 	}

// 	switch assignment.GetValue().GetKind() {
// 	case scrilaAst.IntLiteralNode:
// 		self.writeLnTranspilat(value.ToString())
// 	default:
// 		return NewNullVal(), fmt.Errorf("%s: Object member value '%s' is not supported", self.getPos(assignment.GetValue()), assignment.GetValue().GetKind())
// 	}

// 	runtimeToObjVal(obj).GetProperties()[propName] = value
// 	return value, nil
// }

func (self *Transpiler) evalObjectExpr(object scrilaAst.IObjectLiteral, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

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

// func (self *Transpiler) evalMemberExpr(memberExpr scrilaAst.IMemberExpr, env *Environment) (scrilaAst.IRuntimeVal, error) {
// 	self.printFuncName("")

// 	if memberExpr.GetObject().GetKind() != scrilaAst.IdentifierNode {
// 		return NewNullVal(), fmt.Errorf("%s: Object name is not the right type. Got '%s'", self.getPos(memberExpr.GetObject()), memberExpr.GetObject().GetKind())
// 	}

// 	objName := identNodeGetSymbol(memberExpr.GetObject())
// 	obj, err := env.lookupVar(objName)
// 	if err != nil {
// 		return NewNullVal(), err
// 	}
// 	if obj.GetType() != scrilaAst.ObjValueType {
// 		return NewNullVal(), fmt.Errorf("%s: Variable '%s' is not of type 'object'", self.getPos(memberExpr.GetObject()), objName)
// 	}

// 	if memberExpr.GetProperty().GetKind() != scrilaAst.IdentifierNode {
// 		return NewNullVal(), fmt.Errorf("%s: Object property name is not the right type. Got '%s'", self.getPos(memberExpr.GetProperty()), memberExpr.GetProperty().GetKind())
// 	}

// 	propName := identNodeGetSymbol(memberExpr.GetProperty())
// 	result := runtimeToObjVal(obj).GetProperties()[propName]
// 	return result, nil
// }

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

// TODO Rename and move to helpers?
func (self *Transpiler) getCallerResultVarName(call scrilaAst.ICallExpr, env *Environment) (string, error) {
	self.printFuncName("")

	returnType, err := self.getFuncReturnType(call, env)
	if err != nil {
		return "", err
	}

	if returnType == scrilaAst.VoidNode {
		return "", fmt.Errorf("%s: Func '%s' does not have a return value", self.getPos(call.GetCaller()), identNodeGetSymbol(call.GetCaller()))
	}

	resultVarName, err := scrilaNodeTypeToTmpVarName(returnType)
	if err != nil {
		return "", err
	}
	return resultVarName, nil
}

func (self *Transpiler) evalCallExpr(call scrilaAst.ICallExpr, env *Environment) (scrilaAst.IRuntimeVal, error) {
	if call.GetCaller().GetKind() != scrilaAst.IdentifierNode {
		return NewNullVal(), fmt.Errorf("%s: Function name must be an identifier. Got: '%s'", self.getPos(call.GetCaller()), call.GetCaller().GetKind())
	}

	funcName := identNodeGetSymbol(call.GetCaller())

	self.printFuncName(funcName)

	bashArgs := make([]bashAst.IStatement, 0)

	var args []scrilaAst.IRuntimeVal
	for _, arg := range call.GetArgs() {
		evalArg, err := self.transpile(arg, env)
		if err != nil {
			return NewNullVal(), err
		}
		args = append(args, evalArg)

		bashStmt, err := self.exprToRhsBashStmt(arg, env)
		if err != nil {
			return NewNullVal(), err
		}
		bashArgs = append(bashArgs, bashStmt)
	}

	caller, err := env.lookupFunc(funcName)
	if err != nil {
		return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(call), err)
	}

	switch caller.GetType() {
	case scrilaAst.NativeFnType:
		result, err := runtimeToNativeFunc(caller).GetCall()(call.GetArgs(), env)
		if err != nil {
			return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(call), err)
		}

		self.appendUserBody(bashAst.NewCallExpr(funcName, bashArgs))

		return result, nil

	case scrilaAst.FunctionValueType:
		fn := runtimeToFuncVal(caller)

		if len(fn.GetParams()) != len(args) {
			return NewNullVal(), fmt.Errorf("%s: %s(): The amount of passed parameters does not match with the function declaration. Expected: %d, Got: %d", self.getPos(call), fn.GetName(), len(fn.GetParams()), len(args))
		}

		self.appendUserBody(bashAst.NewCallExpr(funcName, bashArgs))

		for i, param := range fn.GetParams() {
			if !scrilaAst.DoTypesMatch(param.GetParamType(), args[i].GetType()) {
				return NewNullVal(), fmt.Errorf("%s: %s(): Parameter '%s' type does not match. Expected: %s, Got: %s", self.getPos(call), fn.GetName(), param.GetName(), param.GetParamType(), args[i].GetType())
			}
		}

		result, err := scrilaNodeTypeToRuntimeVal(fn.GetReturnType())
		if err != nil {
			return NewNullVal(), err
		}
		return result, nil

	default:
		return NewNullVal(), fmt.Errorf("%s: Cannot call value that is not a function: %s", self.getPos(call), caller)
	}
}

func (self *Transpiler) evalWhileExitKeywords(expr scrilaAst.IExpr, env *Environment) (scrilaAst.IRuntimeVal, error) {
	if !self.contextContains(WhileLoopContext) {
		return NewNullVal(), fmt.Errorf("%s: '%s' is only allowed inside a while loop", self.getPos(expr), expr.GetKind())
	}

	bashStmt, err := self.exprToBashStmt(expr, env)
	if err != nil {
		return NewNullVal(), err
	}
	self.appendUserBody(bashStmt)

	return NewNullVal(), nil
}

func (self *Transpiler) evalReturnExpr(returnExpr scrilaAst.IReturnExpr, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	// Check if transpiler is in function context otherwise `return` is not allowed
	if !self.contextContains(FunctionContext) || self.currentFunc == nil {
		return NewNullVal(), fmt.Errorf("%s: Return is only allowed inside a function", self.getPos(returnExpr))
	}

	// Check if functions of type "void" do not have a return expression with value
	if self.currentFunc.GetReturnType() == scrilaAst.VoidNode {
		if !returnExpr.IsEmpty() {
			return NewNullVal(), fmt.Errorf("%s: %s(): Cannot return value if function type is 'void'", self.getPos(returnExpr), self.currentFunc.GetName())
		}

		self.appendUserBody(bashAst.NewReturnExpr())
		return NewNullVal(), nil
	}

	// Check if functions with type other than "void" do not have a return expression without value
	if returnExpr.IsEmpty() {
		return NewNullVal(), fmt.Errorf("%s: %s(): Cannot return without a value for a function with return value", self.getPos(returnExpr), self.currentFunc.GetName())
	}

	value, err := self.transpile(returnExpr.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}

	// Check if the return value matches with the function type
	if !scrilaAst.DoTypesMatch(self.currentFunc.GetReturnType(), value.GetType()) {
		return NewNullVal(), fmt.Errorf("%s: %s(): Return type does not match with function type. Expected: %s, Got: %s", self.getPos(returnExpr), self.currentFunc.GetName(), self.currentFunc.GetReturnType(), value.GetType())
	}

	resultValue, err := self.exprToRhsBashStmt(returnExpr.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}

	resultVarName, err := scrilaNodeTypeToTmpVarName(self.currentFunc.GetReturnType())
	if err != nil {
		return NewNullVal(), err
	}

	resultVarType, err := scrilaNodeTypeToBashNodeType(self.currentFunc.GetReturnType())
	if err != nil {
		return NewNullVal(), err
	}

	if resultValue.GetKind() != bashAst.VarLiteralNode || bashAst.StmtToVarLiteral(resultValue).GetDataType() != resultVarType {
		self.appendUserBody(bashAst.NewAssignmentExpr(
			bashAst.NewVarLiteral(resultVarName, resultVarType),
			resultValue,
			false,
		))
	}
	self.appendUserBody(bashAst.NewReturnExpr())
	return value, nil
}
