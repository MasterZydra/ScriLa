package bashTranspiler

import (
	"ScriLa/cmd/scrila/lexer"
	"ScriLa/cmd/scrila/scrilaAst"
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"
)

func (self *Transpiler) evalProgram(program scrilaAst.IProgram, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	var lastEvaluated scrilaAst.IRuntimeVal = NewNullVal()

	for _, statement := range program.GetBody() {
		var err error
		lastEvaluated, err = self.transpile(statement, env)
		if err != nil {
			return NewNullVal(), err
		}
	}

	return lastEvaluated, nil
}

func (self *Transpiler) evalVarDeclaration(varDeclaration scrilaAst.IVarDeclaration, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	value, err := self.transpile(varDeclaration.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}

	if varDeclaration.GetValue().GetKind() == scrilaAst.BinaryExprNode && scrilaAst.BinExprReturnsBool(scrilaAst.ExprToBinExpr(varDeclaration.GetValue())) {
		self.writeLnTranspilat(binCompExpValueToBashIf(value))
	}

	if self.contextContains(FunctionContext) {
		self.writeTranspilat("local ")
	}
	if varDeclaration.GetValue().GetKind() == scrilaAst.ObjectLiteralNode {
		self.writeLnTranspilat("declare -A " + varDeclaration.GetIdentifier())
	} else {
		self.writeTranspilat(varDeclaration.GetIdentifier() + "=")
	}

	switch varDeclaration.GetValue().GetKind() {
	case scrilaAst.CallExprNode:
		returnType, err := self.getFuncReturnType(scrilaAst.ExprToCallExpr(varDeclaration.GetValue()), env)
		if err != nil {
			return NewNullVal(), err
		}
		if returnType != varDeclaration.GetVarType() {
			return NewNullVal(), fmt.Errorf("%s: Cannot assign a value of type '%s' to a var of type '%s'", self.getPos(varDeclaration.GetValue()), returnType, varDeclaration.GetVarType())
		}

		varName, err := self.getCallerResultVarName(scrilaAst.ExprToCallExpr(varDeclaration.GetValue()), env)
		if err != nil {
			return NewNullVal(), err
		}
		switch varDeclaration.GetVarType() {
		case lexer.StrType:
			self.writeLnTranspilat(strToBashStr(varName))
			value = NewStrVal("")
		case lexer.IntType:
			self.writeLnTranspilat(varName)
			value = NewIntVal(1)
		case lexer.BoolType:
			self.writeLnTranspilat(strToBashStr(varName))
			value = NewBoolVal(true)
		default:
			return NewNullVal(), fmt.Errorf("%s: Assigning return values is not implemented for variables of type '%s'", self.getPos(varDeclaration), varDeclaration.GetVarType())
		}

	case scrilaAst.IdentifierNode:
		symbol := identNodeGetSymbol(varDeclaration.GetValue())
		if symbol == "null" || scrilaAst.IdentIsBool(scrilaAst.ExprToIdent(varDeclaration.GetValue())) {
			self.writeLnTranspilat(strToBashStr(symbol))
		} else if slices.Contains(reservedIdentifiers, symbol) {
			self.writeLnTranspilat(symbol)
		} else {
			valueVarType, err := env.lookupVarType(identNodeGetSymbol(varDeclaration.GetValue()))
			if err != nil {
				return NewNullVal(), err
			}
			if valueVarType != varDeclaration.GetVarType() {
				return NewNullVal(), fmt.Errorf("%s: Cannot assign a value of type '%s' to a var of type '%s'", self.getPos(varDeclaration.GetValue()), valueVarType, varDeclaration.GetVarType())
			}
			switch varDeclaration.GetVarType() {
			case lexer.StrType:
				self.writeLnTranspilat(strToBashStr(identNodeToBashVar(varDeclaration.GetValue())))
			case lexer.IntType:
				self.writeLnTranspilat(identNodeToBashVar(varDeclaration.GetValue()))
			default:
				return NewNullVal(), fmt.Errorf("%s: Assigning variables is not implemented for variables of type '%s'", self.getPos(varDeclaration), varDeclaration.GetVarType())
			}
		}
	case scrilaAst.BinaryExprNode:
		switch varDeclaration.GetVarType() {
		case lexer.StrType:
			self.writeLnTranspilat(strToBashStr(value.GetTranspilat()))
		case lexer.IntType:
			self.writeLnTranspilat(value.GetTranspilat())
		case lexer.BoolType:
			if scrilaAst.BinExprReturnsBool(scrilaAst.ExprToBinExpr(varDeclaration.GetValue())) {
				self.writeLnTranspilat("${tmpBool}")
			} else {
				self.writeLnTranspilat(value.GetTranspilat())
			}
		default:
			return NewNullVal(), fmt.Errorf("%s: Assigning binary expressions is not implemented for variables of type '%s'", self.getPos(varDeclaration), varDeclaration.GetVarType())
		}
	case scrilaAst.StrLiteralNode:
		if varDeclaration.GetVarType() != lexer.StrType {
			return NewNullVal(), fmt.Errorf("%s: Cannot assign a value of type '%s' to a var of type '%s'", self.getPos(varDeclaration.GetValue()), lexer.StrType, varDeclaration.GetVarType())
		}
		self.writeLnTranspilat(strToBashStr(value.ToString()))
	case scrilaAst.IntLiteralNode:
		if varDeclaration.GetVarType() != lexer.IntType {
			return NewNullVal(), fmt.Errorf("%s: Cannot assign a value of type '%s' to a var of type '%s'", self.getPos(varDeclaration.GetValue()), lexer.IntType, varDeclaration.GetVarType())
		}
		self.writeLnTranspilat(value.ToString())
	case scrilaAst.ObjectLiteralNode:
		for _, prop := range scrilaAst.ExprToObjLit(varDeclaration.GetValue()).GetProperties() {
			self.writeTranspilat(varDeclaration.GetIdentifier() + "[" + strToBashStr(prop.GetKey()) + "]=")
			value, err := self.transpile(prop.GetValue(), env)
			if err != nil {
				return NewNullVal(), err
			}
			switch prop.GetValue().GetKind() {
			case scrilaAst.IntLiteralNode:
				self.writeLnTranspilat(value.ToString())
			case scrilaAst.StrLiteralNode:
				self.writeLnTranspilat(strToBashStr(value.ToString()))
			case scrilaAst.IdentifierNode:
				symbol := identNodeGetSymbol(prop.GetValue())
				if symbol == "null" {
					self.writeLnTranspilat(strToBashStr(symbol))
				} else if slices.Contains(reservedIdentifiers, symbol) {
					self.writeLnTranspilat(symbol)
				} else {
					self.writeLnTranspilat(identNodeToBashVar(prop.GetValue()))
				}
			default:
				return NewNullVal(), fmt.Errorf("%s: Assigning object properties of type '%s' is not implemented", self.getPos(varDeclaration), prop.GetValue().GetKind())
			}
		}
	case scrilaAst.MemberExprNode:
		memberVal, err := self.evalMemberExpr(scrilaAst.ExprToMemberExpr(varDeclaration.GetValue()), env)
		if err != nil {
			return NewNullVal(), err
		}
		self.writeLnTranspilat(memberVal.GetTranspilat())
	default:
		return NewNullVal(), fmt.Errorf("%s: Assigning value of type '%s' is not implemented", self.getPos(varDeclaration), varDeclaration.GetValue().GetKind())
	}

	result, err := env.declareVar(varDeclaration.GetIdentifier(), value, varDeclaration.IsConstant(), varDeclaration.GetVarType())
	if err != nil {
		return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(varDeclaration), err)
	}
	return result, nil
}

func (self *Transpiler) evalIfStatement(ifStatement scrilaAst.IIfStatement, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	_, err := self.transpile(ifStatement.GetCondition(), env)
	if err != nil {
		return NewNullVal(), err
	}

	self.writeTranspilat("if ")
	self.pushContext(IfStmtContext)

	// Transpile condition
	err = self.evalStatementCondition(ifStatement.GetCondition(), env)
	if err != nil {
		return NewNullVal(), err
	}

	// Transpile the body line by line
	err = self.evalStatementBody(ifStatement.GetBody(), env)
	if err != nil {
		return NewNullVal(), err
	}

	// Else block
	self.evalIfStatementElse(ifStatement.GetElse(), env)

	self.popContext()
	self.writeLnTranspilat(self.indent(0) + "fi")
	return NewNullVal(), nil
}

func (self *Transpiler) evalWhileStatement(whileStatement scrilaAst.IWhileStatement, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	_, err := self.transpile(whileStatement.GetCondition(), env)
	if err != nil {
		return NewNullVal(), err
	}

	self.writeTranspilat("while ")
	self.pushContext(WhileLoopContext)

	// Transpile condition
	err = self.evalStatementCondition(whileStatement.GetCondition(), env)
	if err != nil {
		return NewNullVal(), err
	}

	// Transpile the body line by line
	err = self.evalStatementBody(whileStatement.GetBody(), env)
	if err != nil {
		return NewNullVal(), err
	}

	self.popContext()
	self.writeLnTranspilat("done")
	return NewNullVal(), nil
}

func (self *Transpiler) evalIfStatementElse(elseBlock scrilaAst.IIfStatement, env *Environment) error {
	self.printFuncName("")

	if elseBlock == nil {
		return nil
	}

	// Else if
	if elseBlock.GetCondition() != nil {
		self.writeTranspilat("elif ")
		// Transpile condition
		err := self.evalStatementCondition(elseBlock.GetCondition(), env)
		if err != nil {
			return err
		}
	} else {
		self.writeLnTranspilat("else")
	}

	// Transpile the body line by line
	err := self.evalStatementBody(elseBlock.GetBody(), env)
	if err != nil {
		return err
	}

	return self.evalIfStatementElse(elseBlock.GetElse(), env)
}

func (self *Transpiler) evalStatementCondition(condition scrilaAst.IExpr, env *Environment) error {
	self.printFuncName("")

	switch condition.GetKind() {
	case scrilaAst.BinaryExprNode:
		value, err := self.transpile(condition, env)
		if err != nil {
			return err
		}
		if value.GetType() != scrilaAst.BoolValueType {
			return fmt.Errorf("%s: Condition is no boolean expression. Got %s", self.getPos(condition), value.GetType())
		}
		self.writeLnTranspilat(value.GetTranspilat())
	case scrilaAst.CallExprNode:
		returnType, err := self.getFuncReturnType(scrilaAst.ExprToCallExpr(condition), env)
		if err != nil {
			return err
		}
		if returnType != lexer.BoolType {
			return fmt.Errorf("%s: Cannot use a value of type '%s' as condition", self.getPos(condition), returnType)
		}

		varName, err := self.getCallerResultVarName(scrilaAst.ExprToCallExpr(condition), env)
		if err != nil {
			return err
		}
		self.writeLnTranspilat(strToBashStr(varName))
	case scrilaAst.IdentifierNode:
		identifier := scrilaAst.ExprToIdent(condition)
		if scrilaAst.IdentIsBool(identifier) {
			self.writeLnTranspilat(boolIdentToBashComparison(identifier))
		} else {
			valueVarType, err := env.lookupVarType(identNodeGetSymbol(condition))
			if err != nil {
				return err
			}

			if valueVarType != lexer.BoolType {
				return fmt.Errorf("%s: Condition is not of type bool. Got %s", self.getPos(condition), valueVarType)
			}
			self.writeLnTranspilat(varIdentToBashComparison(identifier))
		}
	default:
		return fmt.Errorf("%s: Unsupported type '%s' for condition", self.getPos(condition), condition.GetKind())
	}
	switch self.currentContext() {
	case IfStmtContext:
		self.writeLnTranspilat(self.indent(1) + "then")
	case WhileLoopContext:
		self.writeLnTranspilat(self.indent(1) + "do")
	default:
		return fmt.Errorf("%s: Unsupported context '%s' for condition", self.getPos(condition), self.currentContext())
	}
	return nil
}

func (self *Transpiler) evalStatementBody(body []scrilaAst.IStatement, env *Environment) error {
	self.printFuncName("")

	// Bash does not support an empty if-/while-body. A fix is to up a ":" inside the body.
	if len(body) == 0 {
		self.writeLnTranspilat(self.indent(0) + ":")
		return nil
	}

	scope := NewEnvironment(env, self)
	for _, stmt := range body {
		self.writeTranspilat(self.indent(0))
		_, err := self.transpile(stmt, scope)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Transpiler) evalFunctionDeclaration(funcDeclaration scrilaAst.IFunctionDeclaration, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	fn := NewFunctionVal(funcDeclaration, env)
	scope := NewEnvironment(fn.GetDeclarationEnv(), self)

	self.pushContext(FunctionContext)

	self.writeLnTranspilat(funcDeclaration.GetName() + " () {")
	for i, param := range funcDeclaration.GetParameters() {
		var value scrilaAst.IRuntimeVal
		switch fn.GetParams()[i].GetParamType() {
		case lexer.IntType:
			value = NewIntVal(1)
		case lexer.StrType:
			value = NewStrVal("str")
		default:
			return NewNullVal(), fmt.Errorf("%s: Unsupported type '%s' for parameter '%s'", self.getPos(funcDeclaration), fn.GetParams()[i].GetParamType(), fn.GetParams()[i].GetName())
		}
		_, err := scope.declareVar(fn.GetParams()[i].GetName(), value, false, fn.GetParams()[i].GetParamType())
		if err != nil {
			return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(funcDeclaration), err)
		}
		self.writeLnTranspilat(self.indent(0) + "local " + param.GetName() + "=$" + strconv.Itoa(i+1))
	}

	// Transpile the function body line by line
	self.currentFunc = fn
	var result scrilaAst.IRuntimeVal
	result = NewNullVal()
	for _, stmt := range fn.GetBody() {
		var err error
		self.writeTranspilat(self.indent(0))
		result, err = self.transpile(stmt, scope)
		if err != nil {
			return NewNullVal(), err
		}
	}
	self.popContext()
	self.currentFunc = nil

	self.writeLnTranspilat("}\n")
	_, err := env.declareFunc(funcDeclaration.GetName(), fn)
	if err != nil {
		return result, fmt.Errorf("%s: %s", self.getPos(funcDeclaration), err)
	}
	return result, nil
}
