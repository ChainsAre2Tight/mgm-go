package nonce

import (
	"sync"

	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
)

type NonceGenerator interface {
	Nonce() *bitstrings.BitString128
}

type nonceGenerator struct {
	counterUpper uint64
	counterLower uint64
	mu           *sync.Mutex
}

func NewNonceGenerator() *nonceGenerator {
	return &nonceGenerator{
		counterUpper: 0,
		counterLower: 0,
		mu:           &sync.Mutex{},
	}
}

func (ng *nonceGenerator) Nonce() *bitstrings.BitString128 {
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

	return bitstrings.NewBitString(
		ng.counterUpper,
		ng.counterLower,
	)
}
