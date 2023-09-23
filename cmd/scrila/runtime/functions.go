package runtime

import (
	"fmt"
	"time"
)

func nativeTime(args []IRuntimeVal, env *Environment) IRuntimeVal {
	return NewIntVal(time.Now().UnixMilli())
}

func nativePrint(args []IRuntimeVal, env *Environment) IRuntimeVal {
	for _, arg := range args {
		fmt.Print(arg.ToString(), " ")
		arg.GetType()
	}
	fmt.Println("")
	return NewNullVal()
}
