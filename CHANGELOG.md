# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Added another temporary variabel `tmpIndex` to store the current array index so that the `tmpInts` array can be used without limitations
- Added support for arrays as function return type
- Added support for `for each` loop 
- Added native function `strSplit` #6

### Changed

- Changed the return value of the native function `exec` from `void` to `str` by returing the output of the given command

## v0.2.0-alpha

### Added

- Added support for declaring and assigning arrays #5
- Added support for assigning and adding values to arrays #5

### Changed

- Replaced temporary variables `tmpInt`, `tmpBool` and `tmpStr` with arrays of same type: `tmpInts`, `tmpBools`, `tmpStrs`
- Changed temporary variables array index from always 0 to a dynamic index. This allows multiple functions calls as arguments

### Fixed

- Fixed token types for braces and brackets so that they are not the same
- Fixed error caused if a function was declared inside of a function by preventing function nesting

## v0.1.2-alpha

### Added

- Added native function `exit`
- Added native function `strIsBool`
- Added native function `strToBool`

### Changed

- Changed release pipeline to zip the executables and name each of them "scrila"

### Fixed

- Fixed invalid Bash code if a binary comparision has a binary operation on one side

## v0.1.1-alpha

### Added

- Added output of the error message to tests that should pass and generate Bash code

### Changed

- The "show ..." flags are now stored inside the config
- Improved the output generated if the "show AST" flags are passed
- Moved logic to replace a binary comparison into an if-statement in helper function to remove code duplication

### Fixed

- Fixed invalid Bash code if a binary comparison was used as function argument
- Fixed invalid Bash code if a binary boolean operation was used as function argument or assignment
- Fixed transpilation for the "break" and "continue" keyword #4
- Fixed invalid Bash code if the body of if/while/function was just a comment

### Removed

- Removed now unused helper functions or functions e.g. on the IRuntimeVal

## v0.1.0-alpha

### Added

- Added support for assigning the result of a comparison to a bool variable
- Added helpers that allow to centralize the creation of the Bash AST. That means that the switch statements in the transpiler eval[…] functions could simplified

### Changed

- Move the code generation part from the transpiler into the bash assembler class
- Changed the runtime environment to not store explicit values. The main purpose is to keep track of the data type of a variable, whether it is a const and whether it is declared

### Removed

- Removed / commented out all examples and tests for objects. The reason is that the objects have a lof issues and they need a revision anyway

## v0.0.3-alpha

### Added

- Added check if return expression has a return value if the function has a return value
- Added flag to show call stack
- Added support for call expression in binary expression evaluation
- Store result of binary expression in result property to prevent multiple transpilations

### Fixed

- Fixed arguments so that the flags and filename can be passed
- Fixed runtime value for input function
- Fixed runtime value for user defined function

## v0.0.2-alpha

### Added

- Added check if type of return value and the function type match
- Added support for escaping double quotes, backslashes, etc. with a backslash within a string

### Changed

- Allow `return` expression without a value for functions with return type `void`