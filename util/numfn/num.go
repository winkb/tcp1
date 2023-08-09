package numfn

import (
	"fmt"
	"github.com/winkb/tcp1/util/types"
	"strconv"
)

func ToStr[T types.Number](v T) string {
	if float64(int(v)) != float64(v) {
		return fmt.Sprintf("%v", v)
	}

	return strconv.Itoa(int(v))
}
