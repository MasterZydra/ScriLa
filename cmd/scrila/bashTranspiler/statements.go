package bashTranspiler

import (
	"ScriLa/cmd/scrila/bashAst"
	"ScriLa/cmd/scrila/scrilaAst"
	"fmt"
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

	_, err := self.transpile(varDeclaration.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}

	// Check if variable type and value type match
	doMatch, givenType, err := self.exprIsType(varDeclaration.GetValue(), varDeclaration.GetDataType(), env)
	if err != nil {
		return NewNullVal(), err
	}
	if !doMatch {
		return NewNullVal(), fmt.Errorf("%s: Cannot assign a value of type '%s' to a var of type '%s'", self.getPos(varDeclaration.GetValue()), givenType, varDeclaration.GetDataType())
	}

	// Same logic in evalAssignment -> merge into one function
	bashVarType, err := scrilaNodeTypeToBashNodeType(varDeclaration.GetDataType())
	if err != nil {
		return NewNullVal(), err
	}
	bashStmt, err := self.exprToRhsBashStmt(varDeclaration.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}
	self.appendUserBody(bashAst.NewAssignmentExpr(
		bashAst.NewVarLiteral(varDeclaration.GetIdentifier(), bashVarType),
		bashStmt,
		true,
	))

	result, err := env.declareVar(varDeclaration.GetIdentifier(), varDeclaration.IsConstant(), varDeclaration.GetDataType())
	if err != nil {
		return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(varDeclaration), err)
	}
	return result, nil
}

func (self *Transpiler) evalForStatement(forStmt scrilaAst.IForStatement, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	// Index
	varType := forStmt.GetIndexVarType()
	varName := forStmt.GetIndex().GetSymbol()

	localEnv := NewEnvironment(env, self)

	_, err := localEnv.declareVar(varName, false, varType)
	if err != nil {
		return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(forStmt), err)
	}
	bashVarType, err := scrilaNodeTypeToBashNodeType(varType)
	if err != nil {
		return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(forStmt), err)
	}
	varLiteral := bashAst.NewVarLiteral(varName, bashVarType)

	// Array
	_, err = self.transpile(forStmt.GetArray(), env)
	if err != nil {
		return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(forStmt), err)
	}
	doMatch, err := self.exprIsArray(forStmt.GetArray(), varType, localEnv)
	if err != nil {
		return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(forStmt), err)
	}
	if !doMatch {
		return NewNullVal(), fmt.Errorf("%s: Array data type and index data type is not matching", self.getPos(forStmt))
	}
	_, err = self.transpile(forStmt.GetArray(), localEnv)
	bashArrayStmt, err := self.exprToRhsBashStmt(forStmt.GetArray(), localEnv)

	self.pushContext(ForLoopContext)
	self.pushBashContext(bashAst.NewForStmt(varLiteral, bashArrayStmt))

	// Transpile the body line by line
	err = self.evalStatementBody(forStmt.GetBody(), localEnv)
	if err != nil {
		return NewNullVal(), err
	}

	bashForStmt := self.currentBashContext()
	self.popContext()
	self.popBashContext()
	self.appendUserBody(bashForStmt)

	return NewNullVal(), err
}

func (self *Transpiler) evalIfStatement(ifStatement scrilaAst.IIfStatement, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	// Transpile condition
	self.pushCallArgIndex()
	_, err := self.transpile(ifStatement.GetCondition(), env)
	if err != nil {
		return NewNullVal(), err
	}
	self.popCallArgIndex()
	err = self.evalStatementCondition(ifStatement.GetCondition(), env)
	if err != nil {
		return NewNullVal(), err
	}
	bashCond, ok := self.bashStmtStack[ifStatement.GetCondition().GetId()]
	if !ok {
		return NewNullVal(), fmt.Errorf("evalIfStatement(): Condition is not stored in stack")
	}

	self.pushContext(IfStmtContext)
	self.pushBashContext(bashAst.NewIfStmt(bashCond))

	// Transpile the body line by line
	err = self.evalStatementBody(ifStatement.GetBody(), env)
	if err != nil {
		return NewNullVal(), err
	}

	ifStmt := self.currentBashContext()
	self.popContext()
	self.popBashContext()

	// Else block
	if ifStatement.GetElse() != nil {
		if err = self.evalIfStatementElse(ifStatement.GetElse(), env); err != nil {
			return NewNullVal(), err
		}
		elseBlock, ok := self.bashStmtStack[ifStatement.GetElse().GetId()]
		if !ok {
			return NewNullVal(), fmt.Errorf("evalIfStatement(): ElseIf is not stored in stack")
		}
		bashAst.StmtToIfStmt(ifStmt).SetElse(bashAst.StmtToIfStmt(elseBlock))
	}

	self.appendUserBody(ifStmt)

	return NewNullVal(), nil
}

func (self *Transpiler) evalWhileStatement(whileStatement scrilaAst.IWhileStatement, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	// Transpile condition
	self.pushCallArgIndex()
	_, err := self.transpile(whileStatement.GetCondition(), env)
	if err != nil {
		return NewNullVal(), err
	}
	self.popCallArgIndex()
	err = self.evalStatementCondition(whileStatement.GetCondition(), env)
	if err != nil {
		return NewNullVal(), err
	}
	bashCond, ok := self.bashStmtStack[whileStatement.GetCondition().GetId()]
	if !ok {
		return NewNullVal(), fmt.Errorf("evalWhileStatement(): Condition is not stored in stack")
	}

	self.pushContext(WhileLoopContext)
	self.pushBashContext(bashAst.NewWhileStmt(bashCond))

	// Transpile the body line by line
	err = self.evalStatementBody(whileStatement.GetBody(), env)
	if err != nil {
		return NewNullVal(), err
	}

	whileStmt := self.currentBashContext()
	self.popContext()
	self.popBashContext()
	self.appendUserBody(whileStmt)
	return NewNullVal(), nil
}

func (self *Transpiler) evalIfStatementElse(elseBlock scrilaAst.IIfStatement, env *Environment) error {
	// TODO Merge with evalIfStatement - Add param "isElse bool"
	self.printFuncName("")

	// Else if
	var bashCond bashAst.IStatement
	if elseBlock.GetCondition() != nil {
		// Transpile condition
		self.pushCallArgIndex()
		_, err := self.transpile(elseBlock.GetCondition(), env)
		if err != nil {
			return err
		}
		self.popCallArgIndex()
		err = self.evalStatementCondition(elseBlock.GetCondition(), env)
		if err != nil {
			return err
		}

		var ok bool
		bashCond, ok = self.bashStmtStack[elseBlock.GetCondition().GetId()]
		if !ok {
			return fmt.Errorf("evalIfStatementElse(): Condition is not stored in stack")
		}
	}

	self.pushContext(IfStmtContext)
	self.pushBashContext(bashAst.NewIfStmt(bashCond))

	// Transpile the body line by line
	err := self.evalStatementBody(elseBlock.GetBody(), env)
	if err != nil {
		return err
	}

	ifStmt := self.currentBashContext()
	self.popContext()
	self.popBashContext()

	if elseBlock.GetElse() != nil {
		if err = self.evalIfStatementElse(elseBlock.GetElse(), env); err != nil {
			return err
		}
		elseBlock, ok := self.bashStmtStack[elseBlock.GetElse().GetId()]
		if !ok {
			return fmt.Errorf("evalIfStatementElse(): ElseIf is not stored in stack")
		}
		bashAst.StmtToIfStmt(ifStmt).SetElse(bashAst.StmtToIfStmt(elseBlock))
	}

	self.bashStmtStack[elseBlock.GetId()] = ifStmt

	return nil
}

func (self *Transpiler) evalStatementCondition(condition scrilaAst.IExpr, env *Environment) error {
	self.printFuncName("")

	// Check if condition is of type boolean
	doMatch, givenType, err := self.exprIsType(condition, scrilaAst.BoolLiteralNode, env)
	if err != nil {
		return err
	}
	if !doMatch {
		return fmt.Errorf("%s: Condition is not of type bool. Got %s", self.getPos(condition), givenType)
	}

	bashCond, err := self.exprToBashStmt(condition, env)
	if err != nil {
		return err
	}

	self.bashStmtStack[condition.GetId()] = bashCond
	return nil
}

func (self *Transpiler) evalStatementBody(body []scrilaAst.IStatement, env *Environment) error {
	self.printFuncName("")

	scope := NewEnvironment(env, self)
	for _, stmt := range body {
		_, err := self.transpile(stmt, scope)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Transpiler) evalFunctionDeclaration(funcDeclaration scrilaAst.IFunctionDeclaration, env *Environment) (scrilaAst.IRuntimeVal, error) {
	self.printFuncName("")

	if self.contextContains(FunctionContext) {
		return NewNullVal(), fmt.Errorf("%s: Cannot declare a function inside a function", self.getPos(funcDeclaration))
	}

	fn := NewFunctionVal(funcDeclaration, env)
	scope := NewEnvironment(fn.GetDeclarationEnv(), self)

	self.pushContext(FunctionContext)
	bashReturnType, err := scrilaNodeTypeToBashNodeType(funcDeclaration.GetReturnType())
	if err != nil {
		return NewNullVal(), err
	}
	self.currentBashFunc = bashAst.NewFuncDeclaration(funcDeclaration.GetName(), bashReturnType)
	self.currentFunc = fn

	for i, param := range funcDeclaration.GetParameters() {
		paramType, err := scrilaNodeTypeToBashNodeType(param.GetParamType())
		if err != nil {
			return NewNullVal(), err
		}
		self.currentBashFunc.AppendParams(bashAst.NewFuncParameter(param.GetName(), paramType))

		_, err = scope.declareVar(fn.GetParams()[i].GetName(), false, fn.GetParams()[i].GetParamType())
		if err != nil {
			return NewNullVal(), fmt.Errorf("%s: %s", self.getPos(funcDeclaration), err)
		}
	}

	// Transpile the function body line by line
	var result scrilaAst.IRuntimeVal
	result = NewNullVal()
	for _, stmt := range fn.GetBody() {
		var err error
		result, err = self.transpile(stmt, scope)
		if err != nil {
			return NewNullVal(), err
		}
	}
	self.popContext()
	self.bashProgram.AppendUserBody(self.currentBashFunc)
	self.currentBashFunc = nil
	self.currentFunc = nil

	_, err = env.declareFunc(funcDeclaration.GetName(), fn)
	if err != nil {
		return result, fmt.Errorf("%s: %s", self.getPos(funcDeclaration), err)
	}
	return result, nil
}
