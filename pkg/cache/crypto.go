package cache

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

//TODO: Write some use cases.

// Encrypt takes the plaintext unencrypted and the key (byte arrays) and performs default AES encryption
func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt takes the cyphered text and the key (byte arrays) and performs default AES decryption
func Decrypt(crypt []byte, key []byte) (out []byte, err error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return
	}

	nonceSize := gcm.NonceSize()
	if len(crypt) < nonceSize {
		err = fmt.Errorf("crypt text too short")
		return
	}

	nonce, crypt := crypt[:nonceSize], crypt[nonceSize:]
	out, err = gcm.Open(nil, nonce, crypt, nil)

	return
}
