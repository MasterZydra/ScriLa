package transpiler

func runtimeToIntVal(runtimeVal IRuntimeVal) IIntVal {
	var i interface{} = runtimeVal
	return i.(IIntVal)
}

func runtimeToStrVal(runtimeVal IRuntimeVal) IStrVal {
	var i interface{} = runtimeVal
	return i.(IStrVal)
}

func runtimeToObjVal(runtimeVal IRuntimeVal) IObjVal {
	var i interface{} = runtimeVal
	return i.(IObjVal)
}

func runtimeToNativeFunc(runtimeVal IRuntimeVal) INativeFunc {
	var i interface{} = runtimeVal
	return i.(INativeFunc)
}

func runtimeToFuncVal(runtimeVal IRuntimeVal) IFunctionVal {
	var i interface{} = runtimeVal
	return i.(IFunctionVal)
}
