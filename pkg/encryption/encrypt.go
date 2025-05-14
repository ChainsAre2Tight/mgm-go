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
	panic("unimplemented")
}
