package ast

import "golang.org/x/exp/slices"

var bools = []string{"true", "false"}

func IdentIsBool(ident IIdentifier) bool {
	return slices.Contains(bools, ident.GetSymbol())
}
