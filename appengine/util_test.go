package main

import (
	"fmt"
	"testing"
)

func TestDiffStringMaps(t *testing.T) {
	type testcase struct {
		old    map[string]struct{}
		new    map[string]string
		add    map[string]string
		remove map[string]struct{}
	}

	testcases := []testcase{
		{n(), m(), m(), n()},
		{n(), m("foo"), m("foo"), n()},
		{n("foo", "bar"), m("bar"), m(), n("foo")},
		{n("foo", "bar"), m("bar", "baz", "qux"), m("baz", "qux"), n("foo")},
	}

	for i, tt := range testcases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			gota, gotr := diffStringMaps(tt.old, tt.new)
			if !equalKeys(tt.add, gota) {
				t.Errorf("add maps unequal: expected=%v got=%v", tt.add, gota)
				return
			}
			if !equalKeys2(tt.remove, gotr) {
				t.Errorf("remove maps unequal: expected=%v got=%v", tt.remove, gotr)
				return
			}
		})
	}
}

func m(v ...string) map[string]string {
	ret := make(map[string]string)
	for _, vv := range v {
		ret[vv] = ""
	}
	return ret
}

func n(v ...string) map[string]struct{} {
	ret := make(map[string]struct{})
	for _, vv := range v {
		ret[vv] = struct{}{}
	}
	return ret
}

func equalKeys(a, b map[string]string) bool {
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

func equalKeys2(a, b map[string]struct{}) bool {
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
