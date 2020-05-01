package main

import (
	"fmt"
	"testing"
)

func TestDiffStringSet(t *testing.T) {
	type testcase struct {
		old, new, add, remove map[string]struct{}
	}

	testcases := []testcase{
		{m(), m(), m(), m()},
		{m(), m("foo"), m("foo"), m()},
		{m("foo", "bar"), m("bar"), m(), m("foo")},
		{m("foo", "bar"), m("bar", "baz", "qux"), m("baz", "qux"), m("foo")},
	}

	for i, tt := range testcases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			gota, gotr := diffStringSet(tt.old, tt.new)
			if !equalSet(tt.add, gota) {
				t.Errorf("add sets unequal: expected=%v got=%v", tt.add, gota)
				return
			}
			if !equalSet(tt.remove, gotr) {
				t.Errorf("remove sets unequal: expected=%v got=%v", tt.remove, gotr)
				return
			}
		})
	}
}

func m(v ...string) map[string]struct{} {
	ret := make(map[string]struct{})
	for _, vv := range v {
		ret[vv] = struct{}{}
	}
	return ret
}

func equalSet(a, b map[string]struct{}) bool {
	if len(a) != len(b) {
		return false
	}

	for k := range a {
		if _, ok := b[k]; !ok {
			return false
		}
	}

	return true
}
