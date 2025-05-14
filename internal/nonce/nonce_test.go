package nonce

import (
	"fmt"
	"strings"
	"testing"
)

func TestNonceString(t *testing.T) {
	tt := []struct {
		upper, lower uint64
		result       string
	}{
		{18446744073709551615, 18446744073709551615, strings.Repeat("f", 32)},
	}

	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%d|%d -> %s", td.upper, td.lower, td.result),
			func(t *testing.T) {
				n := &Nonce{
					upper: td.upper,
					lower: td.lower,
				}
				if res := n.String(); res != td.result {
					t.Fatalf("\nGot:  %s, \nWant: %s", res, td.result)
				}
			},
		)
	}
}
