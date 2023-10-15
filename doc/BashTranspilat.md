# Bash Transpilat

This document contains technical details how the structures have been solved in bash syntax.

**Content**
- [Bool](#bool)
- [Function Return values](#function-return-values)

## Bool
There are no native bool types in bash. So boolean expressions are represented as strings.

**Example:**  
ScriLa
```Python
if (true && false) {
	printLn("true");
}
```
Bash transpilat
```bash
if [[ "true" == "true" ]] && [[ "false" == "true" ]]; then
	echo "true"
fi
``` 

## Function Return values
Bash functions can return a status code between 0 and 255. Zero stands for success.  
Bash Example:  
```bash
function fn() {
    return 42
}
```

The bash transpilat for a ScriLa function uses a temporary variable to return the value:
| Return type | Variable |
| ----------- | -------- |
| bool        | `$tmpBool` |
| int         | `$tmpInt` |
| string      | `$tmpStr` |  

**Example:**  
ScriLa
```Python
func add(int a, int b) int {
	return a + b;
}

int result = add(13, 42);
```
Bash transpilat
```bash
add () {
	local a=$1
	local b=$2
	tmpInt=$((${a} + ${b}))
	return
}

add 13 42
result=${tmpInt}
``` 
