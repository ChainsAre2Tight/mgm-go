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

func TestNoceFromString(t *testing.T) {
	tt := []struct {
		in           string
		upper, lower uint64
	}{
		{"f", 0, 15},
		{"f0", 0, 15 * 16},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s -> %d | %d", td.in, td.upper, td.lower),
			func(t *testing.T) {
				if n, err := FromString(td.in); err != nil {
					t.Fatalf("\nError: %s", err)
				} else if n.lower != td.lower || n.upper != td.upper {
					t.Fatalf("\nGot:  %0.16x | %0.16x, \nWant: %0.16x | %0.16x", n.upper, n.lower, td.upper, td.lower)
				}

			},
		)
	}
}
