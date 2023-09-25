package transpiler

import (
	"ScriLa/cmd/scrila/parser"
)

func transpileTestMode(code string) {
	testMode = true
	parser := parser.NewParser()
	env := NewEnvironment(nil)

	program := parser.ProduceAST(code)
	Transpile(program, env, "")
}

func ExampleIntDeclaration() {
	transpileTestMode(`
		int i = 42;
		print(i);
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// i=42
	// echo "$i"
}

func ExampleIntAssignment() {
	transpileTestMode(`
		int i = 42;
		i = 101;
		print(i);
	`)

	// Output:
	// #!/bin/bash
	// # Created by Scrila Transpiler v0.0.1
	// i=42
	// i=101
	// echo "$i"
}

// TODO Rebuild code to work without os.exit
// func ExampleIntAssignmentWithMissingDeclaration() {
// 	transpileTestMode(`
// 		i = 42;
// 		print(i);
// 	`)

// 	// Output:
// 	// "
// }
