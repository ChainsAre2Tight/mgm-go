package bitstrings

import (
	"fmt"

	"github.com/ChainsAre2Tight/mgm-go/internal/interfaces"
)

var _ interfaces.BitString = (*bitString)(nil)

// var _ fmt.Stringer = (*bitString)(nil)

type bitString struct {
	length int
	upper  uint64
	lower  uint64
}

// String implements fmt.Stringer.
func (bs *bitString) String() string {
	return fmt.Sprintf("%0.64b%0.64b", bs.upper, bs.lower)[128-bs.length:]
}

// Length implements interfaces.BitString.
func (bs *bitString) Length() int {
	return bs.length
}

func FromString(str string) (*bitString, error) {
	fail := func(err error) (*bitString, error) {
		return nil, fmt.Errorf("bitstrings.FromString: %s", err)
	}
	if l := len(str); l > 128 {
		return fail(fmt.Errorf("exceeded maximum bitsring length (%d > 128)", l))
	}
	var upper, lower uint64

	for i := 0; i < len(str); i++ {
		r := rune(str[len(str)-i-1])
		switch r {
		case '1':
			if i < 64 {
				lower += 1 << i
			} else {
				upper += 1<<i - 64
			}
		case '0':
		default:
			return fail(fmt.Errorf("unsupported characted %s of string %s at position %d", string(r), str, len(str)-i-1))
		}
	}
	return &bitString{
		length: len(str),
		upper:  upper,
		lower:  lower,
	}, nil
}
