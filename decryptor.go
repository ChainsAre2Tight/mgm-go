package mgmgo

import (
	"context"
	"fmt"
	"reflect"

	kuznechikgo "github.com/ChainsAre2Tight/kuznechik-go"
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
	key []byte,
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
	keys, err := kuznechikgo.Schedule(key)
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
	if b := macRaw.Bytes(); !reflect.DeepEqual(b, mac) {
		// TODO add custom exception
		return fail(&ErrMACsDiffer{
			err:       "",
			Computed:  b,
			Presented: mac,
		})
	}

	// decrypt ciphertext
	plaintextArray, err := encryption.Decrypt(ciphertextArray, keys, nonceRaw, ctx)
	if err != nil {
		return fail(fmt.Errorf("ciphertext decryption: %s", err))
	}

	plaintext, err = bitstrings.TextFromSlice(plaintextArray, ciphertextLength)
	if err != nil {
		return fail(fmt.Errorf("plaintext to bytes: %s", err))
	}

	return plaintext, nil
}
