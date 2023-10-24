# ScriLa

ScriLa is a scripting language with a Go and C++ like syntax that transpiles into Bash and PowerShell.

**Documentation**  
- [Syntax](doc/Syntax.md)
- [AST](doc/AST.md)
- [Development](doc/Development.md)
- [Parser](doc/Parser.md)
- [Bash transpilat](doc/BashTranspilat.md)

## Usage
To transpile a script written in ScriLa into bash execute the following command:  
`> scrila myFileName.scri`  
The bash file will be named `myFileName.scri.sh` and will be placed in the same folder as the passed file.

## Resources
- Playlist [Build a Custom Scripting Language In Typescript](https://www.youtube.com/playlist?list=PL_2VhOvlMk4UHGqYCLWc6GO8FaPl8fQTh) by [tylerlaceby](https://www.youtube.com/@tylerlaceby)
