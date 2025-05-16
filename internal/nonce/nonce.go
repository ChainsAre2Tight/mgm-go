package nonce

import (
	"sync"

	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
	"github.com/ChainsAre2Tight/mgm-go/internal/interfaces"
)

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

func (ng *NonceGenerator) Nonce() interfaces.BitString {
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
