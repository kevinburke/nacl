# go-nacl

[![GoDoc](https://godoc.org/github.com/kevinburke/nacl?status.svg)](https://godoc.org/github.com/kevinburke/nacl)

This is a pure Go implementation of the API's available in NaCL:
https://nacl.cr.yp.to. Compared with the implementation in
golang.org/x/crypto/nacl, this library offers *all* of the API's present in
NaCL, better compatibility with NaCL implementations written in other languages,
as well as some utilities for generating and loading keys and nonces, and
encrypting messages.

Many of them are simple wrappers around functions or libraries available in the
Go standard library, or in the golang.org/x/crypto package. Other code I copied
directly into this library with the appropriate LICENSE; if a function is longer
than, say, 5 lines, I didn't write it myself. There are no dependencies outside
of the standard library or golang.org/x/crypto.

The goal is to both show how to implement the NaCL functions in pure Go, and
to provide interoperability between messages encrypted/hashed/authenticated in
other languages, and available in Go.

Among other benefits, NaCL is designed to be misuse resistant and standardizes
on the use of 32 byte keys and 24 byte nonces everywhere. Several helpers are
present for generating keys/nonces and loading them from configuration, as well
as for encrypting messages. You can generate a key by running `openssl rand -hex
32` and use the helpers in your program like so:

```go
import "github.com/kevinburke/nacl"
import "github.com/kevinburke/nacl/secretbox"

func main() {
    key, err := nacl.Load("6368616e676520746869732070617373776f726420746f206120736563726574")
    if err != nil {
        panic(err)
    }
    encrypted := secretbox.EasySeal([]byte("hello world"), key)
    fmt.Println(base64.StdEncoding.EncodeToString(encrypted))
}
```

The package names match the primitives available in NaCL, with the `crypto_`
prefix removed. Some function names have been changed to match the Go
conventions.

### Installation

```
go get github.com/kevinburke/nacl
```

Or you can Git clone the code directly to $GOPATH/src/github.com/kevinburke/nacl.

### Who am I?

While you probably shouldn't trust random security code from the Internet,
I'm reasonably confident that this code is secure. I did not implement any
of the hard math (poly1305, XSalsa20, curve25519) myself - I call into
golang.org/x/crypto for all of those functions. I also ported over every test
I could find from the C/C++ code, and associated RFC's, and ensured that these
libraries passed those tests.

I'm [a contributor to the Go Standard Library and associated
tools][contributor], and I've also been paid to do [security
consulting][services] for startups, and [found security problems in consumer
sites][capital-one].

[contributor]: https://go-review.googlesource.com/q/owner:kev%2540inburke.com
[capital-one]: https://burke.services/capital-one-open-redirect.html
[services]: https://burke.services

### Errata

- The implementation of `crypto_sign` uses the `ref10` implementation of ed25519
from SUPERCOP, *not* the current implementation in NaCL. The difference is that
the entire 64-byte signature is prepended to the message; in the current version
of NaCL, separate bits are prepended and appended to the message.

- Compared with `golang.org/x/crypto/ed25519`, this library's Sign
implementation returns the message along with the signature, and Verify
expects the first 64 bytes of the message to be the signature. This simplifies
the API and matches the behavior of the ref10 implementation and other NaCL
implementations. Sign also flips the order of the message and the private key:
`Sign(message, privatekey)`, to match the NaCL implementation.

- Compared with `golang.org/x/crypto/nacl/box`, `Precompute` returns the shared
key instead of modifying the input. In several places the code was modified to
call functions that now exist in `nacl`.

- Compared with `golang.org/x/crypto/nacl/secretbox`, `Seal` and `Open`
call the `onetimeauth` package in this library, instead of calling
`golang.org/x/crypto/poly1305` directly.
