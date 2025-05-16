package bitstrings

import (
	"fmt"
	"strconv"
	"strings"
)

func FromGOSTString(str string) *BitString128 {
	sep := strings.Split(str, " ")
	u := strings.Join(sep[:8], "")
	l := strings.Join(sep[8:], "")
	upper, err := strconv.ParseUint(u, 16, 64)
	if err != nil {
		panic(fmt.Errorf("upper: %s", err))
	}
	lower, err := strconv.ParseUint(l, 16, 64)
	if err != nil {
		panic(fmt.Errorf("lower: %s", err))
	}
	// fmt.Println(upper, lower, math.MaxInt64)

	return NewBitString(uint64(upper), uint64(lower))
}
