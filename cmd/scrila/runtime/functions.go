package runtime

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func defineNativeFunctions(env *Environment) {
	var nativeFunctions = map[string]FunctionCall{
		"input":   nativeInput,
		"print":   nativePrint,
		"printLn": nativePrintLn,
		"sleep":   nativeSleep,
	}

	for name, function := range nativeFunctions {
		env.declareFunc(name, NewNativeFunc(function))
	}
}

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

func nativeSleep(args []IRuntimeVal, env *Environment) (IRuntimeVal, error) {
	if len(args) != 1 {
		return NewNullVal(), fmt.Errorf("Expected syntax: sleep(int seconds)")
	}
	switch args[0].GetType() {
	case IntValueType:
		time.Sleep(time.Duration(runtimetoIntVal(args[0]).GetValue()) * time.Second)
	default:
		return NewNullVal(), fmt.Errorf("nativeSleep: Arg type '%s' not supported", args[0].GetType())
	}
	return NewNullVal(), nil
}
