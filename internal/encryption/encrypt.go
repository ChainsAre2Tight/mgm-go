package encryption

import (
	"context"
	"fmt"

	"github.com/ChainsAre2Tight/kuznechik-go/pkg/types"
	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
	"github.com/ChainsAre2Tight/mgm-go/internal/gcm"
)

func Encypt(
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
