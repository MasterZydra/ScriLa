package transpiler

import (
	"ScriLa/cmd/scrila/ast"
	"ScriLa/cmd/scrila/lexer"
	"fmt"

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
	writeToFile(varDeclaration.GetIdentifier() + "=")
	value, err := transpile(varDeclaration.GetValue(), env)
	if err != nil {
		return NewNullVal(), err
	}
	switch varDeclaration.GetValue().GetKind() {
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
	default:
		return NewNullVal(), fmt.Errorf("evalVarDeclaration: value kind '%s' not supported", varDeclaration.GetValue())
	}

	return env.declareVar(varDeclaration.GetIdentifier(), value, varDeclaration.IsConstant(), varDeclaration.GetVarType())
}

func evalFunctionDeclaration(funcDeclaration ast.IFunctionDeclaration, env *Environment) (IRuntimeVal, error) {
	fn := NewFunctionVal(funcDeclaration.GetName(), funcDeclaration.GetParameters(), env, funcDeclaration.GetBody())

	return env.declareFunc(funcDeclaration.GetName(), fn)
}
