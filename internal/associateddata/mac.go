package associateddata

import (
	"encoding/binary"
	"fmt"

	kuznechikgo "github.com/ChainsAre2Tight/kuznechik-go"
	"github.com/ChainsAre2Tight/mgm-go/internal/multiplication"
	"github.com/ChainsAre2Tight/mgm-go/internal/utils"
)

func ComputeADMAC(
	keys kuznechikgo.UintRoundKeys,
	nonceUpper, nonceLower uint64,
	associatedData []byte,
) (
	macUpper, macLower uint64,
	counterUpper, counterLower uint64,
	err error,
) {
	fail := func(err error) (uint64, uint64, uint64, uint64, error) {
		return 0, 0, 0, 0, fmt.Errorf("accociatedData.ComputeADMAC: %s", err)
	}
	nonceUpper |= 1 << 63

	counterUpper, counterLower, err = kuznechikgo.UintEncrypt(nonceUpper, nonceLower, keys)
	if err != nil {
		return fail(fmt.Errorf("initial nonce encryption: %s", err))
	}

	if len(associatedData) == 0 {
		return 0, 0, counterUpper, counterLower, nil
	}

	var upper, lower uint64
	for len(associatedData) >= 16 {
		upper = binary.BigEndian.Uint64(associatedData[:8])
		lower = binary.BigEndian.Uint64(associatedData[8:16])
		associatedData = associatedData[16:]

		h_upper, h_lower, err := kuznechikgo.UintEncrypt(counterUpper, counterLower, keys)
		if err != nil {
			return fail(fmt.Errorf("counter encryption: %s", err))
		}
		counterUpper++

		u, l := multiplication.MultiplyUint128(upper, lower, h_upper, h_lower)
		macUpper ^= u
		macLower ^= l
	}

	if len(associatedData) == 0 {
		return macUpper, macLower, counterUpper, counterLower, nil
	}

	upper, lower = utils.UintsToBytesWithPadding(associatedData)

	h_upper, h_lower, err := kuznechikgo.UintEncrypt(counterUpper, counterLower, keys)
	if err != nil {
		return fail(fmt.Errorf("counter encryption final: %s", err))
	}
	counterUpper++

	u, l := multiplication.MultiplyUint128(upper, lower, h_upper, h_lower)
	macUpper ^= u
	macLower ^= l

	return macUpper, macLower, counterUpper, counterLower, nil
}
