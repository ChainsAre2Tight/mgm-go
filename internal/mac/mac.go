package mac

import (
	"context"
	"fmt"

	"github.com/ChainsAre2Tight/kuznechik-go/pkg/types"
	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
	"github.com/ChainsAre2Tight/mgm-go/internal/gcm"
)

func Compute(
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
