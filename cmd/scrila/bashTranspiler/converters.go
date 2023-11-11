package bashTranspiler

import (
	"ScriLa/cmd/scrila/scrilaAst"
)

func runtimeToBoolVal(runtimeVal scrilaAst.IRuntimeVal) IBoolVal {
	var i interface{} = runtimeVal
	return i.(IBoolVal)
}

func runtimeToIntVal(runtimeVal scrilaAst.IRuntimeVal) IIntVal {
	var i interface{} = runtimeVal
	return i.(IIntVal)
}

func runtimeToStrVal(runtimeVal scrilaAst.IRuntimeVal) IStrVal {
	var i interface{} = runtimeVal
	return i.(IStrVal)
}

func runtimeToObjVal(runtimeVal scrilaAst.IRuntimeVal) IObjVal {
	var i interface{} = runtimeVal
	return i.(IObjVal)
}

func runtimeToNativeFunc(runtimeVal scrilaAst.IRuntimeVal) INativeFunc {
	var i interface{} = runtimeVal
	return i.(INativeFunc)
}

func runtimeToFuncVal(runtimeVal scrilaAst.IRuntimeVal) IFunctionVal {
	var i interface{} = runtimeVal
	return i.(IFunctionVal)
}
