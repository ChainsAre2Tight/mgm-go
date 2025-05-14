package interfaces

import "fmt"

// Represents a 128-bit string
type BitString interface {
	fmt.Stringer
	Length() int
}
