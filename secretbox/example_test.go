// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package secretbox_test

import (
	"fmt"

	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/secretbox"
)

func Example() {
	// Load your secret key from a safe place and reuse it across multiple
	// Seal calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	key, err := nacl.Load("6368616e676520746869732070617373776f726420746f206120736563726574")
	if err != nil {
		panic(err)
	}

	// This encrypts "hello world" and returns an encrypted message with a
	// random nonce prepended. You must use a different nonce for each message
	// you encrypt with the same key. Since the nonce here is 192 bits long, a
	// random value provides a sufficiently small probability of repeats.
	encrypted := secretbox.EasySeal([]byte("hello world"), key)

	// When you decrypt, you must use the same nonce and key you used to
	// encrypt the message. One way to achieve this is to store the nonce
	// alongside the encrypted message. Above, we stored the nonce in the first
	// 24 bytes of the encrypted text.
	decrypted, err := secretbox.EasyOpen(encrypted, key)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(decrypted))
	// Output: hello world
}
