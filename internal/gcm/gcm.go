package gcm

import (
	"context"
	"fmt"

	"github.com/ChainsAre2Tight/kuznechik-go/pkg/block"
	kt "github.com/ChainsAre2Tight/kuznechik-go/pkg/types"
	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
	"github.com/ChainsAre2Tight/mgm-go/internal/interfaces"
	"golang.org/x/sync/errgroup"
)

func Seed(
	iv bitstrings.BitString128,
	depth int,
	increment func(*bitstrings.BitString128),
	keys *kt.RoundKeys,
	ctx context.Context,
) ([]*bitstrings.BitString128, error) {
	fail := func(err error) ([]*bitstrings.BitString128, error) {
		return nil, fmt.Errorf("gcm.Seed: %s", err)
	}

	encryptedIV, err := EncryptBitString(&iv, keys)
	if err != nil {
		return fail(fmt.Errorf("error during initial iv encryption: %s", err))
	}

	result := make([]*bitstrings.BitString128, depth)
	for i := range result {
		new := *encryptedIV
		result[i] = &new
		increment(encryptedIV)
	}

	errs, _ := errgroup.WithContext(ctx)

	for i := range result {
		errs.Go(func() error {
			if new, err := EncryptBitString(result[i], keys); err != nil {
				return fmt.Errorf("error in gouroutine at %d: %s", i, err)
			} else {
				result[i] = new
			}
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
	keys *kt.RoundKeys,
) (*bitstrings.BitString128, error) {
	bytes, err := block.Encrypt(
		target.Bytes(),
		keys,
	)
	if err != nil {
		return nil, fmt.Errorf("gcm.encryptBitString: %s", err)
	}
	return bitstrings.FromBytes(bytes), nil
}
