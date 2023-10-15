# Syntax

**Content**
- [Variables](#variables)
- [Comparisons](#comparisons)
- [If](#if)
- [Native functions](#native-functions)
- [User defined functions](#user-defined-functions)

## Variables
**Integer**
```
int i = 42;
i = 48 / 2;
```

**String**
```
str s = "Hello ";
s += "World";
```

**Bool**
```
bool b = false;
b = true;
```

**Object**
```
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
```
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
```
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

## Native functions
**Print**  
```
print("Hello ");
printLn("World");

printLn("str", 42, false, null);

int i = 43;
printLn("i =", i);
```

**Input**  
```
str userInput = input("Pleas enter your name:");
```

**Sleep**
```
# Sleep for 10 seconds
sleep(10);
```

**IsInt**  
```
# Check if the given input is an integer
isInt(123); # true
isInt("123"); # true
isInt("str"); # false
```

**StrToInt**  
```
str s = "123";
int i = strToInt(s);
```

## User defined functions
**Without parameters**
```
func printHelloWorld() void {
    printLn("Hello World!");
}

printHelloWorld();
```

**With parameters**
```
func printGivenString(str s) void {
    printLn("Given string: '" + s + "'");
}

printGivenString("Hello World");
```

**With return value**
```
func add(int a, int b) int {
    return a + b;
}

int c = add(123, 456);
```