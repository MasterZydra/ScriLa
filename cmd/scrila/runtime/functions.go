package runtime

import (
	"fmt"
)

func nativePrint(args []IRuntimeVal, env *Environment) IRuntimeVal {
	for _, arg := range args {
		fmt.Print(arg.ToString(), " ")
		arg.GetType()
	}
	fmt.Println("")
	return NewNullVal()
}
