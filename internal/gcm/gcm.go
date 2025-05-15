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
	iv interfaces.BitString,
	depth int,
	increment func(interfaces.BitString),
	keys *kt.RoundKeys,
	ctx context.Context,
) ([]interfaces.BitString, error) {
	fail := func(err error) ([]interfaces.BitString, error) {
		return nil, fmt.Errorf("gcm.Seed: %s", err)
	}

	result := make([]interfaces.BitString, depth)
	for i := range result {
		increment(iv)
		result[i] = iv
	}

	errs, _ := errgroup.WithContext(ctx)

	for i := range result {
		errs.Go(func() error {
			if new, err := encryptBitString(result[i], keys); err != nil {
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

	panic("unimplemented")
}

func encryptBitString(
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
