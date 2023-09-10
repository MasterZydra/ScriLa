package runtime

import (
	"ScriLa/cmd/scrila/ast"
)

func evalProgram(program ast.IProgram, env *Environment) IRuntimeVal {
	var lastEvaluated IRuntimeVal = NewNullVal()

	for _, statement := range program.GetBody() {
		lastEvaluated = Evaluate(statement, env)
	}

	return lastEvaluated
}

func evalVarDeclaration(varDeclaration ast.IVarDeclaration, env *Environment) IRuntimeVal {
	value := Evaluate(varDeclaration.GetValue(), env)
	return env.declareVar(varDeclaration.GetIdentifier(), value, varDeclaration.IsConstant())
}
