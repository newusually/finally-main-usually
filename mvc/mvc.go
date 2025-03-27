package mvc

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/markcheno/go-talib"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func Getsymbollist() []string {

	defer func() {
		if err := recover(); err != nil {
			log.Printf("[ERROR]程序异常: %v\n堆栈:%s", err, debug.Stack())
			os.Exit(1) // 非0退出码
		}
	}()
	time.Sleep(time.Millisecond)

	// 获取当前时间 或者使用 time.Date(year, month, ...)
	t := time.Now()
	timeStamp := t.Unix()
	client := &http.Client{

		Transport: &http.Transport{

			TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
		},
	}
	req, err := http.NewRequest("GET", "https://www.okx.com/api/v5/public/instruments?instType=SWAP", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("authority", "www.okx.com")
	req.Header.Set("timeout", "10000")
	req.Header.Set("x-cdn", "https://static.okx.com")
	req.Header.Set("devid", "8ccf140e-e4ab-4a46-8582-738445cad57c")
	req.Header.Set("accept-language", "zh-CN")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36 SE 2.X MetaSr 1.0")
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-utc", "8")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("app-type", "web")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("referer", "https://www.okx.com/trade-swap/btc-usdt-swap")
	req.Header.Set("cookie", "locale=zh_CN; defaultLocale=zh_CN; _gcl_au=1.1.2108489782."+strconv.Itoa(int(timeStamp))+"; _ga=GA1.2.1119126875."+strconv.Itoa(int(timeStamp))+"; _gid=GA1.2.1241734301."+strconv.Itoa(int(timeStamp))+"; amp_56bf9d=4KsK9IqNpGzx7Sx3l_9DPR...1g2ehgo0j.1g2ej7b4i.4.0.4")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)

	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)

	}
	src := string(bodyText)

	data := gjson.Get(src, "data").Array()
	regex := regexp.MustCompile(`^.*-USDT-SWAP$`)
	var validInstIDs []string

	for _, item := range data {
		// 解析基础参数
		instID := item.Get("instId").String()
		if !regex.MatchString(instID) {
			continue
		}

		// 转换为数值类型
		ctMult, _ := strconv.ParseFloat(item.Get("ctMult").String(), 64)
		ctVal, _ := strconv.ParseFloat(item.Get("ctVal").String(), 64)
		lever, _ := strconv.ParseFloat(item.Get("lever").String(), 64)

		// 杠杆上限处理
		lever = math.Min(lever, 50)

		// 获取标记价格
		x, c, _, _, _, _, _, _, _ := GetKline(instID, "1m")

		if x > 50 {
			price := c[x-1]
			// 计算初始保证金
			initialMargin := (ctVal * ctMult * price) / lever

			// 筛选条件
			if initialMargin < 1 {
				//fmt.Printf("符合条件: %-20s 保证金=%.4f USD\n", instID, initialMargin)
				validInstIDs = append(validInstIDs, instID)
			}
		}

	}

	return validInstIDs
}

func Getprice(symbol string) float64 {

	time.Sleep(time.Millisecond)

	// 获取当前时间 或者使用 time.Date(year, month, ...)
	t := time.Now()
	timeStamp := t.Unix()
	client := &http.Client{

		Transport: &http.Transport{

			TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
		},
	}
	req, err := http.NewRequest("GET", "https://www.okx.com/priapi/v5/market/mult-tickers?instIds="+symbol+"&t="+strconv.Itoa(int(timeStamp)), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("authority", "www.okx.com")
	req.Header.Set("timeout", "10000")
	req.Header.Set("x-cdn", "https://static.okx.com")
	req.Header.Set("devid", "8ccf140e-e4ab-4a46-8582-738445cad57c")
	req.Header.Set("accept-language", "zh-CN")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36 SE 2.X MetaSr 1.0")
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-utc", "8")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("app-type", "web")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("referer", "https://www.okx.com/trade-swap/btc-usdt-swap")
	req.Header.Set("cookie", "locale=zh_CN; defaultLocale=zh_CN; _gcl_au=1.1.2108489782."+strconv.Itoa(int(timeStamp))+"; _ga=GA1.2.1119126875."+strconv.Itoa(int(timeStamp))+"; _gid=GA1.2.1241734301."+strconv.Itoa(int(timeStamp))+"; amp_56bf9d=4KsK9IqNpGzx7Sx3l_9DPR...1g2ehgo0j.1g2ej7b4i.4.0.4")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)

	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)

	}
	src := string(bodyText)
	//fmt.Println(src)
	//"instType":"SWAP","instId":"GODS-USDT-SWAP","last":"0.7141","lastSz":"13",
	lastprice := gjson.Get(src, "data.#.last").Array()
	//fmt.Println((lastprice[0]).Float())
	fmt.Println(symbol, lastprice)

	return lastprice[0].Float()

}

func Getsymbols() []gjson.Result {
	time.Sleep(time.Millisecond)

	// 获取当前时间 或者使用 time.Date(year, month, ...)
	t := time.Now()
	timeStamp := t.Unix()
	client := &http.Client{

		Transport: &http.Transport{

			TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
		},
	}
	req, err := http.NewRequest("GET", "https://www.okx.com/priapi/v5/market/tickers?instType=SWAP&t="+strconv.Itoa(int(timeStamp))+"&instType=SWAP", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("authority", "www.okx.com")
	req.Header.Set("timeout", "10000")
	req.Header.Set("x-cdn", "https://static.okx.com")
	req.Header.Set("devid", "8ccf140e-e4ab-4a46-8582-738445cad57c")
	req.Header.Set("accept-language", "zh-CN")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36 SE 2.X MetaSr 1.0")
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-utc", "8")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("app-type", "web")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("referer", "https://www.okx.com/trade-swap/btc-usdt-swap")
	req.Header.Set("cookie", "locale=zh_CN; defaultLocale=zh_CN; _gcl_au=1.1.2108489782."+strconv.Itoa(int(timeStamp))+"; _ga=GA1.2.1119126875."+strconv.Itoa(int(timeStamp))+"; _gid=GA1.2.1241734301."+strconv.Itoa(int(timeStamp))+"; amp_56bf9d=4KsK9IqNpGzx7Sx3l_9DPR...1g2ehgo0j.1g2ej7b4i.4.0.4")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)

	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)

	}
	src := string(bodyText)
	//fmt.Println(src)
	//"instType":"SWAP","instId":"GODS-USDT-SWAP","last":"0.7141","lastSz":"13",
	instId := gjson.Get(src, "data.#.instId").Array()
	//fmt.Println(instId)
	return instId
}

func PlayMusic() {
	// 打开MP3文件
	f, err := os.Open("brave heart.wav")
	if err != nil {
		panic(err)
	}

	// 创建MP3解码器
	streamer, format, err := wav.Decode(f)
	if err != nil {
		panic(err)
	}
	defer streamer.Close()

	// 初始化扬声器
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// 播放音频
	speaker.Play(streamer)

	// 等待音频播放完毕
	select {}
}

func CalculatePSY(closes []float64, n int) float64 {
	var N1 float64
	start := len(closes) - n
	if start < 1 {
		start = 1
	}
	for i := start; i < len(closes); i++ {
		if closes[i] > closes[i-1] {
			N1++
		}
	}
	PSY := (N1 / float64(n)) * 100
	return PSY
}

func CalculateVR(closes []float64, vols []float64, n int) float64 {
	var AVS, BVS, CV float64
	start := len(closes) - n
	if start < 1 {
		start = 1
	}
	for i := start; i < len(closes); i++ {
		if closes[i] > closes[i-1] {
			AVS += vols[i]
		} else if closes[i] < closes[i-1] {
			BVS += vols[i]
		} else {
			CV += vols[i]
		}
	}
	VR := (AVS + 0.5*CV) / (BVS + 0.5*CV) * 100
	return VR
}

func CalculateAR(highs, lows, opens []float64, n int) float64 {
	var sumHighOpen, sumOpenLow float64
	for i := 0; i < n; i++ {
		sumHighOpen += highs[len(highs)-i-1] - opens[len(opens)-i-1]
		sumOpenLow += opens[len(opens)-i-1] - lows[len(lows)-i-1]
	}
	if sumOpenLow == 0 {
		return 0
	}
	return (sumHighOpen / float64(n)) / (sumOpenLow / float64(n)) * 100
}

func CalculateBR(highs, lows, closes []float64, n int) float64 {
	var sumHighPrevClose, sumPrevCloseLow float64
	for i := 1; i <= n; i++ {
		sumHighPrevClose += highs[len(highs)-i-1] - closes[len(closes)-i-1]
		sumPrevCloseLow += closes[len(closes)-i-1] - lows[len(lows)-i-1]
	}
	if sumPrevCloseLow == 0 {
		return 0
	}
	return (sumHighPrevClose / float64(n)) / (sumPrevCloseLow / float64(n)) * 100
}

func CalculateBIAS(closes []float64, n int) []float64 {

	ma := talib.Sma(closes, n)
	bias := make([]float64, len(closes))

	for i := 0; i < len(closes); i++ {
		bias[i] = (closes[i] - ma[i]) / ma[i] * 100

	}
	return bias
}

func GetisBtcMacd(symbol string, minute string) bool {
	x, c, _, _, _, _, _, _, _ := GetKline(symbol, minute)
	//x, c, c[x-2] / o[x-2], c[x-1] / o[x-1], o, l, v,h
	if x < 500 {
		return false

	} else if symbol == "USDC-USDT-SWAP" || symbol == "USTC-USDT-SWAP" {
		return false

	} else {
		diff, dea, _ := talib.Macd(c, 12, 26, 60)
		macd1 := 2 * (diff[x-1] - dea[x-1])
		macd2 := 2 * (diff[x-2] - dea[x-2])
		macd3 := 2 * (diff[x-3] - dea[x-3])
		macd4 := 2 * (diff[x-4] - dea[x-4])

		if c[x-3] > 0.001 && macd1 > 0 && macd2 < 0 && macd3 < 0 && macd4 < 0 {

			return true
		}
		return false
	}
}

func GetisEthisBuy(minute string) bool {
	//_, c1, _, co1_btc, _, _, _ := GetKline("BTC-USDT-SWAP", minute)
	_, _, co2, co1, _, _, _, _, _ := GetKline("ETH-USDT-SWAP", minute)
	//diff1, dea1, _ := talib.Macd(c1, 12, 26, 60)

	if co1 < 0.9 || co2 < 0.9 {

		return true
	} else {
		return false
	}
}

func GetIsBuy(symbol string, minute string) bool {

	//x, c, c[x-2] / o[x-2], c[x-1] / o[x-1], o, l, v,h
	isbuy := GetisBtcMacd(symbol, minute)
	if isbuy {

		_, buy, sell := GetBuySellVol(symbol, minute)

		buy1 := buy[len(buy)-1]
		sell1 := sell[len(sell)-1]
		buy2 := buy[len(buy)-2]
		sell2 := sell[len(sell)-2]
		buy3 := buy[len(buy)-3]
		sell3 := sell[len(sell)-3]
		buy4 := buy[len(buy)-4]
		sell4 := sell[len(sell)-4]

		buys1 := (buy1 + buy2 + buy3) / (sell1 + sell2 + sell3)
		buys2 := (buy2 + buy3 + buy4) / (sell2 + sell3 + sell4)

		if buys1/buys2 > 1 && buys1/buys2 < 2 && buy1/sell1 > 1 && buy1/sell1 < 2 {

			log := "\n" + time.Now().Format("2006-1-2 15:04:02") +
				",symbol--->>" + symbol +
				"\n,buys1/buys2--->>" + fmt.Sprintf("%.5f", buys1/buys2) +
				",buys1-->>" + fmt.Sprintf("%.5f", buys1) +
				",buys2-->>" + fmt.Sprintf("%.5f", buys2) +
				"\n,buy1/sell1--->>" + fmt.Sprintf("%.5f", buy1/sell1) +
				",buy2/sell2--->>" + fmt.Sprintf("%.5f", buy2/sell2) +
				",buy3/sell3-->>" + fmt.Sprintf("%.5f", buy3/sell3) +
				",minute--->>" + minute
			fmt.Println(log)
			GetWriter(log, minute)
			return true
		}

		return false
	}
	return false
}

func GetBuySellVol(symbol string, minute string) ([]string, []float64, []float64) {
	time.Sleep(time.Millisecond * 10)

	// 获取当前时间 或者使用 time.Date(year, month, ...)
	t := time.Now()
	timeStamp := t.Unix()
	client := &http.Client{

		Transport: &http.Transport{

			TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
		},
	}

	req, err := http.NewRequest("GET", "https://www.okx.com/priapi/v5/rubik/public/stat/indicators?bar="+minute+"&instId="+symbol+"&unit=2&limit=1030&before=1725000000000&indicators=takerBuySellVol&t="+strconv.Itoa(int(timeStamp)), nil)
	if err != nil {
		panic(err)

	}
	req.Header.Set("authority", "www.okx.com")
	req.Header.Set("timeout", "10000")
	req.Header.Set("x-cdn", "https://static.okx.com")
	req.Header.Set("devid", "a5ceb850-4efb-4a3f-baff-21da4fce8858")
	req.Header.Set("accept-language", "zh-CN")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36 SE 2.X MetaSr 1.0")
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-utc", "8")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("app-type", "web")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("referer", "https://www.okx.com/trade-swap/btc-usdt-swap")
	req.Header.Set("cookie", "locale=zh_CN; defaultLocale=zh_CN; _gcl_au=1.1.1520807996."+strconv.Itoa(int(timeStamp))+"; _ga=GA1.2.1752370991."+strconv.Itoa(int(timeStamp))+"; _gid=GA1.2.650161560."+strconv.Itoa(int(timeStamp))+"; amp_56bf9d=y9J2I5hN4sKjIiyZROsSAs...1g1isehgq.1g1isehgs.2.0.2; _gat_UA-35324627-3=1")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)

	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	src := string(bodyText)
	//fmt.Println(src)
	//"date", "open", "high", "low", "close", "p", "vol"
	dates := gjson.Get(src, "data.takerBuySellVol.#.0").Array()
	buys := gjson.Get(src, "data.takerBuySellVol.#.1").Array()
	sells := gjson.Get(src, "data.takerBuySellVol.#.2").Array()
	day := make([]string, len(dates))
	buy := make([]float64, len(buys))
	sell := make([]float64, len(sells))
	for i := 0; i < len(dates); i++ {
		a := dates[len(dates)-i-1].Str
		e, _ := strconv.ParseInt(a, 10, 64)
		day[i] = time.Unix(0, e*int64(time.Millisecond)).Format("2006-01-02 15:04:05")

		b := buys[len(buys)-i-1].Str
		f, _ := strconv.ParseFloat(b, 64)
		buy[i] = f

		g := sells[len(sells)-i-1].Str
		j, _ := strconv.ParseFloat(g, 64)
		sell[i] = j

	}

	return day, buy, sell
}

func GetRatio(symbol string, minute string) ([]string, []float64) {
	time.Sleep(time.Millisecond * 10)

	// 获取当前时间 或者使用 time.Date(year, month, ...)
	t := time.Now()
	timeStamp := t.Unix()
	client := &http.Client{

		Transport: &http.Transport{

			TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
		},
	}
	//https://www.okx.com/priapi/v5/rubik/public/stat/indicators?bar=4H&instId=XRP-USDT-SWAP&unit=2&limit=103&before=1725000000000&indicators=eliteLSAccountRatio&t=1725940771406
	req, err := http.NewRequest("GET", "https://www.okx.com/priapi/v5/rubik/public/stat/indicators?bar="+minute+"&instId="+symbol+"&unit=2&limit=1030&before=1725000000000&indicators=eliteLSAccountRatio&t="+strconv.Itoa(int(timeStamp)), nil)
	if err != nil {
		panic(err)

	}
	req.Header.Set("authority", "www.okx.com")
	req.Header.Set("timeout", "10000")
	req.Header.Set("x-cdn", "https://static.okx.com")
	req.Header.Set("devid", "a5ceb850-4efb-4a3f-baff-21da4fce8858")
	req.Header.Set("accept-language", "zh-CN")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36 SE 2.X MetaSr 1.0")
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-utc", "8")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("app-type", "web")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("referer", "https://www.okx.com/trade-swap/btc-usdt-swap")
	req.Header.Set("cookie", "locale=zh_CN; defaultLocale=zh_CN; _gcl_au=1.1.1520807996."+strconv.Itoa(int(timeStamp))+"; _ga=GA1.2.1752370991."+strconv.Itoa(int(timeStamp))+"; _gid=GA1.2.650161560."+strconv.Itoa(int(timeStamp))+"; amp_56bf9d=y9J2I5hN4sKjIiyZROsSAs...1g1isehgq.1g1isehgs.2.0.2; _gat_UA-35324627-3=1")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)

	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	src := string(bodyText)
	//fmt.Println(src)
	//"date", "open", "high", "low", "close", "p", "vol"
	dates := gjson.Get(src, "data.eliteLSAccountRatio.#.0").Array()
	ratios := gjson.Get(src, "data.eliteLSAccountRatio.#.3").Array()
	day := make([]string, len(dates))
	ratio := make([]float64, len(ratios))
	for i := 0; i < len(dates); i++ {
		a := dates[len(dates)-i-1].Str
		e, _ := strconv.ParseInt(a, 10, 64)
		day[i] = time.Unix(0, e*int64(time.Millisecond)).Format("2006-01-02 15:04:05")

		b := ratios[len(ratios)-i-1].Str
		f, _ := strconv.ParseFloat(b, 64)
		ratio[i] = f

	}
	return day, ratio
}

func GetKline(symbol string, minute string) (int, []float64, float64, float64, []float64, []float64, []float64, []float64, []string) {

	//fmt.Println(symboldemo, symbol)
	time.Sleep(time.Millisecond * 10)

	// 获取当前时间 或者使用 time.Date(year, month, ...)
	t := time.Now()
	timeStamp := t.Unix()
	client := &http.Client{

		Transport: &http.Transport{

			TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
		},
	}
	req, err := http.NewRequest("GET", "https://www.okx.com/priapi/v5/market/candles?instId="+symbol+"&bar="+minute+"&after=&limit=1400&t="+strconv.Itoa(int(timeStamp)), nil)
	if err != nil {
		panic(err)

	}
	req.Header.Set("authority", "www.okx.com")
	req.Header.Set("timeout", "10000")
	req.Header.Set("x-cdn", "https://static.okx.com")
	req.Header.Set("devid", "a5ceb850-4efb-4a3f-baff-21da4fce8858")
	req.Header.Set("accept-language", "zh-CN")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36 SE 2.X MetaSr 1.0")
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-utc", "8")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("app-type", "web")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("referer", "https://www.okx.com/trade-swap/btc-usdt-swap")
	req.Header.Set("cookie", "locale=zh_CN; defaultLocale=zh_CN; _gcl_au=1.1.1520807996."+strconv.Itoa(int(timeStamp))+"; _ga=GA1.2.1752370991."+strconv.Itoa(int(timeStamp))+"; _gid=GA1.2.650161560."+strconv.Itoa(int(timeStamp))+"; amp_56bf9d=y9J2I5hN4sKjIiyZROsSAs...1g1isehgq.1g1isehgs.2.0.2; _gat_UA-35324627-3=1")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)

	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	src := string(bodyText)
	//"date", "open", "high", "low", "close", "p", "vol"
	closes := gjson.Get(src, "data.#.4").Array()
	dates := gjson.Get(src, "data.#.0").Array()
	opens := gjson.Get(src, "data.#.1").Array()
	highs := gjson.Get(src, "data.#.2").Array()
	lows := gjson.Get(src, "data.#.3").Array()
	vols := gjson.Get(src, "data.#.6").Array()

	if len(closes) < 500 || closes[10].Float() < 0.0001 {
		return 0, []float64{}, 0, 0, []float64{}, []float64{}, []float64{}, []float64{}, []string{}
	}

	if len(closes) > 500 && closes[10].Float() > 0.0001 {

		day := make([]string, len(dates))
		c := make([]float64, len(closes))
		o := make([]float64, len(opens))
		h := make([]float64, len(highs))
		l := make([]float64, len(lows))
		v := make([]float64, len(vols))

		for i := 0; i < len(closes); i++ {

			a := dates[len(dates)-i-1].Str
			e, _ := strconv.ParseInt(a, 10, 64)
			day[i] = time.Unix(0, e*int64(time.Millisecond)).Format("2006-01-02 15:04:05")

			b := closes[len(closes)-i-1].Str
			f, _ := strconv.ParseFloat(b, 64)
			c[i] = f

			d := opens[len(opens)-i-1].Str
			g, _ := strconv.ParseFloat(d, 64)
			o[i] = g

			p := highs[len(highs)-i-1].Str
			rs, _ := strconv.ParseFloat(p, 64)
			h[i] = rs

			ll := lows[len(lows)-i-1].Str
			ll1, _ := strconv.ParseFloat(ll, 64)
			l[i] = ll1

			vv := vols[len(vols)-i-1].Str
			vol1, _ := strconv.ParseFloat(vv, 64)
			v[i] = vol1

		}

		x := len(c)

		return x, c, c[x-2] / o[x-2], c[x-1] / o[x-1], o, l, v, h, day
	}
	return 0, []float64{}, 0, 0, []float64{}, []float64{}, []float64{}, []float64{}, []string{}
}

//https://www.okx.com/priapi/v5/rubik/stat/taker-volume?instType=SPOT&period=1H&ccy=ETH&t=1742707202163

func GetTakerVolume(symbol string, minute string) bool {
	defer func() {
		if r := recover(); r != nil {
			// 处理异常
			fmt.Println("Exception caught:", r)
		}
	}()

	//fmt.Println(symboldemo, symbol)
	time.Sleep(time.Millisecond * 10)

	// 获取当前时间 或者使用 time.Date(year, month, ...)
	t := time.Now()
	timeStamp := t.Unix()
	client := &http.Client{

		Transport: &http.Transport{

			TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
		},
	}
	//https://www.okx.com/priapi/v5/public/liquidation-orders?instType=SWAP&instFamily=ETH-USDT&state=filled&limit=100&t=1739051712008
	req, err := http.NewRequest("GET", "https://www.okx.com/priapi/v5/rubik/stat/taker-volume?instType=SPOT&period="+minute+"&ccy="+symbol+"&t="+strconv.Itoa(int(timeStamp)), nil)
	if err != nil {
		panic(err)

	}
	req.Header.Set("authority", "www.okx.com")
	req.Header.Set("timeout", "10000")
	req.Header.Set("x-cdn", "https://static.okx.com")
	req.Header.Set("devid", "a5ceb850-4efb-4a3f-baff-21da4fce8858")
	req.Header.Set("accept-language", "zh-CN")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36 SE 2.X MetaSr 1.0")
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-utc", "8")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("app-type", "web")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("referer", "https://www.okx.com/trade-swap/btc-usdt-swap")
	req.Header.Set("cookie", "locale=zh_CN; defaultLocale=zh_CN; _gcl_au=1.1.1520807996."+strconv.Itoa(int(timeStamp))+"; _ga=GA1.2.1752370991."+strconv.Itoa(int(timeStamp))+"; _gid=GA1.2.650161560."+strconv.Itoa(int(timeStamp))+"; amp_56bf9d=y9J2I5hN4sKjIiyZROsSAs...1g1isehgq.1g1isehgs.2.0.2; _gat_UA-35324627-3=1")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)

	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	src := string(bodyText)

	dates := gjson.Get(src, "data.#.0").Array()
	buys := gjson.Get(src, "data.#.2").Array()
	sells := gjson.Get(src, "data.#.1").Array()
	day := make([]string, len(dates))
	buy := make([]float64, len(buys))
	sell := make([]float64, len(sells))
	for i := 0; i < len(dates); i++ {
		a := dates[len(dates)-i-1].Str
		e, _ := strconv.ParseInt(a, 10, 64)
		day[i] = time.Unix(0, e*int64(time.Millisecond)).Format("2006-01-02 15:04:05")

		b := buys[len(buys)-i-1].Str
		f, _ := strconv.ParseFloat(b, 64)
		buy[i] = f

		g := sells[len(sells)-i-1].Str
		j, _ := strconv.ParseFloat(g, 64)
		sell[i] = j

	}

	x := len(day)

	if x < 10 {
		return false
	} else {
		buys1 := (buy[x-1] + buy[x-2] + buy[x-3] + buy[x-4]) - (sell[x-1] + sell[x-2] + sell[x-3] + sell[x-4])
		buys2 := (buy[x-5] + buy[x-6] + buy[x-7] + buy[x-8]) - (sell[x-5] + sell[x-6] + sell[x-7] + sell[x-8])

		if buys1 > 0 && buys2 < 0 {
			return true
		} else {
			return false
		}
	}

}

func GetLiquidation() string {
	defer func() {
		if r := recover(); r != nil {
			// 处理异常
			fmt.Println("Exception caught:", r)
		}
	}()

	//fmt.Println(symboldemo, symbol)
	time.Sleep(time.Millisecond * 10)

	// 获取当前时间 或者使用 time.Date(year, month, ...)
	t := time.Now()
	timeStamp := t.Unix()
	client := &http.Client{

		Transport: &http.Transport{

			TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
		},
	}
	//https://www.okx.com/priapi/v5/public/liquidation-orders?instType=SWAP&instFamily=ETH-USDT&state=filled&limit=100&t=1739051712008
	req, err := http.NewRequest("GET", "https://www.okx.com/priapi/v5/public/liquidation-orders?instType=SWAP&instFamily=ETH-USDT&state=filled&limit=100&t="+strconv.Itoa(int(timeStamp)), nil)
	if err != nil {
		panic(err)

	}
	req.Header.Set("authority", "www.okx.com")
	req.Header.Set("timeout", "10000")
	req.Header.Set("x-cdn", "https://static.okx.com")
	req.Header.Set("devid", "a5ceb850-4efb-4a3f-baff-21da4fce8858")
	req.Header.Set("accept-language", "zh-CN")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.87 Safari/537.36 SE 2.X MetaSr 1.0")
	req.Header.Set("accept", "application/json")
	req.Header.Set("x-utc", "8")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("app-type", "web")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("referer", "https://www.okx.com/trade-swap/btc-usdt-swap")
	req.Header.Set("cookie", "locale=zh_CN; defaultLocale=zh_CN; _gcl_au=1.1.1520807996."+strconv.Itoa(int(timeStamp))+"; _ga=GA1.2.1752370991."+strconv.Itoa(int(timeStamp))+"; _gid=GA1.2.650161560."+strconv.Itoa(int(timeStamp))+"; amp_56bf9d=y9J2I5hN4sKjIiyZROsSAs...1g1isehgq.1g1isehgs.2.0.2; _gat_UA-35324627-3=1")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)

	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	src := string(bodyText)
	//fmt.Println(src)

	prices := gjson.Get(src, "data.#.details.#.price").Array()
	times := gjson.Get(src, "data.#.details.#.time").Array()
	szs := gjson.Get(src, "data.#.details.#.sz").Array()
	sides := gjson.Get(src, "data.#.details.#.side").Array()

	// 获取当前时间
	now := time.Now()

	// 将二维数组展平成一维数组
	var flatSzs []float64
	var flatTimes []string
	var flatSides []string
	var flatPrices []float64
	var isbuy string

	// 合并循环，同时遍历szs和times
	for i := 0; i < len(szs); i++ {
		innerSzs := szs[i].Array()
		innerTimes := times[i].Array()
		innerSides := sides[i].Array()
		innerPrices := prices[i].Array()

		// 遍历内部数组，并同时处理szs和times的元素
		for j := 0; j < len(innerSzs); j++ {
			szNum, _ := strconv.ParseFloat(innerSzs[j].String(), 64)

			flatSzs = append(flatSzs, szNum)

			side := innerSides[j].String()
			flatSides = append(flatSides, side)

			priceNum, _ := strconv.ParseFloat(innerPrices[j].String(), 64)

			flatPrices = append(flatPrices, priceNum)

			timeNum, _ := strconv.ParseFloat(innerTimes[j].String(), 64)

			// 假设这是你从JSON中获取的时间戳（毫秒级）
			timestampInMilliseconds := int64(timeNum)

			// 将毫秒级时间戳转换为秒级时间戳
			timestampInSeconds := timestampInMilliseconds / 1000

			// 将秒级时间戳转换为time.Time类型
			utcTime := time.Unix(timestampInSeconds, 0)

			// 将毫秒级时间戳转换为time.Time类型
			parsedTime := time.Unix(0, timestampInMilliseconds*int64(time.Millisecond))

			// 计算与当前时间的差距
			duration := now.Sub(parsedTime)

			// 格式化并打印时间
			formattedTime := utcTime.Format("2006-01-02 15:04:05")

			flatTimes = append(flatTimes, formattedTime)

			// 检查时间差距是否不超过5分钟，并且sz值是否大于100
			if duration.Minutes() <= 5 && szNum > 3000 && side == "sell" {

				log := "买入ETH--->>> ,time--->>>" + formattedTime +
					",sz--->>>" + fmt.Sprintf("%.5f", szNum) + ",price--->>>" + fmt.Sprintf("%.5f", priceNum)

				fmt.Println(log)
				GetWriter(log, "3m")

				isbuy = "buy"
			} else if duration.Minutes() <= 5 && szNum > 3000 && side == "buy" {

				log := "卖出ETH--->>> ,time--->>>" + formattedTime +
					",sz--->>>" + fmt.Sprintf("%.5f", szNum) + ",price--->>>" + fmt.Sprintf("%.5f", priceNum)

				fmt.Println(log)
				GetWriter(log, "3m")

				isbuy = "sell"
			}

		}
	}

	/**

	for i := 0; i < len(flatSzs); i++ {

		if flatSides[i]=="sell" && flatSzs[i]>100{

			fmt.Println("Flat Times:", flatTimes[i],"Flat Sides:", flatSides[i],"Flat Szs:", flatSzs[i],"Flat Prices:", flatPrices[i])
		}

	}
	*/

	// 打印结果
	//fmt.Println("Flat Szs:", flatSzs)
	//fmt.Println("Flat Times:", flatTimes)
	//fmt.Println("isbuy:", isbuy)
	return isbuy

}

func GetWriter(log string, minute string) {
	filePath := "..\\datas\\log\\buylog_" + minute + ".txt"

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)

	write.WriteString("\n")
	write.WriteString(log)

	//Flush将缓存的文件真正写入到文件中
	write.Flush()

}

func Buy(symbol string, minute string) {
	fmt.Println("-------------------------------买入--------------------------------->>>")

	cmd := exec.Command("python", "gorun.py", symbol, minute)
	res, _ := cmd.Output()
	fmt.Println(string(res))
}

func Getcashbal() {
	cmd := exec.Command("python", "cash.py")
	res, _ := cmd.Output()
	fmt.Println(string(res))
}

func Getcashhistory() {
	cmd := exec.Command("python", "cashhistory.py")
	res, _ := cmd.Output()
	fmt.Println(string(res))
}

func SellAll() {
	cmd := exec.Command("python", "sells.py")
	res, _ := cmd.Output()
	fmt.Println(string(res))
}

func GetuplRatio() {

	cmd := exec.Command("python", "getuplRatio.py")
	res, _ := cmd.Output()
	fmt.Println(string(res))

}

func Savecsv(minute string) {
	cmd := exec.Command("python", "savecsv.py", minute)
	res, _ := cmd.Output()
	fmt.Println(string(res))
}

func SendDingMsg(msg string) {
	//请求地址模板
	webHook := `https://oapi.dingtalk.com/robot/send?access_token=f8195c9e4ad6da4427d67e80dffed5d07ecaca1d1e79462fb5c0a9c6b12e90f2`
	content := `{"msgtype": "text",
        "text": {"content": "` + msg + `"}
    }`
	//创建一个请求
	req, err := http.NewRequest("POST", webHook, strings.NewReader(content))
	if err != nil {
		// handle error
	}

	client := &http.Client{}
	//设置请求头
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	//发送请求
	resp, err := client.Do(req)
	//关闭请求
	defer resp.Body.Close()

	if err != nil {
		// handle error
	}
}
