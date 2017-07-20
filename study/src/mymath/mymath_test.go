package mymath

import (
	"testing"
)

type mathTest struct {
	a, b, ret int
}

var addTest = []mathTest{
	{4, 6, 10},
	{5, 6, 11},
	{2, -6, -4},
}

var maxTest = []mathTest{
	{3, 5, 5},
	{-3, 5, 5},
	{-3, -5, -3},
}

func TestAdd(t *testing.T) {
	for _, v := range addTest {
		ret := Add(v.a, v.b)
		if ret != v.ret {
			t.Errorf("%d add %d, want %d, but get %d", v.a, v.b, v.ret, ret)
		}
	}
}
func TestMax(t *testing.T) {
	for _, v := range maxTest {
		ret := Max(v.a, v.b)
		if ret != v.ret {
			t.Errorf("the max number between %d and %d is want %d, but get %d", v.a, v.b, v.ret, ret)
		}
	}
}
