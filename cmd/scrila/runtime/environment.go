package runtime

import (
	"fmt"
	"os"

	"golang.org/x/exp/slices"
)

func setupScope(env *Environment) {
	// Create Default Global Environment
	env.declareVar("null", NewNullVal(), true)
	env.declareVar("true", NewBoolVal(true), true)
	env.declareVar("false", NewBoolVal(false), true)

	// Define native builtin methods
	env.declareVar("print", NewNativeFunc(nativePrint), true)
	env.declareVar("time", NewNativeFunc(nativeTime), true)
}

type Environment struct {
	parent    *Environment
	variables map[string]IRuntimeVal
	constants []string
}

func NewEnvironment(parentEnv *Environment) *Environment {
	isGlobal := parentEnv == nil
	env := &Environment{
		parent:    parentEnv,
		variables: make(map[string]IRuntimeVal),
		constants: make([]string, 0),
	}

	if isGlobal {
		setupScope(env)
	}
	return env
}

func (self *Environment) declareVar(varName string, value IRuntimeVal, isConstant bool) IRuntimeVal {
	if _, ok := self.variables[varName]; ok {
		fmt.Println("Cannot declare variable '" + varName + "'. As it already is defined.")
		os.Exit(1)
	}

	self.variables[varName] = value

	if isConstant {
		self.constants = append(self.constants, varName)
	}
	return value
}

func (self *Environment) assignVar(varName string, value IRuntimeVal) IRuntimeVal {
	env := self.resolve(varName)

	// Cannot assign to constant
	if slices.Contains(self.constants, varName) {
		fmt.Println("Cannot reasign to variable '" + varName + "' as it was declared constant.")
		os.Exit(1)
	}

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
