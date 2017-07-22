// Package randombytes implements helpers for reading random data.
package randombytes

import (
	"crypto/rand"
	"strconv"
)

// Read fills in with random data.
func Read(in []byte) (int, error) {
	return rand.Read(in)
}

// MustRead fills in entirely with random data, or panics.
func MustRead(in []byte) {
	n, err := Read(in)
	if err != nil {
		panic(err)
	}
	if n != len(in) {
		panic("did not read enough random data: only " + strconv.Itoa(n) + " bytes")
	}
}
