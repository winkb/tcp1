package strfn

import (
	"github.com/winkb/tcp1/util/types"
	"strconv"
)

func ToInt[T types.Number](v string) T {
	i, _ := strconv.Atoi(v)
	return T(i)
}
