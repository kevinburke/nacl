package nacl_test

import (
	"encoding/base64"
	"fmt"

	"github.com/kevinburke/nacl"
)

func ExampleLoad() {
	// Don't use this key for anything real.
	// You can generate one by running openssl rand -hex 32.
	key, err := nacl.Load("6368616e676520746869732070617373776f726420746f206120736563726574")
	if err != nil {
		panic(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(key[:]))
	// Output: Y2hhbmdlIHRoaXMgcGFzc3dvcmQgdG8gYSBzZWNyZXQ=
}
