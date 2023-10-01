package runtime

import (
	"fmt"
)

func nativePrintLn(args []IRuntimeVal, env *Environment) IRuntimeVal {
	nativePrint(args, env)
	fmt.Println("")
	return NewNullVal()
}

func nativePrint(args []IRuntimeVal, env *Environment) IRuntimeVal {
	for _, arg := range args {
		fmt.Print(arg.ToString(), " ")
		arg.GetType()
	}
	fmt.Print("")
	return NewNullVal()
}
