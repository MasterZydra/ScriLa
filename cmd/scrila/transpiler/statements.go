package transpiler

import (
	"ScriLa/cmd/scrila/ast"
)

func evalProgram(program ast.IProgram, env *Environment) IRuntimeVal {
	var lastEvaluated IRuntimeVal = NewNullVal()

	for _, statement := range program.GetBody() {
		lastEvaluated = transpile(statement, env)
	}

	return lastEvaluated
}

func evalVarDeclaration(varDeclaration ast.IVarDeclaration, env *Environment) IRuntimeVal {
	writeToFile(varDeclaration.GetIdentifier() + "=")
	value := transpile(varDeclaration.GetValue(), env)
	if value.GetType() == StrValueType {
		writeLnToFile("\"" + value.ToString() + "\"")
	} else {
		writeLnToFile(value.ToString())
	}
	return env.declareVar(varDeclaration.GetIdentifier(), value, varDeclaration.IsConstant())
}

func evalFunctionDeclaration(funcDeclaration ast.IFunctionDeclaration, env *Environment) IRuntimeVal {
	fn := NewFunctionVal(funcDeclaration.GetName(), funcDeclaration.GetParameters(), env, funcDeclaration.GetBody())

	return env.declareFunc(funcDeclaration.GetName(), fn)
}
