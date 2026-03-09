package main

import (
	"flag"
	"io"
	"os"

	"github.com/kevinburke/nacl"
	"github.com/kevinburke/nacl/box"
)

var w = os.Stdout

func main() {

	encrypt := flag.Bool("encrypt", false, "Encrypt (instead of decrypt)")
	pubhex := flag.String("pubhex", "", "Public key in hex format")
	privhex := flag.String("privhex", "", "Private key in hex format")
	flag.Parse()
	pub, err := nacl.Load(*pubhex)
	if err != nil {
		panic(err)
	}
	prv, err := nacl.Load(*privhex)
	if err != nil {
		panic(err)
	}
	buf, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	if *encrypt {
		w.Write(box.EasySeal(buf, pub, prv))
	} else {
		out, err := box.EasyOpen(buf, pub, prv)
		if err != nil {
			panic(err)
		}
		w.Write((out))
	}

}
