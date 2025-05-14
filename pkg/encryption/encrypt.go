package encryption

import (
	"github.com/ChainsAre2Tight/mgm-go/internal/interfaces"
	"github.com/ChainsAre2Tight/mgm-go/internal/nonce"
)

type encryptor struct {
	ng *nonce.NonceGenerator
}

func NewEncryptor(
	ng *nonce.NonceGenerator,
) interfaces.Encryptor {
	return &encryptor{
		ng: ng,
	}
}

func (e *encryptor) Encrypt(
	key string,
	associatedData string,
	plaintext string,
) (
	nonce *nonce.Nonce,
	ciphertext string,
	mac string,
	err error,
) {
	// schedule keys

	// get nonce

	// step 1: sign associated data
	// use nonce and convert 1||nonce to []byte
	// encrypt this nonce
	// increment result h times, encrypting each step in a goroutine
	// MULTIPLY encryption result with a corresponding accosiated data block
	// convert result to bitString and XOR it with cumulative result

	// step 2: encrypt plaintext and sign it (in parallel)
	// use nonce and convert 0||nonce to []byte
	// encrypt this nonce
	// increment result Q times, encrypting each step in a goroutine
	// XOR encrypted result with a corresponding plaintext block
	// convert []byte to bitString
	// increment MAC counter and multiply new encryption result with encrypted MAC counter result
	// XOR it with cumulative mac result

	// step 3: multiply and xor len(a) || len(c)

	// step 4: Encrypt cumulative result, taking first S bits from it
	panic("unimplemented")
}
