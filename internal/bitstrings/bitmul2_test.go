package bitstrings

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBitMulRawFromStrings(t *testing.T) {
	tt := []struct {
		a string
		b string
		c []uint64
	}{
		{"1", "1", []uint64{0, 0, 0, 1}},
		{"11", "1", []uint64{0, 0, 0, 3}},
		{"1001", "1", []uint64{0, 0, 0, 9}},

		{"1001", "0", []uint64{0, 0, 0, 0}},

		{"1001", "101", []uint64{0, 0, 0, 45}},
		{"1101010", "1010111", []uint64{0, 0, 0, 7478}},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s * %s -> %v", td.a, td.b, td.c),
			func(t *testing.T) {
				a, _ := FromString(td.a)
				b, _ := FromString(td.b)
				if res := bitMulPreRemainder(a, b); !reflect.DeepEqual(res, td.c) {
					t.Fatalf("\nGot:  %v, \nWant: %v", res, td.c)
				}
			},
		)
	}
}

func TestBitMulRaw(t *testing.T) {
	tt := []struct {
		upperA, lowerA uint64
		upperB, lowerB uint64
		result         []uint64
	}{
		{123, 0, 0, 1, []uint64{0, 0, 123, 0}},
		{123, 321, 0, 1, []uint64{0, 0, 123, 321}},
		{1, 0, 1, 0, []uint64{0, 1, 0, 0}},
		{ // x^75 + x^69 + x^32 + x^9 + 1
			// * x^100 + x^90 + x^50
			(1 << 11) + (1 << 5),
			(1 << 32) + (1 << 9) + 1,
			(1 << 36) + (1 << 26),
			1 << 50,
			[]uint64{
				0,
				(1 << 47) + (1 << 41) + (1 << 37) + (1 << 31) + (1 << 4),
				(1 << 61) + (1 << 58) + (1 << 55) + (1 << 45) + (1 << 36) + (1 << 35) + (1 << 26) + (1 << 18),
				(1 << 59) + (1 << 50),
			},
		},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%d|%d * %d|%d -> %v", td.upperA, td.lowerA, td.upperB, td.lowerB, td.result),
			func(t *testing.T) {
				a := &BitString128{
					upper: td.upperA,
					lower: td.lowerA,
				}
				b := &BitString128{
					upper: td.upperB,
					lower: td.lowerB,
				}
				if res := bitMulPreRemainder(a, b); !reflect.DeepEqual(res, td.result) {
					t.Fatalf("\nGot:  %v, \nWant: %v", res, td.result)
				}
			},
		)
	}
}
