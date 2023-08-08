package numfn

import (
	"fmt"
	"strconv"
	"tcp1/util/types"
)

func ToStr[T types.Number](v T) string {
	if float64(int(v)) != float64(v) {
		return fmt.Sprintf("%v", v)
	}

	return strconv.Itoa(int(v))
}
