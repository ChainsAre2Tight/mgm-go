package interfaces

// Represents a bit string of a certain length
type BitString interface {
	Length() int
	Bytes() []byte
	Upper() uint64
	Lower() uint64
}

type Encryptor interface {
	Encrypt(
		key string,
		associatedData string,
		plaintext string,
	) (
		nonce BitString,
		ciphertext string,
		mac string,
		err error,
	)
}
