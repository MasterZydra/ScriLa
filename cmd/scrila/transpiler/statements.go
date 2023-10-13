package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"ScriLa/cmd/scrila/lexer"
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"
)

func evalProgram(program ast.IProgram, env *Environment) (IRuntimeVal, error) {
	var lastEvaluated IRuntimeVal = NewNullVal()

	for _, statement := range program.GetBody() {
		var err error
		lastEvaluated, err = transpile(statement, env)
		if err != nil {
			return NewNullVal(), err
		}
	}

	return lastEvaluated, nil
}

func evalVarDeclaration(varDeclaration ast.IVarDeclaration, env *Environment) (IRuntimeVal, error) {
	value, err := transpile(varDeclaration.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}
	if funcContext {
		writeToFile("local ")
	}
	if varDeclaration.GetValue().GetKind() == ast.ObjectLiteralNode {
		writeLnToFile("declare -A " + varDeclaration.GetIdentifier())
	} else {
		writeToFile(varDeclaration.GetIdentifier() + "=")
	}

	switch varDeclaration.GetValue().GetKind() {
	case ast.CallExprNode:
		varName, err := getCallerResultVarName(ast.ExprToCallExpr(varDeclaration.GetValue()), env)
		if err != nil {
			return NewNullVal(), err
		}
		switch varDeclaration.GetVarType() {
		case lexer.StrType:
			writeLnToFile(strToBashStr(varName))
			value = NewStrVal("")
		case lexer.IntType:
			writeLnToFile(varName)
			value = NewIntVal(1)
		default:
			return NewNullVal(), fmt.Errorf("%s:%d:%d: Assigning return values is not implemented for variables of type '%s'", fileName, varDeclaration.GetLn(), varDeclaration.GetCol(), varDeclaration.GetVarType())
		}

	case ast.IdentifierNode:
		symbol := identNodeGetSymbol(varDeclaration.GetValue())
		if symbol == "null" || ast.IdentIsBool(ast.ExprToIdent(varDeclaration.GetValue())) {
			writeLnToFile(strToBashStr(symbol))
		} else if slices.Contains(reservedIdentifiers, symbol) {
			writeLnToFile(symbol)
		} else {
			valueVarType, err := env.lookupVarType(identNodeGetSymbol(varDeclaration.GetValue()))
			if err != nil {
				return NewNullVal(), err
			}
			if valueVarType != varDeclaration.GetVarType() {
				return NewNullVal(), fmt.Errorf("%s:%d:%d: Cannot assign a value of type '%s' to a var of type '%s'", fileName, varDeclaration.GetValue().GetLn(), varDeclaration.GetValue().GetCol(), valueVarType, varDeclaration.GetVarType())
			}
			switch varDeclaration.GetVarType() {
			case lexer.StrType:
				writeLnToFile(strToBashStr(identNodeToBashVar(varDeclaration.GetValue())))
			case lexer.IntType:
				writeLnToFile(identNodeToBashVar(varDeclaration.GetValue()))
			default:
				return NewNullVal(), fmt.Errorf("%s:%d:%d: Assigning variables is not implemented for variables of type '%s'", fileName, varDeclaration.GetLn(), varDeclaration.GetCol(), varDeclaration.GetVarType())
			}
		}
	case ast.BinaryExprNode:
		switch varDeclaration.GetVarType() {
		case lexer.StrType:
			writeLnToFile(strToBashStr(value.GetTranspilat()))
		case lexer.IntType, lexer.BoolType:
			writeLnToFile(value.GetTranspilat())
		default:
			return NewNullVal(), fmt.Errorf("%s:%d:%d: Assigning binary expressions is not implemented for variables of type '%s'", fileName, varDeclaration.GetLn(), varDeclaration.GetCol(), varDeclaration.GetVarType())
		}
	case ast.StrLiteralNode:
		if varDeclaration.GetVarType() != lexer.StrType {
			return NewNullVal(), fmt.Errorf("%s:%d:%d: Cannot assign a value of type '%s' to a var of type '%s'", fileName, varDeclaration.GetValue().GetLn(), varDeclaration.GetValue().GetCol(), lexer.StrType, varDeclaration.GetVarType())
		}
		writeLnToFile(strToBashStr(value.ToString()))
	case ast.IntLiteralNode:
		if varDeclaration.GetVarType() != lexer.IntType {
			return NewNullVal(), fmt.Errorf("%s:%d:%d: Cannot assign a value of type '%s' to a var of type '%s'", fileName, varDeclaration.GetValue().GetLn(), varDeclaration.GetValue().GetCol(), lexer.IntType, varDeclaration.GetVarType())
		}
		writeLnToFile(value.ToString())
	case ast.ObjectLiteralNode:
		for _, prop := range ast.ExprToObjLit(varDeclaration.GetValue()).GetProperties() {
			writeToFile(varDeclaration.GetIdentifier() + "[" + strToBashStr(prop.GetKey()) + "]=")
			value, err := transpile(prop.GetValue(), env)
			if err != nil {
				return NewNullVal(), err
			}
			switch prop.GetValue().GetKind() {
			case ast.IntLiteralNode:
				writeLnToFile(value.ToString())
			case ast.StrLiteralNode:
				writeLnToFile(strToBashStr(value.ToString()))
			case ast.IdentifierNode:
				symbol := identNodeGetSymbol(prop.GetValue())
				if symbol == "null" {
					writeLnToFile(strToBashStr(symbol))
				} else if slices.Contains(reservedIdentifiers, symbol) {
					writeLnToFile(symbol)
				} else {
					writeLnToFile(identNodeToBashVar(prop.GetValue()))
				}
			default:
				return NewNullVal(), fmt.Errorf("%s:%d:%d: Assigning object properties of type '%s' is not implemented", fileName, varDeclaration.GetLn(), varDeclaration.GetCol(), prop.GetValue().GetKind())
			}
		}
	case ast.MemberExprNode:
		memberVal, err := evalMemberExpr(ast.ExprToMemberExpr(varDeclaration.GetValue()), env)
		if err != nil {
			return NewNullVal(), err
		}
		writeLnToFile(memberVal.GetTranspilat())
	default:
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Assigning value of type '%s' is not implemented", fileName, varDeclaration.GetLn(), varDeclaration.GetCol(), varDeclaration.GetValue().GetKind())
	}

	result, err := env.declareVar(varDeclaration.GetIdentifier(), value, varDeclaration.IsConstant(), varDeclaration.GetVarType())
	if err != nil {
		return NewNullVal(), fmt.Errorf("%s:%d:%d: %s", fileName, varDeclaration.GetLn(), varDeclaration.GetCol(), err)
	}
	return result, nil
}

func evalIfStatement(ifStatement ast.IIfStatement, env *Environment) (IRuntimeVal, error) {
	writeToFile("if ")

	// Transpile condition
	switch ifStatement.GetCondition().GetKind() {
	case ast.BinaryExprNode:
		value, err := transpile(ifStatement.GetCondition(), env)
		if err != nil {
			return NewNullVal(), err
		}
		if value.GetType() != BoolValueType {
			return NewNullVal(), fmt.Errorf("%s:%d:%d: Condition is no boolean expression. Got %s", fileName, ifStatement.GetCondition().GetLn(), ifStatement.GetCondition().GetCol(), value.GetType())
		}
		writeToFile(value.GetTranspilat())
	case ast.IdentifierNode:
		identifier := ast.ExprToIdent(ifStatement.GetCondition())
		if ast.IdentIsBool(identifier) {
			writeToFile(boolIdentToBashComparison(identifier))
		} else {
			valueVarType, err := env.lookupVarType(identNodeGetSymbol(ifStatement.GetCondition()))
			if err != nil {
				return NewNullVal(), err
			}

			if valueVarType != lexer.BoolType {
				return NewNullVal(), fmt.Errorf("%s:%d:%d: Condition is not of type bool. Got %s", fileName, ifStatement.GetCondition().GetLn(), ifStatement.GetCondition().GetCol(), valueVarType)
			}
			writeToFile(varIdentToBashComparison(identifier))
		}
	default:
		return NewNullVal(), fmt.Errorf("%s:%d:%d: Unsupported type '%s' for condition", fileName, ifStatement.GetCondition().GetLn(), ifStatement.GetCondition().GetCol(), ifStatement.GetCondition().GetKind())
	}
	writeLnToFile("; then")

	// Transpile the block line by line
	scope := NewEnvironment(NewEnvironment(env))
	for _, stmt := range ifStatement.GetBody() {
		writeToFile("\t")
		_, err := transpile(stmt, scope)
		if err != nil {
			return NewNullVal(), err
		}
	}

	writeLnToFile("fi")
	return NewNullVal(), nil
}

func evalFunctionDeclaration(funcDeclaration ast.IFunctionDeclaration, env *Environment) (IRuntimeVal, error) {
	fn := NewFunctionVal(funcDeclaration, env)
	scope := NewEnvironment(fn.GetDeclarationEnv())

	writeLnToFile(funcDeclaration.GetName() + " () {")
	for i, param := range funcDeclaration.GetParameters() {
		var value IRuntimeVal
		switch fn.GetParams()[i].GetParamType() {
		case lexer.IntType:
			value = NewIntVal(1)
		case lexer.StrType:
			value = NewStrVal("str")
		default:
			return NewNullVal(), fmt.Errorf("%s:%d:%d: Unsupported type '%s' for parameter '%s'", fileName, funcDeclaration.GetLn(), funcDeclaration.GetCol(), fn.GetParams()[i].GetParamType(), fn.GetParams()[i].GetName())
		}
		_, err := scope.declareVar(fn.GetParams()[i].GetName(), value, false, fn.GetParams()[i].GetParamType())
		if err != nil {
			return NewNullVal(), fmt.Errorf("%s:%d:%d: %s", fileName, funcDeclaration.GetLn(), funcDeclaration.GetCol(), err)
		}
		writeLnToFile("\tlocal " + param.GetName() + "=$" + strconv.Itoa(i+1))
	}

	// Transpile the function body line by line
	funcContext = true
	currentFunc = fn
	var result IRuntimeVal
	result = NewNullVal()
	for _, stmt := range fn.GetBody() {
		var err error
		writeToFile("\t")
		result, err = transpile(stmt, scope)
		if err != nil {
			return NewNullVal(), err
		}
	}
	funcContext = false
	currentFunc = nil

	writeLnToFile("}\n")
	_, err := env.declareFunc(funcDeclaration.GetName(), fn)
	if err != nil {
		return result, fmt.Errorf("%s:%d:%d: %s", fileName, funcDeclaration.GetLn(), funcDeclaration.GetCol(), err)
	}
	return result, nil
}
