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
		writeLnToFile("$?")
	case ast.IdentifierNode:
		var i interface{} = varDeclaration.GetValue()
		identifier, _ := i.(ast.IIdentifier)
		if slices.Contains(reservedIdentifiers, identifier.GetSymbol()) {
			writeLnToFile(identifier.GetSymbol())
		} else {
			writeLnToFile("${" + identifier.GetSymbol() + "}")
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
		writeLnToFile("\"" + value.ToString() + "\"")
	case ast.IntLiteralNode:
		writeLnToFile(value.ToString())
	case ast.ObjectLiteralNode:
		var i interface{} = varDeclaration.GetValue()
		objectLiteral, _ := i.(ast.IObjectLiteral)
		for _, prop := range objectLiteral.GetProperties() {
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
				i = prop.GetValue()
				identifier, _ := i.(ast.IIdentifier)
				if identifier.GetSymbol() == "null" {
					writeLnToFile("\"" + identifier.GetSymbol() + "\"")
				} else if slices.Contains(reservedIdentifiers, identifier.GetSymbol()) {
					writeLnToFile(identifier.GetSymbol())
				} else {
					writeLnToFile("$" + identifier.GetSymbol())
				}
			default:
				return NewNullVal(), fmt.Errorf("evalVarDeclaration - ObjectLiteralNode: property kind '%s' not supported", prop.GetValue().GetKind())
			}
		}
	case ast.MemberExprNode:
		var i interface{} = varDeclaration.GetValue()
		memberExpr, _ := i.(ast.IMemberExpr)
		memberVal, err := evalMemberExpr(memberExpr, env)
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
		// TODO Check the bounds here. Verify arity of function.
		// Which means: len(fn.GetParams()) == len(args)
		// TODO var type - Get from function declaration and validate type against given type
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
