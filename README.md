# ScriLa

ScriLa is a scripting language with a Go and C++ like syntax that transpiles into Bash and PowerShell.

**Content**
- [Syntax](#syntax)
  - [Variables](#variables)
- [Development](#development)

# Syntax
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

# Development
**Compile and run**  
`go run ./...`

**Build executable**  
`go build -o . ./...`

**Resources**  
- Playlist [Build a Custom Scripting Language In Typescript](https://www.youtube.com/playlist?list=PL_2VhOvlMk4UHGqYCLWc6GO8FaPl8fQTh) by [tylerlaceby](https://www.youtube.com/@tylerlaceby)