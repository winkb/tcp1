package chain

import "testing"

func TestVersionCollect_reg(t *testing.T) {
	allString := versionReg.FindAllString("aaa", -1)
	t.Log(allString)
}
