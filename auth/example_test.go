// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package auth_test

import (
	"fmt"

	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/auth"
)

func Example() {
	// Load your secret key from a safe place and reuse it across multiple
	// Sum/Verify calls. (Obviously don't use this example key for anything
	// real.) If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	key, err := nacl.Load("6368616e676520746869732070617373776f726420746f206120736563726574")
	if err != nil {
		panic(err)
	}

	mac := auth.Sum([]byte("hello world"), key)
	fmt.Printf("%x\n", *mac)
	result := auth.Verify(mac, []byte("hello world"), key)
	fmt.Println(result)
	// Output: eca5a521f3d77b63f567fb0cb6f5f2d200641bc8dada42f60c5f881260c30317
	// true
}
