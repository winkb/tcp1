package chain

import "testing"

func TestVersionCollect_reg(t *testing.T) {
	allString := versionReg.FindStringSubmatch("v1.0.1")
	t.Log(allString[1])
}
