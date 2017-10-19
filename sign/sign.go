// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sign can be used to verify messages were signed with a given secret
// key. The implementation uses the Ed25519 signature algorithm. See
// https://nacl.cr.yp.to/sign.html.
package sign

// This code is a port of the public domain, “ref10” implementation of ed25519
// from SUPERCOP.

import (
	"crypto"
	"errors"
	"io"
	"strconv"

	"golang.org/x/crypto/ed25519"
)

const (
	// PublicKeySize is the size, in bytes, of public keys as used in this package.
	PublicKeySize = 32
	// PrivateKeySize is the size, in bytes, of private keys as used in this package.
	PrivateKeySize = 64
	// SignatureSize is the size, in bytes, of signatures generated and verified by this package.
	SignatureSize = 64
)

// PublicKey is the type of Ed25519 public keys.
type PublicKey ed25519.PublicKey

// PrivateKey is the type of Ed25519 private keys. It implements crypto.Signer.
type PrivateKey ed25519.PrivateKey

// Public returns the PublicKey corresponding to priv.
func (priv PrivateKey) Public() crypto.PublicKey {
	pub := ed25519.PrivateKey(priv).Public().(ed25519.PublicKey)
	return PublicKey(pub)
}

// Sign signs the given message with priv.
// Ed25519 performs two passes over messages to be signed and therefore cannot
// handle pre-hashed messages. Thus opts.HashFunc() must return zero to
// indicate the message hasn't been hashed. This can be achieved by passing
// crypto.Hash(0) as the value for opts.
func (priv PrivateKey) Sign(rand io.Reader, message []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	if opts.HashFunc() != crypto.Hash(0) {
		return nil, errors.New("sign: cannot sign hashed message")
	}

	out := Sign(message, priv)
	return out[:], nil
}

// Keypair generates a public/private key pair using entropy from rand.
// If rand is nil, crypto/rand.Reader will be used.
func Keypair(rand io.Reader) (publicKey PublicKey, privateKey PrivateKey, err error) {
	public, private, err := ed25519.GenerateKey(rand)
	if err != nil {
		return nil, nil, err
	}
	return PublicKey(public), PrivateKey(private), nil
}

// Sign signs the message with privateKey. The first SignatureSize bytes of the
// response will be the signature; the rest will be the message. It will panic
// if len(privateKey) is not PrivateKeySize.
func Sign(message []byte, privateKey PrivateKey) []byte {
	sig := ed25519.Sign(ed25519.PrivateKey(privateKey), message)
	response := make([]byte, SignatureSize+len(message))
	copy(response[:SignatureSize], sig)
	copy(response[SignatureSize:], message)
	return response
}

// Verify uses key to report whether signature is a valid signature of message.
// The first SignatureSize bytes of signature should be the signature; the
// remaining bytes are the message to verify.
func (key PublicKey) Verify(signature []byte) bool {
	if len(signature) < SignatureSize || signature[63]&224 != 0 {
		return false
	}
	return Verify(signature, key)
}

// Verify reports whether sig is a valid signature of message by publicKey. It
// will panic if len(publicKey) is not PublicKeySize.
func Verify(sig []byte, publicKey PublicKey) bool {
	if l := len(publicKey); l != PublicKeySize {
		panic("sign: bad public key length: " + strconv.Itoa(l))
	}

	if len(sig) < SignatureSize || sig[63]&224 != 0 {
		return false
	}
	msg := sig[SignatureSize:]
	sig = sig[:SignatureSize]

	return ed25519.Verify(ed25519.PublicKey(publicKey), msg, sig)
}
