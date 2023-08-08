package numfn

import (
	"testing"
)

func TestToStr(t *testing.T) {
	var s = ToStr(1)
	t.Logf("%T", s)
	t.Log(s)
	t.Log(ToStr(1.1))
}
