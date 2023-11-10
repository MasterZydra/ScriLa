# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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