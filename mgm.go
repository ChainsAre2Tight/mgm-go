package mgmgo

type Encryptor interface {
	Encrypt(
		key string,
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
		key string,
		nonce []byte,
		associatedData []byte,
		ciphertext []byte,
		mac []byte,
	) (
		plaintext []byte,
		err error,
	)
}
