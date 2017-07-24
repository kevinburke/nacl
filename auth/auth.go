// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package auth authenticates a message m using a secret key k.

The Sum function, viewed as a function of the message for a uniform random
key, is designed to meet the standard notion of unforgeability. This means
that an attacker cannot find authenticators for any messages not authenticated
by the sender, even if the attacker has adaptively influenced the messages
authenticated by the sender. For a formal definition see, e.g., Section 2.4
of Bellare, Kilian, and Rogaway, "The security of the cipher block chaining
message authentication code," Journal of Computer and System Sciences 61 (2000),
362–399; http://www-cse.ucsd.edu/~mihir/papers/cbc.html.

auth does not make any promises regarding "strong" unforgeability; perhaps
one valid authenticator can be converted into another valid authenticator for
the same message. NaCl also does not make any promises regarding "truncated
unforgeability."

This package is interoperable with NaCl: https://nacl.cr.yp.to/auth.html.
*/
package auth

import (
	"crypto/hmac"
	"crypto/sha512"

	"github.com/kevinburke/nacl"
)

// Size is the size, in bytes, of an authenticated digest.
const Size = 32

// Sum generates an authenticator for m using a secret key and returns the
// 32-byte digest.
func Sum(m []byte, key nacl.Key) *[Size]byte {
	mac := hmac.New(sha512.New, (*key)[:])
	mac.Write(m)
	sum := mac.Sum(nil) // first 256 bits of 512 bit sum
	out := new([32]byte)
	copy(out[:], sum[:32])
	return out
}

// Verify checks that digest is a correct authenticator of a message m under the
// secret key key. If not, the function returns false. Verify does not leak
// timing information.
func Verify(digest *[Size]byte, m []byte, key nacl.Key) bool {
	mac := hmac.New(sha512.New, (*key)[:])
	mac.Write(m)
	expectedMAC := mac.Sum(nil) // first 256 bits of 512 bit sum
	return hmac.Equal((*digest)[:], expectedMAC[:Size])
}
