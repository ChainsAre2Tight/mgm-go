package encryption

import (
	"context"
	"fmt"

	"github.com/ChainsAre2Tight/kuznechik-go/pkg/keyschedule"
	"github.com/ChainsAre2Tight/kuznechik-go/pkg/types"
	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
	"github.com/ChainsAre2Tight/mgm-go/internal/gcm"
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
	nonce interfaces.BitString,
	ciphertext string,
	mac string,
	err error,
) {
	fail := func(err error) (interfaces.BitString, string, string, error) {
		return nil, "", "", fmt.Errorf("encryption.Encrypt: %s", err)
	}
	// schedule keys
	keys, err := keyschedule.Schedule(key)
	if err != nil {
		return fail(fmt.Errorf("key schedule: %s", err))
	}
	_ = keys

	// get nonce
	nonce = e.ng.Nonce()

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

func encyptPlaintext(
	plaintext []*bitstrings.BitString128,
	keys *types.RoundKeys,
	nonce *bitstrings.BitString128,
	ctx context.Context,
) ([]*bitstrings.BitString128, error) {
	fail := func(err error) ([]*bitstrings.BitString128, error) {
		return nil, fmt.Errorf("encrypt.encryptPlaintext: %s", err)
	}

	gamma, err := gcm.Seed(
		*nonce,
		len(plaintext),
		bitstrings.IncrementR,
		keys,
		ctx,
	)

	if err != nil {
		return fail(fmt.Errorf("gamma: %s", err))
	}

	result := make([]*bitstrings.BitString128, len(gamma))
	for i := range gamma {
		result[i] = bitstrings.BitSum128(gamma[i], plaintext[i])
	}

	return result, nil
}

func computeMAC(
	keys *types.RoundKeys,
	nonce *bitstrings.BitString128,
	authenticatedData, ciphertext []*bitstrings.BitString128,
	lengthAuth, lengthPlain uint64,
	ctx context.Context,
) (*bitstrings.BitString128, error) {
	fail := func(err error) (*bitstrings.BitString128, error) {
		return nil, fmt.Errorf("encryption.computeMAC: %s", err)
	}

	gamma, err := gcm.Seed(
		*bitstrings.NewBitString(
			nonce.Upper()|1<<63,
			nonce.Lower(),
		),
		len(ciphertext)+len(authenticatedData)+1,
		bitstrings.IncrementL,
		keys,
		ctx,
	)

	var mac *bitstrings.BitString128 = bitstrings.NewBitString(0, 0)

	if err != nil {
		return fail(fmt.Errorf("gamma: %s", err))
	}

	for i := range authenticatedData {
		mac = bitstrings.BitSum128(
			mac,
			bitstrings.BitMul(
				authenticatedData[i],
				gamma[i],
			))
	}

	for i := range ciphertext {
		mac = bitstrings.BitSum128(
			mac,
			bitstrings.BitMul(
				gamma[i+len(authenticatedData)],
				ciphertext[i],
			))
	}

	mac = bitstrings.BitSum128(
		mac,
		bitstrings.BitMul(
			gamma[len(authenticatedData)+len(ciphertext)],
			bitstrings.NewBitString(
				lengthAuth,
				lengthPlain,
			),
		))

	mac, err = gcm.EncryptBitString(mac, keys)
	if err != nil {
		return fail(fmt.Errorf("mac encryption: %s", err))
	}

	return mac, nil
}
