package sign_test

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/kevinburke/nacl/sign"
)

func Example() {
	// Create a public and private key pair.
	pubkey, privkey, err := sign.Keypair(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}
	message := []byte("Are you taking notes on a criminal conspiracy?")
	signedMessage := sign.Sign(message, privkey)
	// The first SignatureSize bytes will be the signature. The remaining bytes
	// will be the message.
	fmt.Printf("%s\n", signedMessage[sign.SignatureSize:])
	result := sign.Verify(signedMessage, pubkey)
	fmt.Println(result)
	// Output: Are you taking notes on a criminal conspiracy?
	// true
}
