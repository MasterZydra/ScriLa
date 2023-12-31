# Bash Transpilat

This document contains technical details how the structures have been solved in bash syntax.

**Content**
- [Bool](#bool)
- [Bool - Assign comparison](#bool---assign-comparison)
- [Function Return values](#function-return-values)

## Bool
There are no native bool types in bash. So boolean expressions are represented as strings.

**Example:**  
```Python
# ScriLa
if (true && false) {
	printLn("true");
}
```
```bash
# Bash transpilat
if [[ "true" == "true" ]] && [[ "false" == "true" ]]; then
	echo "true"
fi
``` 

## Bool - Assign comparison
There is no native way to assign the result of a comparision as bool to a variable in bash. So this assignment will be replaced with a if statement that uses the temporary variable `$tmpBool` and sets it to `true` or `false` (in the else block).

**Example:**  

```Python
# ScriLa
bool b = 42 > 13;
```
```bash
# Bash transpilat
if [[ 42 -gt 13 ]]
then
	tmpBools[0]="true"
else
	tmpBools[0]="false"
fi
b="${tmpBools[0]}"
```


## Function Return values
Bash functions can return a status code between 0 and 255. Zero stands for success.  
Bash Example:  
```bash
# Bash
function fn() {
    return 42
}
```

The bash transpilat for a ScriLa function uses a temporary variable of type array to return the value:
| Return type | Variable |
| ----------- | -------- |
| bool        | `$tmpBools` |
| int         | `$tmpInts` |
| string      | `$tmpStrs` |  

The index of the array is dynamic so that it can be set from outside of the function. This allows passing multiple function calls as arguments for a function call.

**Example:**  
```Python
# ScriLa
func add(int a, int b) int {
	return a + b;
}

int result = add(13, 42);
```
```bash
# Bash transpilat
add () {
	local a=$1
	local b=$2
	tmpInts[${tmpIndex}]=$((${a} + ${b}))
	return
}

tmpIndex=0
add 13 42
result=${tmpInts[0]}

``` 
