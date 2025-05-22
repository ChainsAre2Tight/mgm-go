package mgmgo

import (
	"context"
	"fmt"

	kuznechikgo "github.com/ChainsAre2Tight/kuznechik-go"
	ad "github.com/ChainsAre2Tight/mgm-go/internal/associateddata"
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

func (e *encryptor) OldEncrypt(
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
	fail := func(err error) ([]byte, []byte, []byte, error) {
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

	// compute authenticated data and ciphertext length
	lengthAuth := uint64(len(associatedData)) * 8
	lengthPlaintext := uint64(len(plaintext)) * 8

	// create MAC of associated data
	macUpper, macLower, counterUpper, counterLower, err := ad.ComputeADMAC(keys, nonceRaw.Upper(), nonceRaw.Lower(), associatedData)
	if err != nil {
		return fail(fmt.Errorf("ad: %s", err))
	}

	ciphertext, mac, err = encryption.EncryptAndComputeMAC(
		keys, nonceRaw.Upper(), nonceRaw.Lower(), counterUpper, counterLower, macUpper, macLower, plaintext, lengthAuth, lengthPlaintext,
	)

	if err != nil {
		return fail(fmt.Errorf("encryption: %s", err))
	}

	return nonceRaw.Bytes(), ciphertext, mac, nil
}
