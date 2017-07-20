// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sign

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"os"
	"strings"
	"testing"

	"github.com/kevinburke/nacl/sign/internal/edwards25519"
)

type zeroReader struct{}

func (zeroReader) Read(buf []byte) (int, error) {
	for i := range buf {
		buf[i] = 0
	}
	return len(buf), nil
}

func TestUnmarshalMarshal(t *testing.T) {
	pub, _, _ := Keypair(rand.Reader)

	var A edwards25519.ExtendedGroupElement
	var pubBytes [32]byte
	copy(pubBytes[:], pub)
	if !A.FromBytes(&pubBytes) {
		t.Fatalf("ExtendedGroupElement.FromBytes failed")
	}

	var pub2 [32]byte
	A.ToBytes(&pub2)

	if pubBytes != pub2 {
		t.Errorf("FromBytes(%v)->ToBytes does not round-trip, got %x\n", pubBytes, pub2)
	}
}

func TestSignVerify(t *testing.T) {
	var zero zeroReader
	public, private, _ := Keypair(zero)

	message := []byte("test message")
	signed := Sign(message, private)
	if !Verify(signed, public) {
		t.Errorf("valid signature rejected")
	}
	signature := signed[:SignatureSize]
	wrongMessage := append(signature, []byte("wrong message")...)
	if Verify(wrongMessage, public) {
		t.Errorf("signature of different message accepted")
	}
}

func TestCryptoSigner(t *testing.T) {
	var zero zeroReader
	public, private, _ := Keypair(zero)

	signer := crypto.Signer(private)

	publicInterface := signer.Public()
	public2, ok := publicInterface.(PublicKey)
	if !ok {
		t.Fatalf("expected PublicKey from Public() but got %T", publicInterface)
	}

	if !bytes.Equal(public, public2) {
		t.Errorf("public keys do not match: original:%x vs Public():%x", public, public2)
	}

	message := []byte("message")
	var noHash crypto.Hash
	signature, err := signer.Sign(zero, message, noHash)
	if err != nil {
		t.Fatalf("error from Sign(): %s", err)
	}

	if !public.Verify(signature) {
		t.Errorf("Verify failed on signature from Sign()")
	}
}

func TestGolden(t *testing.T) {
	// sign.input.gz is a selection of test cases from
	// https://ed25519.cr.yp.to/python/sign.input
	testDataZ, err := os.Open("testdata/sign.input.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer testDataZ.Close()
	testData, err := gzip.NewReader(testDataZ)
	if err != nil {
		t.Fatal(err)
	}
	defer testData.Close()

	scanner := bufio.NewScanner(testData)
	lineNo := 0

	for scanner.Scan() {
		lineNo++

		line := scanner.Text()
		parts := strings.Split(line, ":")
		if len(parts) != 5 {
			t.Fatalf("bad number of parts on line %d", lineNo)
		}

		privBytes, _ := hex.DecodeString(parts[0])
		pubKey, _ := hex.DecodeString(parts[1])
		msg, _ := hex.DecodeString(parts[2])
		sig, _ := hex.DecodeString(parts[3])

		if l := len(pubKey); l != PublicKeySize {
			t.Fatalf("bad public key length on line %d: got %d bytes", lineNo, l)
		}

		var priv [PrivateKeySize]byte
		copy(priv[:], privBytes)
		copy(priv[32:], pubKey)

		sig2 := Sign(msg, priv[:])
		if !bytes.Equal(sig, sig2[:]) {
			t.Errorf("different signature result on line %d: %x vs %x", lineNo, sig, sig2)
		}

		if !Verify(sig2, pubKey) {
			t.Errorf("signature failed to verify on line %d", lineNo)
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("error reading test data: %s", err)
	}
}

func BenchmarkKeyGeneration(b *testing.B) {
	var zero zeroReader
	for i := 0; i < b.N; i++ {
		if _, _, err := Keypair(zero); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSigning(b *testing.B) {
	var zero zeroReader
	_, priv, err := Keypair(zero)
	if err != nil {
		b.Fatal(err)
	}
	message := []byte("Hello, world!")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sign(message, priv)
	}
}

func BenchmarkVerification(b *testing.B) {
	var zero zeroReader
	pub, priv, err := Keypair(zero)
	if err != nil {
		b.Fatal(err)
	}
	message := []byte("Hello, world!")
	signature := Sign(message, priv)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Verify(signature, pub)
	}
}
