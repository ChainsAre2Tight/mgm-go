package additionaldata

import (
	"encoding/binary"

	"github.com/ChainsAre2Tight/mgm-go/internal/multiplication"
	"github.com/ChainsAre2Tight/mgm-go/internal/utils"
)

func ComputeADMAC(
	encryptorFunc func(uint64, uint64) (uint64, uint64),
	nonceUpper, nonceLower uint64,
	associatedData []byte,
) (
	macUpper, macLower uint64,
	counterUpper, counterLower uint64,
) {
	nonceUpper |= 1 << 63

	counterUpper, counterLower = encryptorFunc(nonceUpper, nonceLower)

	if len(associatedData) == 0 {
		return 0, 0, counterUpper, counterLower
	}

	var upper, lower uint64
	for len(associatedData) >= 16 {
		upper = binary.BigEndian.Uint64(associatedData[:8])
		lower = binary.BigEndian.Uint64(associatedData[8:16])
		associatedData = associatedData[16:]

		h_upper, h_lower := encryptorFunc(counterUpper, counterLower)
		counterUpper++

		u, l := multiplication.MultiplyUint128(upper, lower, h_upper, h_lower)
		macUpper ^= u
		macLower ^= l
	}

	if len(associatedData) == 0 {
		return macUpper, macLower, counterUpper, counterLower
	}

	upper, lower = utils.BytesToUint64WithPadding(associatedData)

	h_upper, h_lower := encryptorFunc(counterUpper, counterLower)
	counterUpper++

	u, l := multiplication.MultiplyUint128(upper, lower, h_upper, h_lower)
	macUpper ^= u
	macLower ^= l

	return macUpper, macLower, counterUpper, counterLower
}
