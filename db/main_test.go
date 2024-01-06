package main

import (
	"testing"
)

func Test_getRandStr(t *testing.T) {
	g1 := getRandStr()

	t.Log(g1)
}

func BenchmarkGetRandStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.Log(getRandStr())
	}
}
