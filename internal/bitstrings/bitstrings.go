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

// Lower implements interfaces.BitString.
func (bs *BitString128) Lower() uint64 {
	return bs.lower
}

// Upper implements interfaces.BitString.
func (bs *BitString128) Upper() uint64 {
	return bs.upper
}

func NewBitString(upper, lower uint64) *BitString128 {
	return &BitString128{
		upper: upper,
		lower: lower,
	}
}

// Bytes implements interfaces.BitString.
func (bs *BitString128) Bytes() []byte {
	result := make([]byte, 16)
	for i := range 8 {
		result[i] = byte((bs.upper << (8 * i)) >> 56)
		result[i+8] = byte((bs.lower << (8 * i)) >> 56)
	}
	// fmt.Println(result)
	return result
}

// String implements fmt.Stringer.
func (bs *BitString128) String() string {
	return fmt.Sprintf("%0.64b%0.64b", bs.upper, bs.lower)[128-bs.length:]
}

// Length implements interfaces.BitString.
func (bs *BitString128) Length() int {
	return bs.length
}

func FromBytes(b []byte) *BitString128 {
	var upper, lower uint64

	if len(b) < 16 {
		temp := make([]byte, 16)
		copy(temp, b)
		b = temp
	}

	for i := range 8 {
		upper += uint64(b[i]) << (56 - 8*i)
		lower += uint64(b[i+8]) << (56 - 8*i)
	}

	return &BitString128{
		upper: upper,
		lower: lower,
	}
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
