package util

import (
	"go/types"
	"iter"
	"slices"

	"golang.org/x/tools/go/types/typeutil"
)

type TypeSet struct {
	m typeutil.Map
}

func NewTypeSet() *TypeSet {
	return &TypeSet{}
}

func (t *TypeSet) Size() int {
	return t.m.Len()
}

func (t *TypeSet) Add(typ types.Type) {
	t.m.Set(typ, struct{}{})
}

func (t *TypeSet) Remove(typ types.Type) bool {
	return t.m.Delete(typ)
}

func (t *TypeSet) Contains(typ types.Type) bool {
	return t.m.At(typ) != nil
}

func (t *TypeSet) Values() iter.Seq[types.Type] {
	keys := t.m.Keys()
	return slices.Values(keys)
}

// TypeString returns the string representation of t, including package prefixes.
func TypeString(t types.Type) string {
	if t == nil {
		panic("t nil")
	}
	// this is not the same as t.String(). consider t = *types.Basic,
	// then t.String() is *go/types.Basic, but TypeString(t) is *types.Basic
	// (syntactically correct only when outside the package types, see TypeStringRelative)
	return types.TypeString(t, func(p *types.Package) string {
		return p.Name()
	})
}

// TypeStringRelative returns the string representation of t relative to the package defined
// defined by the path inPkgPath, i.e. if t originates from the package defined by inPkgPath,
// then the package prefix is omitted.
func TypeStringRelative(t types.Type, inPkgPath string) string {
	if t == nil {
		panic("t nil")
	}
	return types.TypeString(t, func(p *types.Package) string {
		if p.Path() == inPkgPath {
			return ""
		}
		return p.Name()
	})
}
