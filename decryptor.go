package mgmgo

import (
	"context"
	"fmt"
	"reflect"

	"github.com/ChainsAre2Tight/kuznechik-go/pkg/keyschedule"
	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
	"github.com/ChainsAre2Tight/mgm-go/internal/encryption"
	"github.com/ChainsAre2Tight/mgm-go/internal/maccomputation"
)

type decryptor struct{}

var _ Decryptor = (*decryptor)(nil)

func NewDecryptor() *decryptor {
	return &decryptor{}
}

// Decrypt implements Decryptor.
func (d *decryptor) Decrypt(
	key string,
	nonce []byte,
	associatedData []byte,
	ciphertext []byte,
	mac []byte,
) (
	plaintext []byte,
	err error,
) {
	ctx, cancel := context.WithCancel(context.Background())
	fail := func(err error) ([]byte, error) {
		cancel()
		return nil, fmt.Errorf("decryptor.Decrypt: %s", err)
	}
	// schedule keys
	keys, err := keyschedule.Schedule(key)
	if err != nil {
		return fail(fmt.Errorf("key schedule: %s", err))
	}

	// get nonce
	nonceRaw := bitstrings.FromBytes(nonce)

	// compute authenticated data and ciphertext length
	authenticatedDataArray, authenticatedDataLength := bitstrings.SliceFromText(associatedData)
	ciphertextArray, ciphertextLength := bitstrings.SliceFromText(ciphertext)

	// compute MAC
	macRaw, err := maccomputation.Compute(
		keys,
		nonceRaw,
		authenticatedDataArray,
		ciphertextArray,
		authenticatedDataLength,
		ciphertextLength,
		ctx,
	)
	if err != nil {
		return fail(fmt.Errorf("mac computation: %s", err))
	}

	// compare MACs
	if !reflect.DeepEqual(macRaw.Bytes(), mac) {
		// TODO add custom exception
		return fail(fmt.Errorf("MACs differ"))
	}

	// decrypt ciphertext
	plaintextArray, err := encryption.Decypt(ciphertextArray, keys, nonceRaw, ctx)
	if err != nil {
		return fail(fmt.Errorf("ciphertext decryption: %s", err))
	}

	plaintext, err = bitstrings.TextFromSlice(plaintextArray, ciphertextLength)
	if err != nil {
		return fail(fmt.Errorf("plaintext to bytes: %s", err))
	}

	return plaintext, nil
}
