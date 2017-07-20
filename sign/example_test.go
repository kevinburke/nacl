package sign_test

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/kevinburke/nacl/sign"
)

func Example() {
	pubkey, privkey, err := sign.Keypair(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	message := []byte("Are you taking notes on a criminal conspiracy?")
	signedMessage := sign.Sign(message, privkey)
	result := sign.Verify(signedMessage, pubkey)
	fmt.Println(result)
	// Output: true
}
