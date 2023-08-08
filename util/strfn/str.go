package strfn

import (
	"strconv"
	"tcp1/util/types"
)

func ToInt[T types.Number](v string) T {
	i, _ := strconv.Atoi(v)
	return T(i)
}
