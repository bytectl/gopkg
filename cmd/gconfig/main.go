package main

import (
	"flag"
	"fmt"

	"github.com/bytectl/gopkg/crypto/gconfig"
)

var (
	key   = flag.String("k", "", "加解密密钥")
	value = flag.String("v", "", "加解密的值")
	e     = flag.Bool("e", false, "是否加密")
)

func main() {
	flag.Parse()
	if *value == "" || *key == "" {
		flag.Usage()
		return
	}
	if *e {
		enc := gconfig.EncryptString(*value, []byte(*key))
		fmt.Println(enc)
	} else {
		dec, err := gconfig.DecryptString(*value, []byte(*key))
		if err != nil {
			fmt.Println(value, " decrypt error:", err)
			return
		}
		fmt.Println(dec)
	}
}
