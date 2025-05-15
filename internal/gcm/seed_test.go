package gcm_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/ChainsAre2Tight/kuznechik-go/pkg/keyschedule"
	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
	"github.com/ChainsAre2Tight/mgm-go/internal/gcm"
)

func TestSeed(t *testing.T) {
	tt := []struct {
		iv        bitstrings.BitString128
		depth     int
		key       string
		increment func(*bitstrings.BitString128)
		result    []*bitstrings.BitString128
	}{
		{
			iv:        *bitstrings.NewBitString(0, 2),
			depth:     1,
			key:       "12345678901234567890123456789012",
			increment: bitstrings.IncrementL,
			result: []*bitstrings.BitString128{
				bitstrings.NewBitString(15496575308278909952, 9583859517063235840),
			},
		}, {
			iv:        *bitstrings.NewBitString(0, 2),
			depth:     2,
			key:       "12345678901234567890123456789012",
			increment: bitstrings.IncrementL,
			result: []*bitstrings.BitString128{
				bitstrings.NewBitString(15496575308278909952, 9583859517063235840),
				bitstrings.NewBitString(18232385423786432256, 1356391677369709056),
			},
		},
	}

	// TODO: fix prints
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s * %d + %s (%s) -> %v", td.iv.String(), td.depth, td.key, reflect.TypeOf(td.increment).Name(), td.result),
			func(t *testing.T) {
				keys, err := keyschedule.Schedule(td.key)
				if err != nil {
					t.Fatalf("Error during keyschedule: %s", err)
				}
				res, err := gcm.Seed(
					td.iv,
					td.depth,
					td.increment,
					keys,
					context.Background(),
				)
				if err != nil {
					t.Fatalf("Error during seeding: %s", err)
				}
				if !reflect.DeepEqual(res, td.result) {
					t.Fatalf("\nGot:  %v, \nWant: %v.", res, td.result)
				}
			},
		)
	}
}
