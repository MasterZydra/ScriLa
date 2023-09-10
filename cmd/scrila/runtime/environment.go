package runtime

import (
	"fmt"
	"os"
)

type Environment struct {
	parent    *Environment
	variables map[string]IRuntimeVal
}

func NewEnvironment(parentEnv *Environment) *Environment {
	env := &Environment{
		parent:    parentEnv,
		variables: make(map[string]IRuntimeVal),
	}

	env.declareVar("x", NewIntVal(100))
	env.declareVar("null", NewNullVal())
	env.declareVar("true", NewBoolVal(true))
	env.declareVar("false", NewBoolVal(false))
	return env
}

func (self *Environment) declareVar(varName string, value IRuntimeVal) IRuntimeVal {
	if _, ok := self.variables[varName]; ok {
		fmt.Println("Cannot declare variable", varName, ". As it already is defined.")
		os.Exit(1)
	}

	self.variables[varName] = value
	return value
}

func (self *Environment) assignVar(varName string, value IRuntimeVal) IRuntimeVal {
	env := self.resolve(varName)
	env.variables[varName] = value
	return value
}

func (self *Environment) resolve(varName string) *Environment {
	if _, ok := self.variables[varName]; ok {
		return self
	}

	if self.parent == nil {
		fmt.Println("Cannot resolve variable '" + varName + "' as it does not exist.")
		os.Exit(1)
	}

	return self.parent.resolve(varName)
}

func (self *Environment) lookupVar(varName string) IRuntimeVal {
	env := self.resolve(varName)
	return env.variables[varName]
}
