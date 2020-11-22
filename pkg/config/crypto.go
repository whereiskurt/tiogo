package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

func (c *Config) decrypt() {
	if c.CryptoKey == "" {
		return
	}
	c.CryptoKey = strings.Repeat(c.CryptoKey, (16/len(c.CryptoKey) + 1))[:16]
	if c.VM.AccessKey != "" {
		raw, err := base64.StdEncoding.DecodeString(c.VM.AccessKey)
		if err != nil {
			c.VM.Log.Fatalf("fatal: cannot decode URL bytes")
		}
		dec, err := Decrypt(raw, []byte(c.CryptoKey))
		if err != nil {
			c.VM.Log.Fatalf("error: cannot decrypt URL with key")
			c.VM.Log.Fatalf("fatal: your cryptokey [--key] is likely wrong.")
		}
		c.VM.AccessKey = string(dec)
	}

	if c.VM.SecretKey != "" {
		raw, err := base64.StdEncoding.DecodeString(c.VM.SecretKey)
		if err != nil {
			c.VM.Log.Fatalf("fatal: cannot decode URL bytes")
		}
		dec, err := Decrypt(raw, []byte(c.CryptoKey))
		if err != nil {
			c.VM.Log.Fatalf("error: cannot decrypt URL with key")
			c.VM.Log.Fatalf("fatal: your cryptokey [--key] is likely wrong.")
		}
		c.VM.SecretKey = string(dec)
	}

	return
}

func (c *Config) encrypt() {
	if c.CryptoKey == "" {
		c.VM.Log.Fatalf("fatal: cannot encrypt config without --key.")
		return
	}
	c.CryptoKey = strings.Repeat(c.CryptoKey, (16/len(c.CryptoKey) + 1))[:16]

	r, err := Encrypt([]byte(c.VM.AccessKey), []byte(c.CryptoKey))
	if err != nil {
		c.VM.Log.Fatalf("fatal: cannot encrypt VM.AccessKey")
	}
	c.VM.AccessKey = string(base64.StdEncoding.EncodeToString([]byte(r)))

	r, err = Encrypt([]byte(c.VM.SecretKey), []byte(c.CryptoKey))
	if err != nil {
		c.VM.Log.Fatalf("fatal: cannot encrypt VM.SecretKey")
	}
	c.VM.SecretKey = string(base64.StdEncoding.EncodeToString([]byte(r)))

}

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
