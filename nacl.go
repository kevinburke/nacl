// Package nacl is a pure Go implementation of the NaCL cryptography library.
//
// Compared with the implementation in golang.org/x/crypto/nacl, this library
// offers all of the API's present in NaCL, as well as some utilities for
// generating and loading keys and nonces, and encrypting messages.
//
// NaCl's goal is to provide all of the core operations needed to build
// higher-level cryptographic tools, as well as to demonstrate how to implement
// these tools in Go.
//
// Compared with the equivalent packages in the Go standard library and x/crypto
// package, we replace some function calls with their equivalents in this
// package, and make more use of return values (versus writing to a byte array
// specified at stdin). Most functions should be compatible with their C/C++
// counterparts in the library here: https://nacl.cr.yp.to/. In many cases the
// tests are ported directly to this library.
package nacl

import (
	"crypto/sha512"
	"crypto/subtle"
	"encoding/hex"
	"fmt"

	"github.com/kevinburke/nacl/randombytes"
)

// The software version.
const Version = "0.2"

// Key represents a private or public key for use in encryption or
// authentication. A key should be random bytes and *not* simply 32 characters
// in the visible ASCII set.
type Key *[32]byte

// Nonce is an arbitrary value that should be used only once per (sender,
// receiver) pair. For example, the lexicographically smaller public key can
// use nonce 1 for its first message to the other key, nonce 3 for its second
// message, nonce 5 for its third message, etc., while the lexicographically
// larger public key uses nonce 2 for its first message to the other key, nonce
// 4 for its second message, nonce 6 for its third message, etc. Nonces are long
// enough that randomly generated nonces have negligible risk of collision.
type Nonce *[24]byte

// Load decodes a 64-byte hex string into a Key. A hex key is suitable for
// representation in a configuration file. You can generate one by running
// "openssl rand -hex 32".
func Load(hexkey string) (Key, error) {
	if len(hexkey) != 64 {
		return nil, fmt.Errorf("nacl: incorrect hex key length: %d, should be 64", len(hexkey))
	}
	keyBytes, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != 32 {
		return nil, fmt.Errorf("nacl: incorrect key length: %d", len(keyBytes))
	}
	key := new([32]byte)
	copy(key[:], keyBytes)
	return key, nil
}

// NewKey returns a new Key with cryptographically random data. NewKey panics if
// we could not read the correct amount of random data into key.
func NewKey() Key {
	key := new([32]byte)
	randombytes.MustRead(key[:])
	return key
}

// NewNonce returns a new Nonce with cryptographically random data. It panics if
// we could not read the correct amount of random data into nonce.
func NewNonce() Nonce {
	nonce := new([24]byte)
	randombytes.MustRead(nonce[:])
	return nonce
}

// Verify returns true if and only if a and b have equal contents.
func Verify(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}

// Verify16 returns true if and only if a and b have equal contents.
func Verify16(a, b *[16]byte) bool {
	if a == nil || b == nil {
		panic("nacl: nil input")
	}
	return subtle.ConstantTimeCompare(a[:], b[:]) == 1
}

// Verify32 returns true if and only if a and b have equal contents.
func Verify32(a, b *[32]byte) bool {
	if a == nil || b == nil {
		panic("nacl: nil input")
	}
	return subtle.ConstantTimeCompare(a[:], b[:]) == 1
}

// HashSize is the size, in bytes, of the result of calling Hash.
const HashSize = sha512.Size

// Hash hashes a message m.
//
// The crypto_hash function is designed to be usable as a strong component of
// DSA, RSA-PSS, key derivation, hash-based message-authentication codes,
// hash-based ciphers, and various other common applications. "Strong" means
// that the security of these applications, when instantiated with crypto_hash,
// is the same as the security of the applications against generic attacks.
// In particular, the crypto_hash function is designed to make finding
// collisions difficult.
//
// Hash is currently an implementation of SHA-512.
func Hash(m []byte) *[HashSize]byte {
	out := sha512.Sum512(m)
	return &out
}
