// generate-sign-keypair generates cryptographically secure public and private
// keys and prints hex encoded representations to stdout.
package main

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/kevinburke/nacl/sign"
)

func main() {
	silent := flag.Bool("silent", false, "Print the public key on the first line and private key on the second line only")
	flag.Parse()
	pub, priv, err := sign.Keypair(rand.Reader)
	if err != nil {
		panic(err)
	}
	pubhex := hex.EncodeToString(pub)
	privhex := hex.EncodeToString(priv)
	if *silent {
		fmt.Fprintf(os.Stdout, "%s\n%s\n", pubhex, privhex)
	} else {
		fmt.Fprintf(os.Stdout, "public:  %s\nprivate: %s\n", pubhex, privhex)
	}
}
