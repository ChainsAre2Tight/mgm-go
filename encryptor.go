package mgmgo

import (
	"context"
	"fmt"

	"github.com/ChainsAre2Tight/kuznechik-go/pkg/keyschedule"
	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
	"github.com/ChainsAre2Tight/mgm-go/internal/encryption"
	"github.com/ChainsAre2Tight/mgm-go/internal/maccomputation"
	"github.com/ChainsAre2Tight/mgm-go/internal/nonce"
)

var _ Encryptor = (*encryptor)(nil)

type encryptor struct {
	ng nonce.NonceGenerator
}

func NewEncryptor(
	ng nonce.NonceGenerator,
) *encryptor {
	return &encryptor{
		ng: ng,
	}
}

func (e *encryptor) Encrypt(
	key string,
	associatedData []byte,
	plaintext []byte,
) (
	nonce []byte,
	ciphertext []byte,
	mac []byte,
	err error,
) {
	ctx, cancel := context.WithCancel(context.Background())
	fail := func(err error) ([]byte, []byte, []byte, error) {
		cancel()
		return nil, nil, nil, fmt.Errorf("encryptor.Encrypt: %s", err)
	}
	// schedule keys
	keys, err := keyschedule.Schedule(key)
	if err != nil {
		return fail(fmt.Errorf("key schedule: %s", err))
	}

	// get nonce
	nonceRaw := e.ng.Nonce()
	nonce = nonceRaw.Bytes()

	// compute authenticated data and ciphertext length
	authenticatedDataArray, authenticatedDataLength := bitstrings.SliceFromText(associatedData)
	plaintextArray, plaintextLength := bitstrings.SliceFromText(plaintext)

	// encrypt plaintext
	ciphertextArray, err := encryption.Encypt(plaintextArray, keys, nonceRaw, ctx)
	if err != nil {
		return fail(fmt.Errorf("plaintext encryption: %s", err))
	}

	// make MSB(u) out of the last ciphertext block
	if u := int(plaintextLength % 128); u != 0 {
		ciphertextArray[len(ciphertextArray)-1], err = bitstrings.MSB(ciphertextArray[len(ciphertextArray)-1], u)
		if err != nil {
			return fail(fmt.Errorf("MSB: %s", err))
		}
	}

	fmt.Println(bitstrings.RepresentPointerArray(ciphertextArray))

	// compute MAC
	macRaw, err := maccomputation.Compute(
		keys,
		nonceRaw,
		authenticatedDataArray,
		ciphertextArray,
		authenticatedDataLength,
		plaintextLength,
		ctx,
	)
	if err != nil {
		return fail(fmt.Errorf("mac computation: %s", err))
	}
	mac = macRaw.Bytes()

	ciphertext, err = bitstrings.TextFromSlice(ciphertextArray, plaintextLength)
	if err != nil {
		return fail(fmt.Errorf("ciphertex to bytes: %s", err))
	}

	return nonce, ciphertext, mac, nil
}
