package bitstrings

func IncrementL(bs *BitString128) {
	bs.upper++
}

func IncrementR(bs *BitString128) {
	bs.lower++
}
