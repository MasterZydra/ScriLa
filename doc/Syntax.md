# Syntax

**Content**  
- [Variables](#variables)
  - [Array variables](#array-variables)
  - [Boolean variables](#boolean-variables)
  - [Integer variables](#integer-variables)
  - [String variables](#string-variables)
- [Comparisons](#comparisons)
  - [Comparing Booleans](#comparing-booleans)
  - [Comparing Integers](#comparing-integers)
  - [Comparing Strings](#comparing-strings)
- [Control structures](#control-structures)
  - [If](#if)
  - [While](#while)
- [Native functions](#native-functions)
  - [Exec](#exec)
  - [Exit](#exit)
  - [Input](#input)
  - [Print](#print)
  - [Sleep](#sleep)
  - [StrIsBool](#strisbool)
  - [StrIsInt](#strisint)
  - [StrToBool](#strtobool)
  - [StrToInt](#strtoint)
- [User defined functions](#user-defined-functions)
  - [Without parameters](#without-parameters)
  - [With parameters](#with-parameters)
  - [With return value](#with-return-value)

# Variables
A variable can store a specified type of value e.g. `int`, `string`, `bool`. This type cannot be changed later in the program.

[Back to top](#syntax)

## Array variables
An array can be declared from every base data type e.g. `int`, `string`, `bool`.

**Syntax**  
```Python
dataType[] variableName = value;
```

**Example**  
```Python
str[] strs = [];       # Empty string array
int[] ints = [43, 44]; # Integer array with values
ints[0] = 42;          # Change value at index 0 to 42
ints[] = 45;           # Append array with new value 45
```

[Back to top](#syntax)

## Boolean variables
A boolean variable can only store the values `true` and `false`.

**Example**  
```Python
bool b = false;
b = true;
```

[Back to top](#syntax)

## Integer variables
An integer variable can store values between -2^63 and 2^63-1 on 64-bit computers.

**Example**  
```Python
int i = 42;
i = 48 / 2;
```

[Back to top](#syntax)

## String variables
A string variable can store a string value. The limit of long a string can be depends on the environment where the bash script will be executed. 

**Example**
```Python
str s = "Hello ";
s += "World";
```

[Back to top](#syntax)

# Comparisons
Two values can be compared with a boolean as return value. This way the result of the comparison can be used in a condition of an `if` or `while` control structure.

[Back to top](#syntax)

## Comparing Booleans
The comparison of boolean values allows the equal and unequal operation.

**Example**  
```Python
bool b = true;

# Equal
if (b == false) {
}
# Unequal
if (b != true) {
}
```

[Back to top](#syntax)

## Comparing Integers
The comparison of integer values allows the equal, unequal, greater (or equal) and smaller (or equal) operation.

**Example**  
```Python
int i = 42;
# Equal
if (i == 48) {
}
# Unequal
if (i != 48) {
}
# Smaller
if (i < 48) {
}
# Smaller or equal
if (i <= 48) {
}
# Greater
if (i > 48) {
}
# Greater or equal
if (i >= 48) {
}
```

[Back to top](#syntax)

## Comparing Strings
The comparison of string values allows the equal, unequal, greater and smaller operation.

**Example**  
```Python
str s = "abc";
# Equal
if (s == "abc") {
}
# Unequal
if (s != "abc") {
}
# Smaller
if (s < "bcd") {
}
# Greater
if (s > "bcd") {
}
```

[Back to top](#syntax)

# Control structures
Control structures allow to change a purely linear program flow e.g. to "branches" (`if`) depending on a condition or repeat a code block multiple times until a condition is true (`while`).

[Back to top](#syntax)

## If
Use the `if` statement to execute a specified block of code if a condition is `true`.

**Syntax**  
```Python
if (condition) {
    # block of code that is executed if condition is true
}
```

**Examples**  
```Python
if (true) {
    printLn("true");
} else if (true && true) {
    printLn("true");
}

if (42 > 13) {
    printLn("true");
} else {
    printLn("false");
}

bool b = true;
if (b) {
    printLn("true");
}
```

[Back to top](#syntax)

## While
The `while` loop executes the block of code until the given condition is `true`.

**Syntax**  
```Python
while (condition) {
    # block of code that is executed repeatedly while the condition is true
}
```

**Example**  
```Python
while (true && true) {
    printLn("true");
}

int i = 0;
while (i < 10) {
    i += 1;
}

bool b = true;
while (b) {
    printLn("true");
}
```

[Back to top](#syntax)

# Native functions
The following functions are provided for use in your scripts. These are directly transpiled into bash or represented by a function implemented in bash.

[Back to top](#syntax)

## Exec
The native function `exec` allows to directly add bash code into the transpilat. 

**Syntax**  
```Python
exec(str command) void
```

**Example**  
```Python
exec("touch test.txt");

str cmd = "touch test.txt";
exec(cmd);
```

[Back to top](#syntax)

## Exit
The native function `exit` exits the current script with a status code.

**Syntax**  
```Python
exit(int code) void
```

**Example**  
```Python
exit(0); # Success
exit(1); # Error
```

[Back to top](#syntax)

## Input
The native function `input` waits for the user of the script to input a string and returns it. 

**Syntax**  
```Python
input(str prompt) str
```

**Example**  
```Python
str userInput = input("Pleas enter your name:");
```

[Back to top](#syntax)

## Print
The native functions `print` and `printLn` write the given values to terminal. The difference between `print` and `printLn` is that `printLn` adds new line.

**Syntax**  
```Python
print(str|int|bool value, ...) void
printLn(str|int|bool value, ...) void
```

**Example**  
```Python
print("Hello ");
printLn("World");

printLn("str", 42, false, null);

int i = 43;
printLn("i =", i);
```

[Back to top](#syntax)

## Sleep
The native function `sleep` waits for the given amount of seconds. After that the program flow continues.

**Syntax**  
```Python
sleep(int seconds) void
```

**Example**  
```Python
# Sleep for 10 seconds
sleep(10);
```

[Back to top](#syntax)

## StrIsBool
The native function `strIsBool` checks if the given string is a boolean so that it could be converted to a string e.g. with `strToBool`.

**Syntax**  
```Python
strIsBool(str value) bool
```

**Example**  
```Python
strIsBool("true"); # true
strIsBool("str"); # false
```

## StrIsInt
The native function `strIsInt` checks if the given string is a number so that it could be converted to a string e.g. with `strToInt`.

**Syntax**  
```Python
strIsInt(str value) bool
```

**Example**  
```Python
strIsInt("123"); # true
strIsInt("str"); # false
```

[Back to top](#syntax)

## StrToBool
The native function `strToBool` takes a given string and tries to convert it into a boolean value.

**Syntax**  
```Python
strToBool(str value) bool
```

**Example**  
```Python
str s = "123";
bool b1 = strToBool(s); # b = false
bool b2 = strTobool("true") # b = true
bool b3 = strTobool("false") # b = false
```

[Back to top](#syntax)

## StrToInt
The native function `strToInt` takes a given string and tries to convert it into an integer value.

**Syntax**  
```Python
strToInt(str value) int
```

**Example**  
```Python
str s = "123";
int i = strToInt(s); # i = 123
```

[Back to top](#syntax)

# User defined functions
A function can be used to reuse code and make it easier to read.

**Syntax**  
```Python
func functionName(type param) type {
    # block of code executed when the function is called
}
```

[Back to top](#syntax)

## Without parameters
A function can be defined without parameters.

**Example**  
```Python
func printHelloWorld() void {
    printLn("Hello World!");
}

printHelloWorld();
```

[Back to top](#syntax)

## With parameters
The parameters must have a fixed type e.g. `int`, `str`, `bool`.

**Example**  
```Python
func printGivenString(str s) void {
    printLn("Given string: '" + s + "'");
}

printGivenString("Hello World");
```

[Back to top](#syntax)

## With return value
A function can return a value. The type of the value must be fixed e.g. `int`, `str`, `bool`.

**Example**  
```Python
func add(int a, int b) int {
    return a + b;
}

int c = add(123, 456);
```

[Back to top](#syntax)
