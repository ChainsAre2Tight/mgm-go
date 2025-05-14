package bitstrings

import (
	"fmt"

	"github.com/ChainsAre2Tight/mgm-go/internal/interfaces"
)

var _ interfaces.BitString = (*BitString128)(nil)
var _ fmt.Stringer = (*BitString128)(nil)

type BitString128 struct {
	length int
	upper  uint64
	lower  uint64
}

// String implements fmt.Stringer.
func (bs *BitString128) String() string {
	return fmt.Sprintf("%0.64b%0.64b", bs.upper, bs.lower)[128-bs.length:]
}

// Length implements interfaces.BitString.
func (bs *BitString128) Length() int {
	return bs.length
}

func FromString(str string) (*BitString128, error) {
	fail := func(err error) (*BitString128, error) {
		return nil, fmt.Errorf("BitString128s.FromString: %s", err)
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
	return &BitString128{
		length: len(str),
		upper:  upper,
		lower:  lower,
	}, nil
}
