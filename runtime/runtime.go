package runtime

import (
	"finally-main/mvc"
	"fmt"
	"os/exec"
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

		x, c, c2, c1, o, v, _, _, _ := mvc.GetKline(symbol, minute)

		if x < 200 {
			continue
		} else if symbol == "USDC-USDT-SWAP" || symbol == "USTC-USDT-SWAP" {
			continue

		} else if c[x-3] < 0.001 {
			continue

		} else {
			// 改良 MACD 参数
			diff, dea, _ := talib.Macd(c, 10, 20, 8)
			macd := make([]float64, len(diff))
			for i := range diff {
				macd[i] = 2 * (diff[i] - dea[i])
			}

			// 趋势判断：日线 MACD 柱体连续 3 根绿柱，且长度递减
			trendCheck := macd[x-2] < 0 && macd[x-3] < 0 && macd[x-4] < 0 && macd[x-2] > macd[x-3] && macd[x-3] > macd[x-4] && macd[x-4] < macd[x-5]

			// 量能验证：近 5 日成交量均值较前 10 日萎缩 40% 以上
			avgVolLast5 := talib.Sma(v[x-5:], 5)[0]
			avgVolLast10 := talib.Sma(v[x-10:], 10)[0]
			volumeCheck := avgVolLast5 < avgVolLast10*0.6

			// K 线确认：出现 “早晨之星” 或 “底分型” 形态（实体阳线突破 5 日均线）
			sma5 := talib.Sma(c, 5)
			klineCheck := c[x-2] > o[x-2] && c[x-2] > sma5[x-2] && c[x-3] < o[x-3] && c[x-2] > o[x-3]

			// 绿柱缩短 + 成交量萎缩至地量 + 阳线反包前日阴线实体
			macdShortening := macd[x-2] > macd[x-3] && macd[x-3] < 0
			volumeLow := avgVolLast5 < avgVolLast10*0.2
			bullishEngulfing := c[x-2] > o[x-2] && c[x-2] > o[x-3] && o[x-2] < c[x-3]

			isbuy := trendCheck && volumeCheck && klineCheck && macdShortening && volumeLow && bullishEngulfing && c[x-1] > 0.0001 && c1 > 0.99 && c1 < 1.1 && c2 > 1.002 && c2 < 1.1

			if isbuy {

				log := "\n" + time.Now().Format("2006-1-2 15:04:02") +
					",symbol--->>" + symbol +
					"\n,c1--->>" + fmt.Sprintf("%.5f", c1) +
					",c2--->>" + fmt.Sprintf("%.5f", c2) +
					",minute--->>" + minute
				fmt.Println(log)
				mvc.GetWriter(log, minute)

				fmt.Println(time.Now().Format("2006-1-2 15:04:02") + "-------------------------------买入--------------------------------->>>" + ",symbol--->>>" + symbol + ",minute--->>" + minute)

				cmd := exec.Command("python", "gorun.py", symbol, minute)
				res, _ := cmd.Output()
				fmt.Println(string(res))
			} else {
				continue
			}

		}

	}
}
