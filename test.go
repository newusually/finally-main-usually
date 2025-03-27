package main

import (
	"finally-main/mvc"
	"fmt"
	"strings"
)

func main() {

	symbollist := mvc.Getsymbollist()
	for _, symbol := range symbollist {
		isbuy := mvc.GetTakerVolume(strings.TrimSuffix(symbol, "-USDT-SWAP"), "1H")
		fmt.Println(symbol, isbuy)

	}
}
