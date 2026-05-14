// Package crypto provides HIPAA-compliant AES-256-GCM encryption for Protected Health Information (PHI).
//
// All patient message content, notes, and medication data stored in the database
// are encrypted at rest using AES-256-GCM, which provides both confidentiality and
// authenticated integrity. The nonce is randomly generated per operation and prepended
// to the ciphertext, making each encrypted blob self-contained.
//
// HIPAA Technical Safeguard references:
//
//	45 CFR § 164.312(a)(2)(iv) – Encryption and decryption
//	45 CFR § 164.312(e)(2)(ii) – Encryption of data in transit and at rest
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// Encryptor holds an AES-256-GCM AEAD cipher, keyed once at startup.
type Encryptor struct {
	aead cipher.AEAD
}

// NewEncryptor creates an Encryptor from a 32-byte AES-256 key.
// Panics if the key length is incorrect – this is a programming error, not a runtime error.
func NewEncryptor(key []byte) (*Encryptor, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("crypto: AES-256 key must be 32 bytes, got %d", len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("crypto: failed to create AES cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("crypto: failed to create GCM wrapper: %w", err)
	}
	return &Encryptor{aead: aead}, nil
}

// Encrypt encrypts plaintext using AES-256-GCM and returns a base64-encoded string
// of the form: base64(nonce || ciphertext || gcm_tag).
// A fresh random 12-byte nonce is generated for every call.
func (e *Encryptor) Encrypt(plaintext string) (string, error) {
	nonce := make([]byte, e.aead.NonceSize()) // 12 bytes for GCM
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("crypto: nonce generation failed: %w", err)
	}

	// Seal appends ciphertext+tag to nonce, producing nonce||ciphertext||tag
	ciphertext := e.aead.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decodes a base64 blob produced by Encrypt and returns the original plaintext.
// Returns an error if the ciphertext has been tampered with (GCM authentication failure).
func (e *Encryptor) Decrypt(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("crypto: base64 decode failed: %w", err)
	}

	nonceSize := e.aead.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("crypto: ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := e.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// Do NOT wrap the underlying error message – it may contain timing info.
		return "", errors.New("crypto: decryption failed (authentication error or corrupted data)")
	}
	return string(plaintext), nil
}

// MustEncrypt encrypts or panics. Use only at startup for static values.
func (e *Encryptor) MustEncrypt(plaintext string) string {
	s, err := e.Encrypt(plaintext)
	if err != nil {
		panic(err)
	}
	return s
}

// GenerateKey generates a new cryptographically random 32-byte AES-256 key
// and returns it as a base64-encoded string. Run this once during setup:
//
//	go run -v ./scripts/genkey
func GenerateKey() (string, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", fmt.Errorf("crypto: failed to generate key: %w", err)
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
