package chain

import (
	"fmt"
	"regexp"
)

// 1.0.0
var versionReg = regexp.MustCompile(`^v?(\d)(\.\d){0,2}$`)

type VersionCollect struct {
}

func NewVersionCollect() *VersionCollect {
	return &VersionCollect{}
}

func (l *VersionCollect) Title(c *Collect) {
	fmt.Println(fmt.Sprintf("当前版本%s:请输入版本号:", c.data.Version))
}

func (l *VersionCollect) Collect(txt string, data *InputData) bool {
	allString := versionReg.FindAllString(txt, -1)
	if len(allString) == 0 {
		return false
	}

	version := allString[0]

	data.Version = version

	return true
}
