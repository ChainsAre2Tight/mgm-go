package gcm

import (
	"fmt"
	"testing"

	"github.com/ChainsAre2Tight/kuznechik-go/pkg/keyschedule"
	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
)

func TestEncryptBitString(t *testing.T) {
	tt := []struct {
		key          string
		upper, lower uint64
		resU, resL   uint64
	}{
		{"12345678901234567890123456789012", 1, 2, 15496575308278909952, 9583859517063235840},
		{"12345678901234567890123456789012", 2, 2, 18232385423786432256, 1356391677369709056},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s + %d | %d ->  %d | %d", td.key, td.upper, td.lower, td.resU, td.resL),
			func(t *testing.T) {
				rawbs := bitstrings.NewBitString(
					td.upper,
					td.lower,
				)
				keys, err := keyschedule.Schedule(td.key)
				if err != nil {
					t.Fatalf("Error during keyscheduling: %s", err)
				}
				bs, err := encryptBitString(rawbs, keys)
				if err != nil {
					t.Fatalf("Error during encryption: %s", err)
				}
				if td.resU != bs.Upper() || td.resL != bs.Lower() {
					t.Fatalf("\nGot:  %d | %d, \nWant: %d | %d", bs.Upper(), bs.Lower(), td.resU, td.resL)
				}
			},
		)
	}
}
