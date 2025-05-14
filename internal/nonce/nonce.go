package nonce

import (
	"fmt"
	"sync"
)

var _ fmt.Stringer = (*Nonce)(nil)

type Nonce struct {
	upper uint64
	lower uint64
}

// String implements fmt.Stringer.
func (n *Nonce) String() string {
	return fmt.Sprintf("%0.16x%0.16x", n.upper, n.lower)
}

// func FromString(str string) *Nonce {

// }

type NonceGenerator struct {
	counterUpper uint64
	counterLower uint64
	mu           *sync.Mutex
}

func NewNonceGenerator() *NonceGenerator {
	return &NonceGenerator{
		counterUpper: 0,
		counterLower: 0,
		mu:           &sync.Mutex{},
	}
}

func (ng *NonceGenerator) Nonce() *Nonce {
	ng.mu.Lock()
	defer ng.mu.Unlock()

	// we increment upper counter when lower overflows
	ng.counterLower++
	if ng.counterLower == 0 {
		ng.counterUpper++
	}

	// we need only 63 bits of upper counter so reset the counter when it exceeds this value
	if ng.counterUpper >= 1<<63 {
		ng.counterLower = 1
		ng.counterUpper = 1
	}

	return &Nonce{
		upper: ng.counterUpper,
		lower: ng.counterLower,
	}
}
