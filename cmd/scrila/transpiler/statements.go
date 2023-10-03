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
			writeLnToFile("\"" + varName + "\"")
			value = NewStrVal("")
		case lexer.IntType:
			writeLnToFile(varName)
			value = NewIntVal(1)
		default:
			return NewNullVal(), fmt.Errorf("evalVarDeclaration - CallExprNode: Unsupported varType '%s'", varDeclaration.GetVarType())
		}

	case ast.IdentifierNode:
		symbol := identNodeGetSymbol(varDeclaration.GetValue())
		if symbol == "null" {
			writeLnToFile("\"" + symbol + "\"")
		} else if slices.Contains(reservedIdentifiers, symbol) {
			writeLnToFile(symbol)
		} else {
			valueVarType, err := env.lookupVarType(identNodeGetSymbol(varDeclaration.GetValue()))
			if err != nil {
				return NewNullVal(), err
			}
			if valueVarType != varDeclaration.GetVarType() {
				return NewNullVal(), fmt.Errorf("Cannot assign a value of type '%s' to a var of type '%s'", valueVarType, varDeclaration.GetVarType())
			}
			switch varDeclaration.GetVarType() {
			case lexer.StrType:
				writeLnToFile("\"" + identNodeToBashVar(varDeclaration.GetValue()) + "\"")
			case lexer.IntType:
				writeLnToFile(identNodeToBashVar(varDeclaration.GetValue()))
			default:
				return NewNullVal(), fmt.Errorf("evalVarDeclaration - Identifier: Unsupported varType '%s'", varDeclaration.GetVarType())
			}
		}
	case ast.BinaryExprNode:
		switch varDeclaration.GetVarType() {
		case lexer.StrType:
			writeLnToFile("\"" + value.GetTranspilat() + "\"")
		case lexer.IntType:
			writeLnToFile(value.GetTranspilat())
		default:
			return NewNullVal(), fmt.Errorf("evalVarDeclaration - BinaryExpr: Unsupported varType '%s'", varDeclaration.GetVarType())
		}
	case ast.StrLiteralNode:
		if varDeclaration.GetVarType() != lexer.StrType {
			return NewNullVal(), fmt.Errorf("Cannot assign a value of type '%s' to a var of type '%s'", lexer.StrType, varDeclaration.GetVarType())
		}
		writeLnToFile("\"" + value.ToString() + "\"")
	case ast.IntLiteralNode:
		if varDeclaration.GetVarType() != lexer.IntType {
			return NewNullVal(), fmt.Errorf("Cannot assign a value of type '%s' to a var of type '%s'", lexer.IntType, varDeclaration.GetVarType())
		}
		writeLnToFile(value.ToString())
	case ast.ObjectLiteralNode:
		for _, prop := range ast.ExprToObjLit(varDeclaration.GetValue()).GetProperties() {
			writeToFile(varDeclaration.GetIdentifier() + "[\"" + prop.GetKey() + "\"]=")
			value, err := transpile(prop.GetValue(), env)
			if err != nil {
				return NewNullVal(), err
			}
			switch prop.GetValue().GetKind() {
			case ast.IntLiteralNode:
				writeLnToFile(value.ToString())
			case ast.StrLiteralNode:
				writeLnToFile("\"" + value.ToString() + "\"")
			case ast.IdentifierNode:
				symbol := identNodeGetSymbol(prop.GetValue())
				if symbol == "null" {
					writeLnToFile("\"" + symbol + "\"")
				} else if slices.Contains(reservedIdentifiers, symbol) {
					writeLnToFile(symbol)
				} else {
					writeLnToFile(identNodeToBashVar(prop.GetValue()))
				}
			default:
				return NewNullVal(), fmt.Errorf("evalVarDeclaration - ObjectLiteralNode: property kind '%s' not supported", prop.GetValue().GetKind())
			}
		}
	case ast.MemberExprNode:
		memberVal, err := evalMemberExpr(ast.ExprToMemberExpr(varDeclaration.GetValue()), env)
		if err != nil {
			return NewNullVal(), err
		}
		writeLnToFile(memberVal.GetTranspilat())
	default:
		return NewNullVal(), fmt.Errorf("evalVarDeclaration: value kind '%s' not supported", varDeclaration.GetValue().GetKind())
	}

	return env.declareVar(varDeclaration.GetIdentifier(), value, varDeclaration.IsConstant(), varDeclaration.GetVarType())
}

func evalFunctionDeclaration(funcDeclaration ast.IFunctionDeclaration, env *Environment) (IRuntimeVal, error) {
	fn := NewFunctionVal(funcDeclaration.GetName(), funcDeclaration.GetParameters(), env, funcDeclaration.GetBody())
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
			return NewNullVal(), fmt.Errorf("evalFunctionDeclaration: Param type '%s' not supported. (%s)", fn.GetParams()[i].GetParamType(), fn.GetParams()[i])
		}
		scope.declareVar(fn.GetParams()[i].GetName(), value, false, fn.GetParams()[i].GetParamType())
		writeLnToFile("\tlocal " + param.GetName() + "=$" + strconv.Itoa(i+1))
	}

	// Transpile the function body line by line
	funcContext = true
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

	writeLnToFile("}\n")
	_, err := env.declareFunc(funcDeclaration.GetName(), fn)
	if err != nil {
		return result, err
	}
	return result, nil
}
