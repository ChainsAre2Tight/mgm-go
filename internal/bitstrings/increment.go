package bitstrings

// IncrementL implements interfaces.BitString.
func (bs *BitString128) IncrementL() {
	bs.upper++
}

// IncremtntR implements interfaces.BitString.
func (bs *BitString128) IncremtntR() {
	bs.lower++
}
