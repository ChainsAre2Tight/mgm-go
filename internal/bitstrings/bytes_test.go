package bitstrings

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestToBytes(t *testing.T) {
	tt := []struct {
		upper, lower uint64
		bytes        []byte
	}{
		{1, 2, []byte{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2}},
		{math.MaxUint64, math.MaxUint64, []byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}},
		{0x1122334455667700, 0xFFEEDDCCBBAA9988, []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x00, 0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA, 0x99, 0x88}},
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

func TestFromBytes(t *testing.T) {
	tt := []struct {
		upper, lower uint64
		bytes        []byte
	}{
		{0x1122334455667700, 0xFFEEDDCCBBAA9988, []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x00, 0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA, 0x99, 0x88}},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%v -> %x | %x", td.bytes, td.upper, td.lower),
			func(t *testing.T) {
				b := FromBytes(td.bytes)
				if b.upper != td.upper || b.lower != td.lower {
					t.Fatalf("\nGot:  %x, %x, \nWant: %x, %x.", b.upper, b.lower, td.upper, td.lower)
				}
			},
		)
	}
}
