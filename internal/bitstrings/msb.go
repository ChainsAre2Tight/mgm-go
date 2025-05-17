package bitstrings

import "fmt"

func MSB(in *BitString128, length int) (*BitString128, error) {
	if length >= 128 {
		return nil, fmt.Errorf("bitstrings.MSB: Maximum length for MSB excceded: %d > 128", length)
	}
	if length > 64 {
		return NewBitString(
			in.Upper(),
			(in.Lower()>>(128-length))<<(128-length),
		), nil
	} else {
		return NewBitString(
			(in.Upper()>>(64-length))<<(64-length),
			0,
		), nil
	}
}
