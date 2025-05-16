package bitstrings

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func FromGOSTString(str string) *BitString128 {
	sep := strings.Split(str, " ")
	u := strings.Join(sep[:8], "")
	l := strings.Join(sep[8:], "")
	upper, _ := strconv.ParseUint(u, 16, 64)
	lower, _ := strconv.ParseUint(l, 16, 64)
	fmt.Println(upper, lower, math.MaxInt64)

	return NewBitString(uint64(upper), uint64(lower))
}
