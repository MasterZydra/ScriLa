# Bash Transpilat

This document contains technical details how the structures have been solved in bash syntax.

**Content**
- [Function Return values](#function-return-values)

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