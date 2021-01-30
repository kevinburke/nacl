package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/kevinburke/nacl"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, `generate-key generates a key for use with symmetric encryption (auth, secretbox).

`)
		flag.PrintDefaults()
	}
	flag.Parse()
	key := nacl.NewKey()
	keyhex := hex.EncodeToString((*key)[:])
	fmt.Fprintf(os.Stdout, "%s\n", keyhex)
}
