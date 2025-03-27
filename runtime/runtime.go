package runtime

import (
	"finally-main/mvc"
	"fmt"
	"os/exec"
	"strings"
	"time"
	//"regexp"
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

		x, c, c2, c1, _, _, _, _, _ := mvc.GetKline(strings.TrimSuffix(symbol, "-USDT-SWAP"), "15m")

		isbuy := c2 < 0.95 && c2 > 0.9 && c1 > 0.9 && c1 < 0.98 && c[x-1] > 0.0001
		if isbuy {

			fmt.Println(time.Now().Format("2006-1-2 15:04:02") + "-------------------------------买入--------------------------------->>>")

			cmd := exec.Command("python", "gorun.py", symbol, minute)
			res, _ := cmd.Output()
			fmt.Println(string(res))
		}
	}

}
