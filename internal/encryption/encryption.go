package encryption

import (
	"encoding/binary"
	"fmt"

	kuznechikgo "github.com/ChainsAre2Tight/kuznechik-go"
	"github.com/ChainsAre2Tight/mgm-go/internal/multiplication"
	"github.com/ChainsAre2Tight/mgm-go/internal/utils"
)

func EncryptAndComputeMAC(
	keys kuznechikgo.UintRoundKeys,
	nonceUpper, nonceLower uint64,
	counterUpper, counterLower uint64,
	macUpper, macLower uint64,
	plaintext []byte,
	lengthAuth, lengthPlaintext uint64,
) (
	ciphertext []byte,
	mac []byte,
	err error,
) {
	fail := func(err error) ([]byte, []byte, error) {
		return nil, nil, fmt.Errorf("ecnryption.EncryptAndComputeMAC: %s", err)
	}

	nonceUpper, nonceLower, err = kuznechikgo.UintEncrypt((nonceUpper<<1)>>1, nonceLower, keys)
	if err != nil {
		return fail(fmt.Errorf("initial nonce encryption: %s", err))
	}

	ciphertext = make([]byte, 0, (len(plaintext)%8)*8)

	var upper, lower uint64
	for len(plaintext) >= 16 {
		upper = binary.BigEndian.Uint64(plaintext[:8])
		lower = binary.BigEndian.Uint64(plaintext[8:16])
		plaintext = plaintext[16:]

		y_upper, y_lower, err := kuznechikgo.UintEncrypt(nonceUpper, nonceLower, keys)
		if err != nil {
			return fail(fmt.Errorf("h counter encryption: %s", err))
		}
		nonceLower++

		upper ^= y_upper
		lower ^= y_lower

		ciphertext = binary.BigEndian.AppendUint64(ciphertext, upper)
		ciphertext = binary.BigEndian.AppendUint64(ciphertext, lower)

		h_upper, h_lower, err := kuznechikgo.UintEncrypt(counterUpper, counterLower, keys)
		if err != nil {
			return fail(fmt.Errorf("h counter encryption: %s", err))
		}
		counterUpper++

		u, l := multiplication.MultiplyUint128(upper, lower, h_upper, h_lower)
		macUpper ^= u
		macLower ^= l
	}

	upper, lower = utils.UintsToBytesWithPadding(plaintext)

	if len(plaintext) > 0 {
		y_upper, y_lower, err := kuznechikgo.UintEncrypt(nonceUpper, nonceLower, keys)
		if err != nil {
			return fail(fmt.Errorf("y counter encryption final: %s", err))
		}

		upper ^= y_upper
		lower ^= y_lower

		ciphertext = binary.BigEndian.AppendUint64(ciphertext, upper)
		ciphertext = binary.BigEndian.AppendUint64(ciphertext, lower)
		ciphertext = ciphertext[:len(ciphertext)-16+len(plaintext)]

		if len(plaintext) > 8 {
			shift := uint64((16 - len(plaintext)) * 8)
			lower = (lower >> shift) << shift
		} else {
			lower = 0
			shift := uint64((8 - len(plaintext)) * 8)
			upper = (upper >> shift) << shift
		}

		h_upper, h_lower, err := kuznechikgo.UintEncrypt(counterUpper, counterLower, keys)
		if err != nil {
			return fail(fmt.Errorf("h counter encryption final: %s", err))
		}
		counterUpper++

		u, l := multiplication.MultiplyUint128(upper, lower, h_upper, h_lower)
		macUpper ^= u
		macLower ^= l
	}

	h_upper, h_lower, err := kuznechikgo.UintEncrypt(counterUpper, counterLower, keys)
	if err != nil {
		return fail(fmt.Errorf("h counter encryption final (length): %s", err))
	}
	// finalize
	u, l := multiplication.MultiplyUint128(lengthAuth, lengthPlaintext, h_upper, h_lower)
	macUpper ^= u
	macLower ^= l

	macUpper, macLower, err = kuznechikgo.UintEncrypt(macUpper, macLower, keys)
	if err != nil {
		return fail(fmt.Errorf("h counter encryption final (transform): %s", err))
	}

	mac = make([]byte, 0, 16)
	mac = binary.BigEndian.AppendUint64(mac, macUpper)
	mac = binary.BigEndian.AppendUint64(mac, macLower)

	return ciphertext, mac, nil
}
