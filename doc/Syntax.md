# Syntax

**Content**
- [Variables](#variables)
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

## Native functions
**Print**  
```
print("Hello ");
printLn("World);

printLn("str", 42, false, null);

int i = 43;
printLn("i =", i);
```

## User defined functions
**Without parameters**
```
func printHelloWorld() {
    printLn("Hello World!");
}

printHelloWorld();
```

**With parameters**
```
func printGivenString(str s) {
    printLn("Given string: '" + s + "'");
}

printGivenString("Hello World");
```

**With return value**
```
func add(int a, int b) {
    return a + b;
}

int c = add(123, 456);
```