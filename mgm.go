package mgmgo

type Encryptor interface {
	Encrypt(
		key []byte,
		associatedData []byte,
		plaintext []byte,
	) (
		nonce []byte,
		ciphertext []byte,
		mac []byte,
		err error,
	)
}

type Decryptor interface {
	Decrypt(
		key []byte,
		nonce []byte,
		associatedData []byte,
		ciphertext []byte,
		mac []byte,
	) (
		plaintext []byte,
		err error,
	)
}
