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
	cryptorand "crypto/rand"
	"crypto/sha512"
	"errors"
	"io"
	"strconv"

	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/sign/internal/edwards25519"
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
type PublicKey []byte

// PrivateKey is the type of Ed25519 private keys. It implements crypto.Signer.
type PrivateKey []byte

// Public returns the PublicKey corresponding to priv.
func (priv PrivateKey) Public() crypto.PublicKey {
	publicKey := make([]byte, PublicKeySize)
	copy(publicKey, priv[32:])
	return PublicKey(publicKey)
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
	if rand == nil {
		rand = cryptorand.Reader
	}

	privateKey = make([]byte, PrivateKeySize)
	publicKey = make([]byte, PublicKeySize)
	_, err = io.ReadFull(rand, privateKey[:32])
	if err != nil {
		return nil, nil, err
	}

	digest := nacl.Hash(privateKey[:32])
	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64

	var A edwards25519.ExtendedGroupElement
	var hBytes [32]byte
	copy(hBytes[:], digest[:])
	edwards25519.GeScalarMultBase(&A, &hBytes)
	var publicKeyBytes [32]byte
	A.ToBytes(&publicKeyBytes)

	copy(privateKey[32:], publicKeyBytes[:])
	copy(publicKey, publicKeyBytes[:])

	return publicKey, privateKey, nil
}

// Sign signs the message with privateKey. The first SignatureSize bytes of the
// response will be the signature; the rest will be the message. It will panic
// if len(privateKey) is not PrivateKeySize.
func Sign(message []byte, privateKey PrivateKey) []byte {
	if l := len(privateKey); l != PrivateKeySize {
		panic("sign: bad private key length: " + strconv.Itoa(l))
	}

	h := sha512.New()
	h.Write(privateKey[:32])

	var digest1, messageDigest, hramDigest [64]byte
	var expandedSecretKey [32]byte
	h.Sum(digest1[:0])
	copy(expandedSecretKey[:], digest1[:])
	expandedSecretKey[0] &= 248
	expandedSecretKey[31] &= 63
	expandedSecretKey[31] |= 64

	h.Reset()
	h.Write(digest1[32:])
	h.Write(message)
	h.Sum(messageDigest[:0])

	var messageDigestReduced [32]byte
	edwards25519.ScReduce(&messageDigestReduced, &messageDigest)
	var R edwards25519.ExtendedGroupElement
	edwards25519.GeScalarMultBase(&R, &messageDigestReduced)

	var encodedR [32]byte
	R.ToBytes(&encodedR)

	h.Reset()
	h.Write(encodedR[:])
	h.Write(privateKey[32:])
	h.Write(message)
	h.Sum(hramDigest[:0])
	var hramDigestReduced [32]byte
	edwards25519.ScReduce(&hramDigestReduced, &hramDigest)

	var s [32]byte
	edwards25519.ScMulAdd(&s, &hramDigestReduced, &expandedSecretKey, &messageDigestReduced)

	response := make([]byte, SignatureSize+len(message))
	copy(response[:], encodedR[:])
	copy(response[32:], s[:])
	copy(response[SignatureSize:], message)

	return response
}

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

	var A edwards25519.ExtendedGroupElement
	var publicKeyBytes [32]byte
	copy(publicKeyBytes[:], publicKey)
	if !A.FromBytes(&publicKeyBytes) {
		return false
	}
	edwards25519.FeNeg(&A.X, &A.X)
	edwards25519.FeNeg(&A.T, &A.T)

	h := sha512.New()
	h.Write(sig[:32])
	h.Write(publicKey[:])
	h.Write(msg)
	var digest [64]byte
	h.Sum(digest[:0])

	var hReduced [32]byte
	edwards25519.ScReduce(&hReduced, &digest)

	var R edwards25519.ProjectiveGroupElement
	var b [32]byte
	copy(b[:], sig[32:])
	edwards25519.GeDoubleScalarMultVartime(&R, &hReduced, &A, &b)

	var checkR [32]byte
	R.ToBytes(&checkR)
	return nacl.Verify(sig[:32], checkR[:])
}
