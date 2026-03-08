package mcclient

import (
	"crypto/aes"
	"crypto/rand"
	"testing"
)

func TestCFB8RoundTripInPlace(t *testing.T) {
	key := make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("rand key: %v", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("new cipher: %v", err)
	}

	plain := []byte("hello minecraft cfb8 in-place stream")
	ciphertext := make([]byte, len(plain))
	copy(ciphertext, plain)

	enc := newCFB8(block, key, true)
	enc.XORKeyStream(ciphertext, ciphertext)

	dec := newCFB8(block, key, false)
	dec.XORKeyStream(ciphertext, ciphertext)

	if string(ciphertext) != string(plain) {
		t.Fatalf("round trip mismatch: got %q want %q", ciphertext, plain)
	}
}
