package transpiler

import (
	"ScriLa/cmd/scrila/ast"
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
	if value.GetType() == StrValueType {
		writeLnToFile("\"" + value.ToString() + "\"")
	} else {
		writeLnToFile(value.ToString())
	}
	return env.declareVar(varDeclaration.GetIdentifier(), value, varDeclaration.IsConstant())
}

func evalFunctionDeclaration(funcDeclaration ast.IFunctionDeclaration, env *Environment) (IRuntimeVal, error) {
	fn := NewFunctionVal(funcDeclaration.GetName(), funcDeclaration.GetParameters(), env, funcDeclaration.GetBody())

	return env.declareFunc(funcDeclaration.GetName(), fn)
}
