package stream

import (
	"github.com/kevinburke/nacl"
	"golang.org/x/crypto/salsa20"
	"golang.org/x/crypto/salsa20/salsa"
)

const NonceSize = 24

// Stream produces a l-byte stream as a function of a secret key k and a nonce
// n. Note that it is the caller's responsibility to ensure the uniqueness of
// nonces—for example, by using nonce 1 for the first message, nonce 2 for the
// second message, etc. Nonces are long enough that randomly generated nonces
// have negligible risk of collision.
func Stream(l int, nonce nacl.Nonce, key nacl.Key) []byte {
	subKey, counter := nacl.Setup(nonce, key)
	out := make([]byte, l)
	salsa.XORKeyStream(out, out, counter, subKey)
	return out
}

// XOR encrypts a message m using a secret key k and a nonce n. XOR returns
// the ciphertext c. Note that it is the caller's responsibility to ensure the
// uniqueness of nonces—for example, by using nonce 1 for the first message,
// nonce 2 for the second message, etc. Nonces are long enough that randomly
// generated nonces have negligible risk of collision.
//
// Note also that encrypting a message with XOR does not protect against
// that message being tampered with in transit. To detect tampering, combine
// encryption with an authenticator, like the one provided by nacl/secretbox.
func XOR(message []byte, nonce nacl.Nonce, key nacl.Key) []byte {
	out := make([]byte, len(message))
	salsa20.XORKeyStream(out, message, nonce[:], key)
	return out
}
