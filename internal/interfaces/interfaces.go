package interfaces

// Represents a bit string of a certain length
type BitString interface {
	Length() int
	Bytes() []byte
	Upper() uint64
	Lower() uint64
}
