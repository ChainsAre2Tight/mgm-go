package bitstrings_test

import (
	"fmt"
	"testing"

	"github.com/ChainsAre2Tight/mgm-go/internal/bitstrings"
)

func TestBitstringMultiplication(t *testing.T) {
	tt := []struct {
		a string
		b string
		c string
	}{
		{"1010", "1", "1010"},
		{"1101", "1011", "1111111"},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s x %s -> %s", td.a, td.b, td.c),
			func(t *testing.T) {
				a, _ := bitstrings.FromString(td.a)
				b, _ := bitstrings.FromString(td.b)
				if res := bitstrings.BitMul(a, b); res.String() != td.c {
					t.Fatalf("\nGot:  %s, \nWant: %s", res, td.c)
				}
			},
		)
	}
}

func TestBitMul(t *testing.T) {
	tt := []struct {
		a, b, c *bitstrings.BitString128
	}{
		{
			a: bitstrings.NewBitString(1<<63, 0b1000011),
			b: bitstrings.NewBitString(0, 0b10),
			c: bitstrings.NewBitString(0, 1),
		}, {
			a: bitstrings.NewBitString(0, 0b10000111),
			b: bitstrings.NewBitString(1<<63, 0),
			c: bitstrings.NewBitString(1<<63, (1<<13)+(1<<6)+0b1001),
		}, {
			a: bitstrings.FromGOSTString("8D B1 87 D6 53 83 0E A4 BC 44 64 76 95 2C 30 0B"),
			b: bitstrings.FromGOSTString("02 02 02 02 02 02 02 02 01 01 01 01 01 01 01 01"),
			c: bitstrings.FromGOSTString("4C F4 27 F4 AD B7 5C F4 C0 DA 39 D5 AB 48 CF 38"),
		}, {
			a: bitstrings.FromGOSTString("7A 24 F7 26 30 E3 76 37 21 C8 F3 CD B1 DA 0E 31"),
			b: bitstrings.FromGOSTString("04 04 04 04 04 04 04 04 03 03 03 03 03 03 03 03"),
			c: bitstrings.BitSum128(
				bitstrings.FromGOSTString("4C F4 27 F4 AD B7 5C F4 C0 DA 39 D5 AB 48 CF 38"),
				bitstrings.FromGOSTString("94 95 44 0E F6 24 A1 DD C6 F5 D9 77 28 50 C5 73"),
			),
		}, {
			b: bitstrings.FromGOSTString("EA 05 05 05 05 05 05 05 05 00 00 00 00 00 00 00"),
			a: bitstrings.FromGOSTString("44 11 96 21 17 D2 06 35 C5 25 E0 A2 4D B4 B9 0A"),
			c: bitstrings.BitSum128(
				bitstrings.FromGOSTString("94 95 44 0E F6 24 A1 DD C6 F5 D9 77 28 50 C5 73"),
				bitstrings.FromGOSTString("A4 9A 8C D8 A6 F2 74 23 DB 79 E4 4A B3 06 D9 42"),
			),
		}, {
			a: bitstrings.FromGOSTString("D8 C9 62 3C 4D BF E8 14 CE 7C 1C 0C EA A9 59 DB"), // h4
			b: bitstrings.FromGOSTString("A9 75 7B 81 47 95 6E 90 55 B8 A3 3D E8 9F 42 FC"), // c1
			c: bitstrings.BitSum128(
				bitstrings.FromGOSTString("A4 9A 8C D8 A6 F2 74 23 DB 79 E4 4A B3 06 D9 42"),
				bitstrings.FromGOSTString("09 FE 3F 6A 83 3C 21 B3 90 27 D0 20 6A 84 E1 5A"),
			),
		}, { //c2
			a: bitstrings.FromGOSTString("A5 E1 F1 95 33 3E 14 82 96 99 31 BF BE 6D FD 43"), // h5
			b: bitstrings.FromGOSTString("80 75 D2 21 2B F9 FD 5B D3 F7 06 9A AD C1 6B 39"), //c2
			c: bitstrings.BitSum128(
				bitstrings.FromGOSTString("09 FE 3F 6A 83 3C 21 B3 90 27 D0 20 6A 84 E1 5A"),
				bitstrings.FromGOSTString("B5 DA 26 BB 00 EB A8 04 35 D7 97 6B C6 B5 46 4D"),
			),
		}, { //c3
			a: bitstrings.FromGOSTString("B4 CA 80 8C AC CF B3 F9 17 24 E4 8A 2C 7E E9 D2"), // h6
			b: bitstrings.FromGOSTString("49 7A B1 59 15 A6 BA 85 93 6B 5D 0E A9 F6 85 1C"), //c3
			c: bitstrings.BitSum128(
				bitstrings.FromGOSTString("B5 DA 26 BB 00 EB A8 04 35 D7 97 6B C6 B5 46 4D"),
				bitstrings.FromGOSTString("DD 1C 0E EE F7 83 C8 EB 2A 33 F3 58 D7 23 0E E5"),
			),
		}, { //c4
			a: bitstrings.FromGOSTString("72 90 8F C0 74 E4 69 E8 90 1B D1 88 EA 91 C3 31"), // h7
			b: bitstrings.FromGOSTString("C6 0C 14 D4 D3 F8 83 D0 AB 94 42 06 95 C7 6D EB"), // c4
			c: bitstrings.BitSum128(
				bitstrings.FromGOSTString("DD 1C 0E EE F7 83 C8 EB 2A 33 F3 58 D7 23 0E E5"),
				bitstrings.FromGOSTString("89 6C E1 08 32 EB EA F9 06 9F 3F 73 76 59 4D 40"),
			),
		}, { // c5
			a: bitstrings.FromGOSTString("23 CA 27 15 B0 2C 68 31 3B FD AC B3 9E 4D 0F B8"), //h8
			b: bitstrings.FromGOSTString("2C 75 52 00 00 00 00 00 00 00 00 00 00 00 00 00"), //c5
			c: bitstrings.BitSum128(
				bitstrings.FromGOSTString("89 6C E1 08 32 EB EA F9 06 9F 3F 73 76 59 4D 40"),
				bitstrings.FromGOSTString("99 1A F5 C9 D0 80 F7 63 87 FE 64 9E 7C 93 C6 42"),
			),
		}, { // len
			a: bitstrings.FromGOSTString("00 00 00 00 00 00 01 48 00 00 00 00 00 00 02 18"),
			b: bitstrings.FromGOSTString("BC BC E6 C4 1A A3 55 A4 14 88 62 BF 64 BD 83 0D"),
			c: bitstrings.BitSum128(
				bitstrings.FromGOSTString("99 1A F5 C9 D0 80 F7 63 87 FE 64 9E 7C 93 C6 42"),
				bitstrings.FromGOSTString("C0 C7 22 DB 5E 0B D6 DB 25 76 73 83 3D 56 71 28"),
			),
		},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%x, %x * %x, %x -> %x, %x", td.a.Upper(), td.a.Lower(), td.b.Upper(), td.b.Lower(), td.c.Upper(), td.c.Lower()),
			func(t *testing.T) {
				if res := bitstrings.BitMul(td.a, td.b); res.Upper() != td.c.Upper() || res.Lower() != td.c.Lower() {
					t.Fatalf("\nGot:  %0.64b, %0.64b, \nWant: %0.64b, %0.64b", res.Upper(), res.Lower(), td.c.Upper(), td.c.Lower())
				}
			},
		)
	}
}
