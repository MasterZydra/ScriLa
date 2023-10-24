# Syntax

**Content**
- [Variables](#variables)
- [Comparisons](#comparisons)
- [If](#if)
- [While](#while)
- [Native functions](#native-functions)
- [User defined functions](#user-defined-functions)

## Variables
**Integer**
```Python
int i = 42;
i = 48 / 2;
```

**String**
```Python
str s = "Hello ";
s += "World";
```

**Bool**
```Python
bool b = false;
b = true;
```

**Object**
```Python
obj o = {
    a: 1,
    b: null,
    c: false,
};
o.a = 2;
o.d = "d";
```

## Comparisons
**Integer**
```Python
int i = 42;
if (i == 48) {
    # Do smth
}
if (i != 48) {
    # Do smth
}
if (i <= 48) {
    # Do smth
}
if (i >= 48) {
    # Do smth
}
if (i < 48) {
    # Do smth
}
if (i > 48) {
    # Do smth
}
```

## If
```Python
if (true) {
    printLn("true");
} else if (true && true) {
    printLn("true");
}
if (true || false) {
    printLn("true");
} else {
    printLn("false");
}
bool b = true;
if (b) {
    printLn("true");
}
```

## While
```Python
while (true && true) {
    printLn("true");
}
while (true || false) {
    printLn("true");
}
bool b = true;
while (b) {
    printLn("true");
}
```

## Native functions
**Exec**  
```Python
exec("touch test.txt");

str cmd = "touch test.txt";
exec(cmd);
```

**Print**  
```Python
print("Hello ");
printLn("World");

printLn("str", 42, false, null);

int i = 43;
printLn("i =", i);
```

**Input**  
```Python
str userInput = input("Pleas enter your name:");
```

**Sleep**
```Python
# Sleep for 10 seconds
sleep(10);
```

**StrIsInt**  
```Python
# Check if the given input is an integer
strIsInt(123); # true
strIsInt("123"); # true
strIsInt("str"); # false
```

**StrToInt**  
```Python
str s = "123";
int i = strToInt(s);
```

## User defined functions
**Without parameters**
```Python
func printHelloWorld() void {
    printLn("Hello World!");
}

printHelloWorld();
```

**With parameters**
```Python
func printGivenString(str s) void {
    printLn("Given string: '" + s + "'");
}

printGivenString("Hello World");
```

**With return value**
```Python
func add(int a, int b) int {
    return a + b;
}

int c = add(123, 456);
```
