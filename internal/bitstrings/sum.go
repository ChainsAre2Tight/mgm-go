package bitstrings

func BitSum128(a, b *BitString128) *BitString128 {
	length := max(a.Length(), b.Length())

	return &BitString128{
		length: length,
		lower:  a.lower ^ b.lower,
		upper:  a.upper ^ b.upper,
	}
}
