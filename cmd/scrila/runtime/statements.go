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

func evalFunctionDeclaration(funcDeclaration ast.IFunctionDeclaration, env *Environment) IRuntimeVal {
	fn := NewFunctionVal(funcDeclaration.GetName(), funcDeclaration.GetParameters(), env, funcDeclaration.GetBody())

	return env.declareVar(funcDeclaration.GetName(), fn, true)
}
