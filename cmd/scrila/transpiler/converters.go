package transpiler

import "ScriLa/cmd/scrila/ast"

func runtimeToBoolVal(runtimeVal ast.IRuntimeVal) IBoolVal {
	var i interface{} = runtimeVal
	return i.(IBoolVal)
}

func runtimeToIntVal(runtimeVal ast.IRuntimeVal) IIntVal {
	var i interface{} = runtimeVal
	return i.(IIntVal)
}

func runtimeToStrVal(runtimeVal ast.IRuntimeVal) IStrVal {
	var i interface{} = runtimeVal
	return i.(IStrVal)
}

func runtimeToObjVal(runtimeVal ast.IRuntimeVal) IObjVal {
	var i interface{} = runtimeVal
	return i.(IObjVal)
}

func runtimeToNativeFunc(runtimeVal ast.IRuntimeVal) INativeFunc {
	var i interface{} = runtimeVal
	return i.(INativeFunc)
}

func runtimeToFuncVal(runtimeVal ast.IRuntimeVal) IFunctionVal {
	var i interface{} = runtimeVal
	return i.(IFunctionVal)
}
