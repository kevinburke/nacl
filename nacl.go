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
	"golang.org/x/crypto/salsa20/salsa"
)

// The software version.
const Version = "0.5"

// Size of a public or private key in bytes.
const KeySize = 32

// Size of a nonce in bytes.
const NonceSize = 24

// Key represents a private or public key for use in encryption or
// authentication. A key should be random bytes and *not* simply 32 characters
// in the visible ASCII set.
type Key *[KeySize]byte

// Nonce is an arbitrary value that should be used only once per (sender,
// receiver) pair. For example, the lexicographically smaller public key can
// use nonce 1 for its first message to the other key, nonce 3 for its second
// message, nonce 5 for its third message, etc., while the lexicographically
// larger public key uses nonce 2 for its first message to the other key, nonce
// 4 for its second message, nonce 6 for its third message, etc. Nonces are long
// enough that randomly generated nonces have negligible risk of collision.
type Nonce *[NonceSize]byte

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
	if len(keyBytes) != KeySize {
		return nil, fmt.Errorf("nacl: incorrect key length: %d", len(keyBytes))
	}
	key := new([KeySize]byte)
	copy(key[:], keyBytes)
	return key, nil
}

// Load decodes a 128-byte hex string into a Key. A hex key is suitable for
// representation in a configuration file. You can generate one by running
// nacl/sign.GenerateKey(nil).
func Load64(hexkey string) (*[64]byte, error) {
	if len(hexkey) != 128 {
		return nil, fmt.Errorf("nacl: incorrect hex key length: %d, should be 64", len(hexkey))
	}
	keyBytes, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != 64 {
		return nil, fmt.Errorf("nacl: incorrect key length: %d", len(keyBytes))
	}
	key := new([64]byte)
	copy(key[:], keyBytes)
	return key, nil
}

// NewKey returns a new Key with cryptographically random data. NewKey panics if
// we could not read the correct amount of random data into key.
func NewKey() Key {
	key := new([KeySize]byte)
	randombytes.MustRead(key[:])
	return key
}

// NewNonce returns a new Nonce with cryptographically random data. It panics if
// we could not read the correct amount of random data into nonce.
func NewNonce() Nonce {
	nonce := new([NonceSize]byte)
	randombytes.MustRead(nonce[:])
	return nonce
}

// Verify returns true if and only if a and b have equal contents. The time
// taken is a function of the length of the slices and is independent of the
// contents. If an attacker controls the length of a, they may be able to
// determine the length of b (and vice versa).
func Verify(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}

// Verify16 returns true if and only if a and b have equal contents, without
// leaking timing information.
func Verify16(a, b *[16]byte) bool {
	if a == nil || b == nil {
		panic("nacl: nil input")
	}
	return subtle.ConstantTimeCompare(a[:], b[:]) == 1
}

// Verify32 returns true if and only if a and b have equal contents, without
// leaking timing information.
func Verify32(a, b *[KeySize]byte) bool {
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

// Setup produces a sub-key and Salsa20 counter given a nonce and key.
func Setup(nonce Nonce, key Key) (Key, *[16]byte) {
	// We use XSalsa20 for encryption so first we need to generate a
	// key and nonce with HSalsa20.
	var hNonce [16]byte
	copy(hNonce[:], nonce[:])
	var subKey [32]byte
	salsa.HSalsa20(&subKey, &hNonce, key, &salsa.Sigma)

	// The final 8 bytes of the original nonce form the new nonce.
	var counter [16]byte
	copy(counter[:], nonce[16:])
	return &subKey, &counter
}
