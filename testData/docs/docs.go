// Package docs exercises the documentation features of golang 1.19 and above at
// the package documentation level.
//
// # This is a heading
//
// This heading has a paragraph with a reference to the standard library
// [math/rand] as well as a function in the file [Func], a type [Type], a type's
// function [Type.Func], a non-standard library package
// [golang.org/x/crypto@v0.5.0/bcrypt.Cost], an external link [Outside Link] and
// a [broken link].
//
// It also has a numbered list:
//  1. First
//  2. Second
//  3. Third
//
// Plus one with blank lines:
//
//  1. First
//
//  2. Second
//
//  3. Third
//
// Non-numbered lists
//   - First
//     another line
//   - Second
//   - Third
//
// Plus blank lines:
//
//   - First
//
//     another paragraph
//
//   - Second
//
//   - Third
//
// And a golang code block:
//
//	func GolangCode(t int) int {
//		return t + 1
//	}
//
// And a random code block:
//
//	something
//		preformatted
//	in a random
//			way
//
// [Outside Link]: https://golang.org/doc/articles/json_and_go.html
package docs

// Func is present in this file.
func Func(param int) int {
	return param
}

// Type is a type in this file.
type Type struct{}

// TypeFunc is a func within a type in this file.
func (t *Type) Func() {}
