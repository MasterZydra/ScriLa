package runtime

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func nativePrintLn(args []IRuntimeVal, env *Environment) (IRuntimeVal, error) {
	nativePrint(args, env)
	fmt.Println("")
	return NewNullVal(), nil
}

func nativePrint(args []IRuntimeVal, env *Environment) (IRuntimeVal, error) {
	for _, arg := range args {
		fmt.Print(arg.ToString(), " ")
		arg.GetType()
	}
	fmt.Print("")
	return NewNullVal(), nil
}

func nativeInput(args []IRuntimeVal, env *Environment) (IRuntimeVal, error) {
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: input(str prompt)")
	}
	switch args[0].GetType() {
	case StrValueType:
		fmt.Print(args[0].ToString() + " ")
	default:
		return NewNullVal(), fmt.Errorf("nativeInput: Arg type '%s' not supported", args[0].GetType())
	}
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	input = strings.ReplaceAll(input, "\n", "")
	if err != nil {
		return NewNullVal(), err
	}

	return NewStrVal(input), nil
}
