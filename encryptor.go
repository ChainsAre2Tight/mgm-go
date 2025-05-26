package mgmgo

import (
	"bytes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"

	kuznechikgo "github.com/ChainsAre2Tight/kuznechik-go"
	ad "github.com/ChainsAre2Tight/mgm-go/internal/associateddata"
	"github.com/ChainsAre2Tight/mgm-go/internal/encryption"
)

var _ cipher.AEAD = (*MGM)(nil)

type EncryptorFunc func(uint64, uint64) (uint64, uint64)

type MGM struct {
	encryptorFunc EncryptorFunc
	keys          kuznechikgo.UintRoundKeys
}

func New(key []byte) (cipher.AEAD, error) {
	k, err := kuznechikgo.Schedule(key)
	if err != nil {
		return nil, fmt.Errorf("mgm.New: Error during keyschedule: %s", err)
	}
	res := MGM{
		keys: kuznechikgo.KeysToUints(k),
	}
	res.encryptorFunc = func(u, l uint64) (uint64, uint64) {
		return kuznechikgo.UintEncrypt(u, l, res.keys)
	}
	return &res, nil
}

func (m *MGM) NonceSize() int {
	return 16
}

func (m *MGM) Overhead() int {
	return 1 << 62
}

func (m *MGM) Open(dst []byte, nonce []byte, ciphertext []byte, additionalData []byte) ([]byte, error) {

	// todo delete
	lengthAuth := uint64(len(additionalData)) * 8
	lengthPlaintext := uint64(len(ciphertext)-16) * 8

	nonceUpper := binary.BigEndian.Uint64(nonce[:8])
	nonceLower := binary.BigEndian.Uint64(nonce[8:])

	macUpper, macLower, counterUpper, counterLower := ad.ComputeADMAC(m.encryptorFunc, nonceUpper, nonceLower, additionalData)

	head, tail := sliceForAppend(dst, len(ciphertext)-16)

	mac := make([]byte, 16)

	encryption.DecryptAndComputeMAC(
		m.encryptorFunc,
		nonceUpper, nonceLower,
		counterUpper, counterLower,
		macUpper, macLower,
		ciphertext[:len(ciphertext)-16],
		lengthAuth, lengthPlaintext,
		tail[:len(ciphertext)-16],
		mac,
	)

	if !bytes.Equal(mac, ciphertext[len(ciphertext)-16:]) {
		return head, fmt.Errorf("MACs differ")
	}

	return head, nil
}

// Taken from go/src/crypto/cipher/gcm.go
func sliceForAppend(in []byte, n int) (head, tail []byte) {
	if total := len(in) + n; cap(in) >= total {
		head = in[:total]
	} else {
		head = make([]byte, total)
		copy(head, in)
	}
	tail = head[len(in):]
	return
}

func (m *MGM) Seal(dst []byte, nonce []byte, plaintext []byte, additionalData []byte) []byte {

	// todo delete
	lengthAuth := uint64(len(additionalData)) * 8
	lengthPlaintext := uint64(len(plaintext)) * 8

	nonceUpper := binary.BigEndian.Uint64(nonce[:8])
	nonceLower := binary.BigEndian.Uint64(nonce[8:])

	macUpper, macLower, counterUpper, counterLower := ad.ComputeADMAC(m.encryptorFunc, nonceUpper, nonceLower, additionalData)

	head, tail := sliceForAppend(dst, len(plaintext)+16)

	encryption.EncryptAndComputeMAC(
		m.encryptorFunc,
		nonceUpper, nonceLower,
		counterUpper, counterLower,
		macUpper, macLower,
		plaintext,
		lengthAuth, lengthPlaintext,
		tail[:len(plaintext)],
		tail[len(plaintext):],
	)

	return head
}
