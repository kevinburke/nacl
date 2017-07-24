// Package onetimeauth provides primitives for secret-key, single-message
// authentication.
//
// The onetimeauth function, viewed as a function of the message for a
// uniform random key, is designed to meet the standard notion of unforgeability
// after a single message. After the sender authenticates one message, an
// attacker cannot find authenticators for any other messages.
//
// The sender must not use onetimeauth to authenticate more than one message
// under the same key. Authenticators for two messages under the same key should
// be expected to reveal enough information to allow forgeries of authenticators
// on other messages.
//
// The selected primitive is poly1305, an authenticator specified in
// "Cryptography in NaCl", Section 9. This authenticator is proven to meet the
// standard notion of unforgeability after a single message.
//
// This package is interoperable with NaCL: https://nacl.cr.yp.to/onetimeauth.html
package onetimeauth

import (
	"github.com/kevinburke/nacl"
	"golang.org/x/crypto/poly1305"
)

// Size is the size, in bytes, of the result of a call to Sum.
const Size = poly1305.TagSize

// Sum generates an authenticator for m using a one-time key and puts the result
// (of length Size) into out.
// Authenticating two different messages with the same key allows an attacker to
// forge messages at will.
func Sum(m []byte, key nacl.Key) *[Size]byte {
	out := new([Size]byte)
	poly1305.Sum(out, m, key)
	return out
}

// Verify returns true if mac is a valid authenticator for m with the given
// key, without leaking timing information.
func Verify(mac *[Size]byte, m []byte, key nacl.Key) bool {
	return poly1305.Verify(mac, m, key)
}
