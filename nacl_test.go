package nacl

import (
	"encoding/hex"
	"testing"
)

func TestHash(t *testing.T) {
	in := []byte("testing\n")
	out := Hash(in)
	want := "24f950aac7b9ea9b3cb728228a0c82b67c39e96b4b344798870d5daee93e3ae5931baae8c7cacfea4b629452c38026a81d138bc7aad1af3ef7bfd5ec646d6c28"
	if got := hex.EncodeToString(out[:]); got != want {
		t.Errorf("Hash(%q): got %q, want %q", in, got, want)
	}
}

func TestSecretKey(t *testing.T) {
	t.Parallel()
	_, err := Load("")
	if err == nil {
		t.Errorf("expected non-nil error, got nil")
	}

	if _, err := Load("wrong length"); err.Error() != "nacl: incorrect hex key length: 12, should be 64" {
		t.Errorf("expected wrong-length error, got %q", err)
	}

	_, err = Load("zzzzzz6e676520746869732070617373776f726420746f206120736563726574")
	if err == nil || err.Error() != "encoding/hex: invalid byte: U+007A 'z'" {
		t.Errorf("expected invalid hex error, got %v", err)
	}

	key, err := Load("6368616e676520746869732070617373776f726420746f206120736563726574")
	if err != nil {
		t.Fatal(err)
	}
	h := hex.EncodeToString(key[:])
	if h != "6368616e676520746869732070617373776f726420746f206120736563726574" {
		t.Errorf("could not roundtrip decoded key: %s", h)
	}
}
