# go-nacl

This is a pure Go implementation of the API's available in NaCL:
https://nacl.cr.yp.to. Compared with the implementation in
golang.org/x/crypto/nacl, this library offers *all* of the API's present in
NaCL, better compatibility with NaCL implementations written in other languages,
as well as some utilities for generating and loading keys and nonces, and
encrypting messages.

Many of them are simple wrappers around functions or libraries available in
the Go standard library, or in the golang.org/x/crypto package. There are no
dependencies outside of the standard library or golang.org/x/crypto.

The goal is to both show how to implement the NaCL functions in pure Go, and
to provide interoperability between messages encrypted/hashed/authenticated in
other languages, and available in Go.

Among other benefits, NaCL is designed to be misuse resistant and standardizes
on the use of 32 byte keys and 24 byte nonces everywhere. Several helpers are
present for generating keys/nonces and loading them from configuration. You can
generate a key by running `openssl rand -hex 32` and use the helpers in your
program like so:

```go
key, err := nacl.Load("6368616e676520746869732070617373776f726420746f206120736563726574")
if err != nil {
    panic(err)
}
nonce := nacl.NewNonce()
encrypted := secretbox.Seal(nonce[:], []byte("hello world"), nonce, key)
fmt.Println(base64.StdEncoding.EncodeToString(encrypted))
```

The package names match the primitives available in NaCL, with the `crypto_`
prefix removed. Some function names have been changed to match the Go
conventions.

### Installation

```
go get github.com/kevinburke/nacl
```

Or you can Git clone the code directly to $GOPATH/src/github.com/kevinburke/nacl.

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
