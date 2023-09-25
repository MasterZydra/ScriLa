package transpiler

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
	env.declareFunc("print", NewNativeFunc(nativePrint))
	env.declareFunc("time", NewNativeFunc(nativeTime))
}

type Environment struct {
	parent    *Environment
	functions map[string]IRuntimeVal
	variables map[string]IRuntimeVal
	constants []string
}

func NewEnvironment(parentEnv *Environment) *Environment {
	isGlobal := parentEnv == nil
	env := &Environment{
		parent:    parentEnv,
		functions: make(map[string]IRuntimeVal),
		variables: make(map[string]IRuntimeVal),
		constants: make([]string, 0),
	}

	if isGlobal {
		setupScope(env)
	}
	return env
}

func (self *Environment) declareFunc(funcName string, value IRuntimeVal) IRuntimeVal {
	if self.isFuncDeclared(funcName) {
		fmt.Println("Cannot declare function '" + funcName + "'. As it already is defined.")
		os.Exit(1)
	}

	self.functions[funcName] = value

	return value
}

func (self *Environment) isFuncDeclared(funcName string) bool {
	if _, ok := self.functions[funcName]; ok {
		return true
	}

	if self.parent == nil {
		return false
	}

	return self.parent.isFuncDeclared(funcName)
}

func (self *Environment) resolveFunc(funcName string) *Environment {
	if _, ok := self.functions[funcName]; ok {
		return self
	}

	if self.parent == nil {
		fmt.Println("Cannot resolve function '" + funcName + "' as it does not exist.")
		os.Exit(1)
	}

	return self.parent.resolveFunc(funcName)
}

func (self *Environment) lookupFunc(funcName string) IRuntimeVal {
	env := self.resolveFunc(funcName)
	return env.functions[funcName]
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
