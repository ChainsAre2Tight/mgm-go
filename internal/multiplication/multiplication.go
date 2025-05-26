package multiplication

const GaloisPolynomial uint64 = 0b10000111

func MultiplyUint128(v1_upper, v1_lower, v2_upper, v2_lower uint64) (uint64, uint64) {
	if v1_upper == 0 && v1_lower == 0 || v2_upper == 0 && v2_lower == 0 {
		return 0, 0
	}

	var v0, v1, v2, v3 uint64
	var lower, middle, upper uint64
	for i := range 64 {
		if low, high := (1<<i)&v1_lower > 0, (1<<i)&v1_upper > 0; low || high {
			lower = v2_lower << i
			middle = (v2_lower >> (64 - i)) + (v2_upper << i)
			upper = v2_upper >> (64 - i)
			if low {
				v3 ^= lower
				v2 ^= middle
				v1 ^= upper
			}
			if high {
				v2 ^= lower
				v1 ^= middle
				v0 ^= upper
			}
		}
	}

	upper, lower = v2, v3
	for i := 63; i >= 0; i-- {
		if v0&(1<<i) > 0 {
			upper ^= GaloisPolynomial << i

			// it is also apparantly necessary to xor overlaps with bits 128-191
			v1 ^= GaloisPolynomial >> (64 - i)
		}
		if v1&(1<<i) > 0 {
			lower ^= GaloisPolynomial << i
			upper ^= GaloisPolynomial >> (64 - i)
		}
	}

	return upper, lower
}
