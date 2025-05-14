package bitstrings

import "fmt"

func BitMul(a, b *BitString128) (*BitString128, error) {
	// calcualte full, 256 bit multiplication result
	raw := bitMulPreRemainder(a, b)

	// bring result back into the field
	// by calculating the remainder
	fmt.Println(raw)

	return nil, fmt.Errorf("not implemented")
}

func bitMulPreRemainder(a, b *BitString128) []uint64 {
	var res = make([]uint64, 4)

	// iterate through lower a
	// if bit at i-th position is present,
	// xor corresponding result bits
	var lower, middle, upper uint64
	for i := range 64 {
		if (1<<i)&a.lower > 0 {
			lower = b.lower << i
			middle = (b.lower >> (64 - i)) + (b.upper << i)
			upper = b.upper >> (64 - i)

			res[3] ^= lower
			res[2] ^= middle
			res[1] ^= upper
		}
	}

	for i := range 64 {
		if (1<<i)&a.upper > 0 {
			lower = b.lower << i
			middle = (b.lower >> (64 - i)) + (b.upper << i)
			upper = b.upper >> (64 - i)

			res[2] ^= lower
			res[1] ^= middle
			res[0] ^= upper
		}
	}

	return res
}
