package utils

func BytesToUint64WithPadding(in []byte) (uint64, uint64) {
	var upper, lower uint64 = 0, 0
	var i int = 0
	for ; i < 8 && i < len(in); i++ {
		upper += uint64(in[i]) << uint64(56-(8*i))
	}
	for ; i < 16 && i < len(in); i++ {
		lower += uint64(in[i]) << uint64(120-(8*i))
	}
	return upper, lower
}

func Uint64ToBytesWithPadding(upper, lower uint64, dst []byte) {
	var i int = 0
	for ; i < 8 && i < len(dst); i++ {
		dst[i] = byte(upper >> (56 - 8*i))
	}
	for ; i < 16 && i < len(dst); i++ {
		dst[i] = byte(lower >> (120 - 8*i))
	}
}
