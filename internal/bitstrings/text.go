package bitstrings

import "fmt"

func SliceFromText(in []byte) ([]*BitString128, uint64) {
	length := uint64(len(in)) * 8

	var result []*BitString128
	var flag bool // indicates that there is data to pad
	if len(in)%16 == 0 {
		result = make([]*BitString128, len(in)/16)
		flag = false
	} else {
		result = make([]*BitString128, len(in)/16, len(in)/16+1)
		flag = true
	}

	for i := range result {
		result[i] = FromBytes(in[16*i : 16*i+16])
	}
	if flag {
		result = append(result, FromBytes(in[len(result)*16:]))
	}

	return result, length
}

func TextFromSlice(in []*BitString128, length uint64) ([]byte, error) {
	fail := func(err error) ([]byte, error) {
		return nil, fmt.Errorf("bitstrings.TextFromSlice: %s", err)
	}
	l := len(in) * 16
	if length > uint64(l)*8 {
		return fail(fmt.Errorf("provided length is greated that length of provided slice: %d > %d (%d * 8)", length, l*8, l))
	}

	result := make([]byte, 0, l)
	for _, val := range in {
		result = append(result, val.Bytes()...)
	}

	return result[:length/8], nil
}
