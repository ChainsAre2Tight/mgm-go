package bitstrings

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestBytes(t *testing.T) {
	tt := []struct {
		upper, lower uint64
		bytes        []byte
	}{
		{1, 2, []byte{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2}},
		{math.MaxUint64, math.MaxUint64, []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}},
	}

	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%d | %d -> %v", td.upper, td.lower, td.bytes),
			func(t *testing.T) {
				bs := &BitString128{
					upper: td.upper,
					lower: td.lower,
				}
				if res := bs.Bytes(); !reflect.DeepEqual(res, td.bytes) {
					t.Fatalf("\nGot:  %v,\nWant: %v\n", res, td.bytes)
				}
			},
		)
	}
}
