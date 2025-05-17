# MGM: Аутентифицированное шифрование (RFC 9058) с использованием алгоритма Кузнечик

Проект реализует режим аутентифицированного шифрования MGM (RFC 9058) на языке Go. Для шифрования данных используется алгоритм **Кузнечик** (ГОСТ 34.12-2015), и ключ для шифрования должен иметь длину **256 бит (32 байта)**. Для генерации уникальных значений `nonce` используется потоковый генератор. Блочное шифрование/дешифрование производится параллельно.

## Установка

```bash
go get github.com/ChainsAre2Tight/mgm-go
```

## Пример использования
### Шифрование

```GO
package main

import (
	"fmt"
	"log"

	mgmgo "github.com/ChainsAre2Tight/mgm-go"
)

func main() {
	nonceGenerator := mgmgo.NewNonceGenerator()

	encryptor := mgmgo.NewEncryptor(nonceGenerator)

	// Ключ должен быть длиной 256 бит (64 hex символа)
	key := "8899AABBCCDDEEFF0011223344556677FEDCBA98765432100123456789ABCDEF"
	associatedData := []byte("your-associated-data")
	plaintext := []byte("your-message")

	nonce, ciphertext, mac, err := encryptor.Encrypt(key, associatedData, plaintext)

	if err != nil {
		log.Fatalf("Encryption failed: %s", err)
	}

	fmt.Printf("Nonce: %x\n", nonce)
	fmt.Printf("Ciphertext: %x\n", ciphertext)
	fmt.Printf("MAC: %x\n", mac)
}

```

### Расшифрование и проверка подлиности

```Go
package main

import (
	"encoding/hex"
	"fmt"
	"log"

	mgmgo "github.com/ChainsAre2Tight/mgm-go"
)

func main() {
	decryptor := mgmgo.NewDecryptor()

	// Ключ должен быть длиной 256 бит (64 hex символа)
	key := "8899AABBCCDDEEFF0011223344556677FEDCBA98765432100123456789ABCDEF"
	nonce, err := hex.DecodeString("00000000000000000000000000000001")
	if err != nil {
		log.Fatalf("Error during nonce decoding: %s", err)
	}
	associatedData := []byte("your-associated-data")
	ciphertext, err := hex.DecodeString("5c44d197d9aa123feb46d896")
	if err != nil {
		log.Fatalf("Error during ciphertext decoding: %s", err)
	}
	mac, err := hex.DecodeString("5f11114a5ee24fdd5085d6ca11a249fe")
	if err != nil {
		log.Fatalf("Error during mac decoding: %s", err)
	}

	plaintext, err := decryptor.Decrypt(key, nonce, associatedData, ciphertext, mac)

	if err != nil {
		log.Fatalf("Decryption failed: %s", err)
	}

	fmt.Printf("Plaintext: %s\n", string(plaintext))
}

```