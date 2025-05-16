package bitstrings

import "fmt"

func BitMul(a, b *BitString128) *BitString128 {
	if a.lower == 0 && a.upper == 0 || b.lower == 0 && b.upper == 0 {
		return &BitString128{
			length: 0,
			upper:  0,
			lower:  0,
		}
	}
	// calcualte full, 256 bit multiplication result
	raw := bitMulPreRemainder(a, b)

	// bring result back into the field
	// by calculating the remainder
	upper, lower := bitMulRemainder(raw)

	// find first non-zero bit to calculate lesngth
	var length int
	if upper > 0 {
		length = firstNonZeroBit(upper) + 64
	} else if lower > 0 {
		length = firstNonZeroBit(lower)
	} else {
		length = 0
	}

	fmt.Println(upper, lower)

	return &BitString128{
		length: length,
		upper:  upper,
		lower:  lower,
	}
}

func firstNonZeroBit(n uint64) int {
	fmt.Println(n)
	for i := 63; i >= 0; i-- {
		if n&uint64(1<<i) > 0 {
			return i + 1
		}
	}
	return 0
}

func bitMulPreRemainder(a, b *BitString128) []uint64 {
	var res = make([]uint64, 4)

	// iterate through lower and upper a's bits
	// if bit at i-th position is present,
	// calculate shifted b's bits and
	// xor corresponding result bits
	var lower, middle, upper uint64
	for i := range 64 {
		if low, high := (1<<i)&a.lower > 0, (1<<i)&a.upper > 0; low || high {
			lower = b.lower << i
			middle = (b.lower >> (64 - i)) + (b.upper << i)
			upper = b.upper >> (64 - i)
			if low {
				res[3] ^= lower
				res[2] ^= middle
				res[1] ^= upper
			}
			if high {
				res[2] ^= lower
				res[1] ^= middle
				res[0] ^= upper
			}
		}
	}

	return res
}

type galoisGenPolynomial struct {
	high uint64 // x^128
	low  uint64 // x^7+x^2+x+1
}

var GaloisGenPolynomial = &galoisGenPolynomial{
	high: 1,                                  // x^128
	low:  (1 << 7) + (1 << 2) + (1 << 1) + 1, // x^7+x^2+x+1
}

func bitMulRemainder(raw []uint64) (upper, lower uint64) {
	if len(raw) != 4 {
		panic(fmt.Sprintf("bitstrings.bitMulRemainder: unexpected multiplication slice length, expected 4, got %d", len(raw)))
	}

	upper = raw[2]
	lower = raw[3]

	for i := 63; i >= 0; i-- {
		if raw[0]&1<<i > 0 {
			upper ^= GaloisGenPolynomial.low << i
		}
		if raw[1]&1<<i > 0 {
			lower ^= GaloisGenPolynomial.low << i
		}
	}

	return upper, lower
}
