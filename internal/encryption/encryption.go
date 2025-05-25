package encryption

import (
	"encoding/binary"

	"github.com/ChainsAre2Tight/mgm-go/internal/multiplication"
	"github.com/ChainsAre2Tight/mgm-go/internal/utils"
)

func EncryptAndComputeMAC(
	encryptorFunc func(uint64, uint64) (uint64, uint64),
	nonceUpper, nonceLower uint64,
	counterUpper, counterLower uint64,
	macUpper, macLower uint64,
	plaintext []byte,
	lengthAuth, lengthPlaintext uint64,
	ciphertext []byte,
	mac []byte,
) {
	nonceUpper, nonceLower = encryptorFunc((nonceUpper<<1)>>1, nonceLower)

	var upper, lower uint64
	for len(plaintext) >= 16 {
		upper = binary.BigEndian.Uint64(plaintext[:8])
		lower = binary.BigEndian.Uint64(plaintext[8:16])
		plaintext = plaintext[16:]

		y_upper, y_lower := encryptorFunc(nonceUpper, nonceLower)
		nonceLower++

		upper ^= y_upper
		lower ^= y_lower

		binary.BigEndian.PutUint64(ciphertext[:8], upper)
		binary.BigEndian.PutUint64(ciphertext[8:16], lower)
		ciphertext = ciphertext[16:]

		h_upper, h_lower := encryptorFunc(counterUpper, counterLower)
		counterUpper++

		u, l := multiplication.MultiplyUint128(upper, lower, h_upper, h_lower)
		macUpper ^= u
		macLower ^= l
	}

	upper, lower = utils.UintsToBytesWithPadding(plaintext)

	if len(plaintext) > 0 {
		y_upper, y_lower := encryptorFunc(nonceUpper, nonceLower)

		upper ^= y_upper
		lower ^= y_lower

		// careful, those bytes may overlap with auth Tag
		binary.BigEndian.PutUint64(ciphertext[:8], upper)
		binary.BigEndian.PutUint64(ciphertext[8:16], lower)
		ciphertext = ciphertext[:len(plaintext)]

		if len(plaintext) > 8 {
			shift := uint64((16 - len(plaintext)) * 8)
			lower = (lower >> shift) << shift
		} else {
			lower = 0
			shift := uint64((8 - len(plaintext)) * 8)
			upper = (upper >> shift) << shift
		}

		h_upper, h_lower := encryptorFunc(counterUpper, counterLower)
		counterUpper++

		u, l := multiplication.MultiplyUint128(upper, lower, h_upper, h_lower)
		macUpper ^= u
		macLower ^= l
	}

	h_upper, h_lower := encryptorFunc(counterUpper, counterLower)

	// finalize
	u, l := multiplication.MultiplyUint128(lengthAuth, lengthPlaintext, h_upper, h_lower)
	macUpper ^= u
	macLower ^= l

	macUpper, macLower = encryptorFunc(macUpper, macLower)

	binary.BigEndian.PutUint64(mac[0:8], macUpper)
	binary.BigEndian.PutUint64(mac[8:16], macLower)
}
