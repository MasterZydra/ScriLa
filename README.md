# ScriLa

ScriLa is a scripting language that adopts a syntax reminiscent of Go and C++, and it compiles down to Bash.  
The primary objective is to craft ScriLa in a manner that ensures intuitive usage. Furthermore, it aims to provide type safety similar to TypeScript's relationship with JavaScript.

**Example**  
```Python
func calculateAge(int birthYear) int {
    return 2023 - birthYear;
}

str yearStr = input("Please enter your birth year:");

if (strIsInt(yearStr)) {
    int birthYear = strToInt(yearStr);
    printLn("Your age is", calculateAge(birthYear));
} else {
    printLn("The input '" + yearStr + "' was not a number");
}
```

**Usage**  
To transpile a script written in ScriLa into a bash script, execute the following command:  
`> scrila -f myFileName.scri`  

The bash file will be named `myFileName.scri.sh` and will be placed in the same folder as the passed file.

**Creation**  
The foundation of this project was created by following the playlist [Build a Custom Scripting Language In Typescript](https://www.youtube.com/playlist?list=PL_2VhOvlMk4UHGqYCLWc6GO8FaPl8fQTh) by [tylerlaceby](https://www.youtube.com/@tylerlaceby). Afterwards it was extended with types, control structures and transpilation.

**Documentation**  
- [Syntax](doc/Syntax.md)
- [Development](doc/Development.md)
- [Steps](doc/Steps.md)
- [AST](doc/AST.md)
- [Parser](doc/Parser.md)
- [Bash transpilat](doc/BashTranspilat.md)

**Possible features**  
- Transpile into PowerShell
