# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Added output of the error message to tests that should pass and generate Bash code

### Changed

- The "show ..." flags are now stored inside the config
- Improved the output generated if the "show AST" flags are passed
- Moved logic to replace a binary comparison into an if-statement in helper function to remove code duplication

### Fixed

- Fixed invalid Bash code if a binary comparison was used as function argument
- Fixed invalid bash code if a binary boolean operation was used as function argument or assignment

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