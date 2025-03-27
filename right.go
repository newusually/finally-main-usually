package main

import (
	"finally-main/mvc"
	"fmt"
)

func main() {
	symbol := "GPT-USDT-SWAP"
	minute := "1H"

	//  计算

	x_eth, c, _, _, _, _, _, h, _ := mvc.GetKline(symbol, minute)
	_, buy, sell := mvc.GetBuySellVol(symbol, minute)
	day, ratiosValues := mvc.GetRatio(symbol, minute)

	for i := 3; i < len(day); i++ {
		today := day[i]

		fmt.Println(today, "close:--->>>", (c[x_eth-len(ratiosValues)+i]), "high:--->>>",
			(h[x_eth-len(ratiosValues)+i]), "\nbuy/sell:--->>>",
			buy[len(buy)-len(ratiosValues)+i]/sell[len(buy)-len(ratiosValues)+i],
			"result--->>>", buy[len(buy)-len(ratiosValues)+i]/sell[len(buy)-len(ratiosValues)+i]-buy[len(buy)-len(ratiosValues)+i-1]/sell[len(buy)-len(ratiosValues)+i-1],
			"ratios:--->>>", ratiosValues[i])
		fmt.Println("===========================================")

	}

}
