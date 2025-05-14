package interfaces

import "github.com/ChainsAre2Tight/mgm-go/internal/nonce"

// Represents a bit string of a certain length
type BitString interface {
	Length() int
}

type Encryptor interface {
	Encrypt(
		key string,
		associatedData string,
		plaintext string,
	) (
		nonce *nonce.Nonce,
		ciphertext string,
		mac string,
		err error,
	)
}
