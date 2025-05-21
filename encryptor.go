package mgmgo

import (
	"context"
	"fmt"

	kuznechikgo "github.com/ChainsAre2Tight/kuznechik-go"
	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
	"github.com/ChainsAre2Tight/mgm-go/internal/encryption"
	"github.com/ChainsAre2Tight/mgm-go/internal/maccomputation"
)

var _ Encryptor = (*encryptor)(nil)

type encryptor struct {
	ng NonceGenerator
}

func NewEncryptor(
	ng NonceGenerator,
) *encryptor {
	return &encryptor{
		ng: ng,
	}
}

func (e *encryptor) Encrypt(
	key []byte,
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
	k, err := kuznechikgo.Schedule(key)
	if err != nil {
		return fail(fmt.Errorf("key schedule: %s", err))
	}
	// convert keys for new Uint64 functions added in kuznechik-go v1.1
	keys := kuznechikgo.KeysToUints(k)

	// get nonce
	nonceRaw := e.ng.Nonce()
	nonce = nonceRaw.Bytes()

	// compute authenticated data and ciphertext length
	authenticatedDataArray, authenticatedDataLength := bitstrings.SliceFromText(associatedData)
	plaintextArray, plaintextLength := bitstrings.SliceFromText(plaintext)

	// encrypt plaintext
	ciphertextArray, err := encryption.Encrypt(plaintextArray, keys, nonceRaw, ctx)
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
