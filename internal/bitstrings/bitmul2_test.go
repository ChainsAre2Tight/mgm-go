package bitstrings

import (
	"fmt"
	"reflect"
	"testing"
)

func TestBitMulRaw(t *testing.T) {
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

func TestBitSum(t *testing.T) {
	tt := []struct {
		upperA, lowerA uint64
		upperB, lowerB uint64
		result         []uint64
	}{
		{123, 0, 0, 1, []uint64{0, 0, 123, 0}},
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
