package randombytes

import (
	"fmt"
	"testing"
)

func TestRead(t *testing.T) {
	var p [32]byte
	for j := 0; j < 100; j++ {
		n, err := Read(p[:])
		if err != nil {
			t.Fatal(err)
		}
		if n != 32 {
			t.Errorf("Read: expected to read 32 bytes, got %d", n)
		}
		v := uint64(0)
		for i := 0; i < 32; i++ {
			v += uint64(p[i])
			p[i] = 0
		}
		if v < 100 {
			t.Errorf("expected p to be filled with random bytes, got sum %d: %x", v, p)
		}
	}
}

func TestMustRead(t *testing.T) {
	var p [32]byte
	for j := 0; j < 100; j++ {
		MustRead(p[:])
		v := uint64(0)
		for i := 0; i < 32; i++ {
			v += uint64(p[i])
			p[i] = 0
		}
		fmt.Println("v", v)
		if v < 100 {
			t.Errorf("expected p to be filled with random bytes, got sum %d: %x", v, p)
		}
	}
}
