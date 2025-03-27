package runtime

import (
	"finally-main/mvc"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/markcheno/go-talib"
)

func Run1() {
	Run("1m")
}

func Run3() {
	Run("3m")
}

func Run5() {
	Run("5m")
}

func Run15() {
	Run("15m")
}

func Run1H() {
	Run("1H")
}

func Run2H() {
	Run("2H")
}

func Run4H() {
	Run("4H")
}

func Savecsvfinal() {
	//mvc.Savecsv("1m")
	//mvc.Savecsv("3m")
	//mvc.Savecsv("5m")
	mvc.Savecsv("15m")

}

func Run6H() {
	Run("6H")

}

func Run12H() {
	Run("12H")
}

func Run1D() {
	Run("1D")
}
func Run(minute string) {
	defer func() {
		if r := recover(); r != nil {
			// 处理异常
			fmt.Println("Exception caught:", r)
		}
	}()

	symbollist := mvc.Getsymbollist()
	for _, symbol := range symbollist {

		x, c, c2, c1, o, _, _, _, _ := mvc.GetKline(strings.TrimSuffix(symbol, "-USDT-SWAP"), minute)
		if x < 10 {
			continue
		}
		c3 := c[x-3] / o[x-3]
		c4 := c[x-4] / o[x-4]
		c5 := c[x-5] / o[x-5]
		c6 := c[x-6] / o[x-6]
		c7 := c[x-7] / o[x-7]
		diff, dea, _ := talib.Macd(c, 12, 26, 60)
		macd1 := 2 * (diff[x-1] - dea[x-1])
		macd2 := 2 * (diff[x-2] - dea[x-2])
		macd3 := 2 * (diff[x-3] - dea[x-3])

		buy1 := c2 < 0.95 && c2 > 0.9 && c1 > 0.9 && c1 < 0.98
		buy2 := c3 > 0.9 && c3 < 0.99 && c2 > 0.9 && c2 < 0.99 && c1 > 0.9 && c1 < 0.99
		buy3 := macd1 < 0 && macd2 < 0 && macd3 < 0 && c1 < 0.99 && c1 > 0.95
		buy4 := macd1 < macd2 && macd2 < macd3 && c1 < 0.99 && c1 > 0.95
		buy5 := c1 > 0.9 && c1 < 1 && c2 > 0.9 && c2 < 1 && c3 > 0.9 && c3 < 1 && c4 > 0.9 && c4 < 1 && c5 > 0.9 && c5 < 1 && c6 > 0.9 && c6 < 1 && c7 > 0.9 && c7 < 1

		isbuy := (buy1 || buy2 || buy3 || buy4 || buy5) && c[x-1] > 0.0001
		if isbuy {

			fmt.Println(time.Now().Format("2006-1-2 15:04:02") + "-------------------------------买入--------------------------------->>>")

			cmd := exec.Command("python", "gorun.py", symbol, minute)
			res, _ := cmd.Output()
			fmt.Println(string(res))
		}
	}

}
