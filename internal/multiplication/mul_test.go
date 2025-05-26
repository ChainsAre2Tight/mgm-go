package multiplication_test

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ChainsAre2Tight/mgm-go/internal/multiplication"
)

func TestMultiplication(t *testing.T) {
	tt := []struct {
		a, b, c string
	}{
		{
			a: "8DB187D653830EA4BC446476952C300B",
			b: "02020202020202020101010101010101",
			c: "4CF427F4ADB75CF4C0DA39D5AB48CF38",
			// }, {
			// 	a: "7A 24 F7 26 30 E3 76 37 21 C8 F3 CD B1 DA 0E 31",
			// 	b: "04 04 04 04 04 04 04 04 03 03 03 03 03 03 03 03",
			// 	c: bitstrings.BitSum128(
			// 		"4C F4 27 F4 AD B7 5C F4 C0 DA 39 D5 AB 48 CF 38",
			// 		"94 95 44 0E F6 24 A1 DD C6 F5 D9 77 28 50 C5 73",
			// 	),
			// }, {
			// 	b: "EA 05 05 05 05 05 05 05 05 00 00 00 00 00 00 00",
			// 	a: "44 11 96 21 17 D2 06 35 C5 25 E0 A2 4D B4 B9 0A",
			// 	c: bitstrings.BitSum128(
			// 		"94 95 44 0E F6 24 A1 DD C6 F5 D9 77 28 50 C5 73",
			// 		"A4 9A 8C D8 A6 F2 74 23 DB 79 E4 4A B3 06 D9 42",
			// 	),
			// }, {
			// 	a: "D8 C9 62 3C 4D BF E8 14 CE 7C 1C 0C EA A9 59 DB", // h4
			// 	b: "A9 75 7B 81 47 95 6E 90 55 B8 A3 3D E8 9F 42 FC", // c1
			// 	c: bitstrings.BitSum128(
			// 		"A4 9A 8C D8 A6 F2 74 23 DB 79 E4 4A B3 06 D9 42",
			// 		"09 FE 3F 6A 83 3C 21 B3 90 27 D0 20 6A 84 E1 5A",
			// 	),
			// }, { //c2
			// 	a: "A5 E1 F1 95 33 3E 14 82 96 99 31 BF BE 6D FD 43", // h5
			// 	b: "80 75 D2 21 2B F9 FD 5B D3 F7 06 9A AD C1 6B 39", //c2
			// 	c: bitstrings.BitSum128(
			// 		"09 FE 3F 6A 83 3C 21 B3 90 27 D0 20 6A 84 E1 5A",
			// 		"B5 DA 26 BB 00 EB A8 04 35 D7 97 6B C6 B5 46 4D",
			// 	),
			// }, { //c3
			// 	a: "B4 CA 80 8C AC CF B3 F9 17 24 E4 8A 2C 7E E9 D2", // h6
			// 	b: "49 7A B1 59 15 A6 BA 85 93 6B 5D 0E A9 F6 85 1C", //c3
			// 	c: bitstrings.BitSum128(
			// 		"B5 DA 26 BB 00 EB A8 04 35 D7 97 6B C6 B5 46 4D",
			// 		"DD 1C 0E EE F7 83 C8 EB 2A 33 F3 58 D7 23 0E E5",
			// 	),
			// }, { //c4
			// 	a: "72 90 8F C0 74 E4 69 E8 90 1B D1 88 EA 91 C3 31", // h7
			// 	b: "C6 0C 14 D4 D3 F8 83 D0 AB 94 42 06 95 C7 6D EB", // c4
			// 	c: bitstrings.BitSum128(
			// 		"DD 1C 0E EE F7 83 C8 EB 2A 33 F3 58 D7 23 0E E5",
			// 		"89 6C E1 08 32 EB EA F9 06 9F 3F 73 76 59 4D 40",
			// 	),
			// }, { // c5
			// 	a: "23 CA 27 15 B0 2C 68 31 3B FD AC B3 9E 4D 0F B8", //h8
			// 	b: "2C 75 52 00 00 00 00 00 00 00 00 00 00 00 00 00", //c5
			// 	c: bitstrings.BitSum128(
			// 		"89 6C E1 08 32 EB EA F9 06 9F 3F 73 76 59 4D 40",
			// 		"99 1A F5 C9 D0 80 F7 63 87 FE 64 9E 7C 93 C6 42",
			// 	),
			// }, { // len
			// 	a: "00 00 00 00 00 00 01 48 00 00 00 00 00 00 02 18",
			// 	b: "BC BC E6 C4 1A A3 55 A4 14 88 62 BF 64 BD 83 0D",
			// 	c: bitstrings.BitSum128(
			// 		"99 1A F5 C9 D0 80 F7 63 87 FE 64 9E 7C 93 C6 42",
			// 		"C0 C7 22 DB 5E 0B D6 DB 25 76 73 83 3D 56 71 28",
			// 	),
		},
	}
	for _, td := range tt {
		t.Run(
			fmt.Sprintf("%s * %s -> %s", td.a, td.b, td.c),
			func(t *testing.T) {
				a_b, err := hex.DecodeString(td.a)
				if err != nil {
					t.Fatalf("err a_b: %s", err)
				}
				b_b, err := hex.DecodeString(td.b)
				if err != nil {
					t.Fatalf("err b_b: %s", err)
				}
				c_b, err := hex.DecodeString(td.c)
				if err != nil {
					t.Fatalf("err c_b: %s", err)
				}
				v1_upper := binary.BigEndian.Uint64(a_b[:8])
				v1_lower := binary.BigEndian.Uint64(a_b[8:])
				v2_upper := binary.BigEndian.Uint64(b_b[:8])
				v2_lower := binary.BigEndian.Uint64(b_b[8:])
				v3_upper := binary.BigEndian.Uint64(c_b[:8])
				v3_lower := binary.BigEndian.Uint64(c_b[8:])
				if upper, lower := multiplication.MultiplyUint128(v1_upper, v1_lower, v2_upper, v2_lower); upper != v3_upper || lower != v3_lower {
					t.Fatalf("\nGot:  %0.64b, %0.64b, \nWant: %0.64b, %0.64b", upper, lower, v3_upper, v3_lower)
				}
			},
		)
	}
}
