package bitstrings

func BitSum(a, b *bitString) (*bitString, error) {
	length := max(a.Length(), b.Length())

	return &bitString{
		length: length,
		lower:  a.lower ^ b.lower,
		upper:  a.upper ^ b.upper,
	}, nil
}
