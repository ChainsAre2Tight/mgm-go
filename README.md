# MGM: Аутентифицированное шифрование (RFC 9058) с использованием алгоритма Кузнечик

Проект реализует режим аутентифицированного шифрования MGM (RFC 9058) на языке Go. Для шифрования данных используется алгоритм **Кузнечик** (ГОСТ 34.12-2015), и ключ для шифрования должен иметь длину **256 бит (32 байта)**.

Модуль реализует инрерфейс cipher.AEAD

## Установка

```bash
go get github.com/ChainsAre2Tight/mgm-go
```

## Пример использования

Смотри раздел [примеры](https://github.com/ChainsAre2Tight/mgm-go/tree/master/examples)

### Шифрование

```GO
package main

import (
	"fmt"
	"log"

	mgmgo "github.com/ChainsAre2Tight/mgm-go"
)

func main() {
	key := make([]byte, 32)

	mgm, err := mgmgo.New(key)
	if err != nil {
		log.Fatalf("Error during encryptor creation: %s", err)
	}

	additionalData := []byte("your-associated-data")
	plaintext := []byte("your-message")
	nonce := make([]byte, mgm.NonceSize())

	result := mgm.Seal(plaintext[:0], nonce, plaintext, additionalData)

	fmt.Printf("Ciphertext: %x\n", result[:len(plaintext)])
	fmt.Printf("MAC: %x\n", result[len(plaintext):])
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

	// Ключ должен быть длиной 256 бит (64 hex символа)
	key := make([]byte, 32)
	nonce := make([]byte, 16)
	additionalData := []byte("your-associated-data")
	ciphertext, _ := hex.DecodeString("705007f92ecfb7cffaf6f009")
	mac, _ := hex.DecodeString("df70a7a01caf7e134f9d6613df9c06c2")

	mgm, err := mgmgo.New(key)
	if err != nil {
		log.Fatalf("Error during encryptor creation: %s", err)
	}

	result, err := mgm.Open(ciphertext[:0], nonce, append(ciphertext, mac...), additionalData)

	if err != nil {
		log.Fatalf("Error: MACs are different")
	}

	fmt.Printf("Plaintext: %s\n", string(result))
}
```