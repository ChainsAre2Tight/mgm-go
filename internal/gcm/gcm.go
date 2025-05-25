package gcm

import (
	"context"
	"fmt"

	kuznechikgo "github.com/ChainsAre2Tight/kuznechik-go"
	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
	"github.com/ChainsAre2Tight/mgm-go/internal/interfaces"
	"golang.org/x/sync/errgroup"
)

func Seed(
	iv bitstrings.BitString128,
	depth int,
	increment func(*bitstrings.BitString128),
	keys kuznechikgo.UintRoundKeys,
	ctx context.Context,
) ([]*bitstrings.BitString128, error) {
	fail := func(err error) ([]*bitstrings.BitString128, error) {
		return nil, fmt.Errorf("gcm.Seed: %s", err)
	}

	encryptedIV := EncryptBitString(&iv, keys)

	result := make([]*bitstrings.BitString128, depth)
	for i := range result {
		new := *encryptedIV
		result[i] = &new
		increment(encryptedIV)
	}

	errs, _ := errgroup.WithContext(ctx)

	for i := range result {
		errs.Go(func() error {
			result[i] = EncryptBitString(result[i], keys)
			return nil
		})

	}

	if err := errs.Wait(); err != nil {
		return fail(fmt.Errorf("parallel encryption: %s", err))
	}

	return result, nil
}

func EncryptBitString(
	target interfaces.BitString,
	keys kuznechikgo.UintRoundKeys,
) *bitstrings.BitString128 {
	upper, lower := kuznechikgo.UintEncrypt(
		target.Upper(),
		target.Lower(),
		keys,
	)
	return bitstrings.NewBitString(upper, lower)
}
