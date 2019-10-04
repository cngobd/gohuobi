package hbhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	TradePairArr   []TradePair
	TradePairCount int
)
var (
	Currency       []string
	currencysCount int
)

var tpr All
var cur All2
var tp All3

type All struct {
	Status string      `json:"status"`
	Data   []TradePair `json:"data"`
}
type All2 struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}
type All3 struct {
	Status string `json:"status"`
	Data   int    `json:"data"`
}
type Kl struct {
	Status string `json:"status"`
	Ch     string `json:"ch"`
	Ts     int    `json:"ts"`
	Data   []Kline
}
type TickerJS struct {
	Status string `json:"status"`
	Ch     string `json:"ch"`
	Ts     int    `json:"ts"`
	Tick   ticker `json:"tick"`
}
type Tickers struct {
	Status string    `json:"status"`
	Ts     int       `json:"ts"`
	Data   []tickers `json:"data"`
}
type Depth struct {
	Status string    `json:"status"`
	Ch     string    `json:"ch"`
	Ts     int       `json:"ts"`
	Tick   depthTick `json:"tick"`
}
type CurTrade struct {
	Status string   `json:"status"`
	Ch     string   `json:"ch"`
	Ts     int      `json:"ts"`
	Tick   curTrade `json:"tick"`
}
type HistoryTrade struct {
	Status string             `json:"status"`
	Ch     string             `json:"ch"`
	Ts     int                `json:"ts"`
	Data   []dataHistoryTrade `json:"data"`
}
type Detail24 struct {
	Status string   `json:"status"`
	Ch     string   `json:"ch"`
	Ts     int      `json:"ts"`
	Tick   detail24 `json:"tick"`
}
type ticker struct {
	ID     int       `json:"id"`
	Amount float64   `json:"amount"`
	Open   float64   `json:"open"`
	Close  float64   `json:"close"`
	Low    float64   `json:"low"`
	High   float64   `json:"high"`
	Vol    float64   `json:"vol"`
	Bid    []float64 `json:"bid"`
	Ask    []float64 `json:"ask"`
}

type tickers struct {
	Amount float64 `json:"amount"`
	Count  int     `json:"count"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	Low    float64 `json:"low"`
	High   float64 `json:"high"`
	Vol    float64 `json:"vol"`
	Symble string  `json:"symble"`
}
type depthTick struct {
	Version int         `json:"version"`
	Ts      int         `json:"ts"`
	Bids    [][]float64 `json:"bids"`
	Asks    [][]float64 `json:"asks"`
}
type curTrade struct {
	Id   int            `json:"id"`
	Ts   int            `json:"ts"`
	Data []dataCurTrade `json:"data"`
}
type dataCurTrade struct {
	Amount float64 `json:"amount"`
	Ts     int     `json:"ts"`
	//Id        int     `json:"id"`
	Price     float64 `json:"price"`
	Direction string  `json:"direction"`
}

type dataHistoryTrade struct {
	Id   int           `json:"id"`
	Ts   int           `json:"ts"`
	Data []historyData `json:"data"`
}
type historyData struct {
	//Id int `json:"id"`
	Amount    float64 `json:"amount"`
	Price     float64 `json:"price"`
	Ts        int     `json:"ts"`
	Direction string  `json:"direction"`
}

type detail24 struct {
	//Id int `json:"id"`
	Amount  float64 `json:"amount"`
	Count   int     `json:"count"`
	Open    float64     `json:"open"`
	Close   float64 `json:"close"`
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Vol     float64 `json:"vol"`
	Version int     `json:"version"`
}

//{"base-currency":"snt","quote-currency":"usdt","price-precision":6,
// "amount-precision":4,"symbol-partition":"innovation","symbol":"sntusdt",
// "state":"online","value-precision":8,"min-order-amt":0.1,
// "max-order-amt":1000000,"min-order-value":0.1},
type TradePair struct {
	BaseCurrency    string  `json:"base-currency"`    //交易对基础币种
	QuoteCurrency   string  `json:"quote-currency"`   //交易对报价币种
	AmountPrecision int8    `json:"amount-precision"` //交易对基础币种精度，小数点后的位数
	PricePrecision  int8    `json:"price-precision"`  //交易对报价币种精度，小数点后的位数
	SymbolPartition string  `json:"symbol-partition"` //币种交易区 main/innovation
	State           string  `json:"state"`            //交易对状态 online/offline/suspend
	ValuePrecision  int8    `json:"value-precision"`  //交易额精度位数
	MinOrderAmt     float32 `json:"min-order-amt"`    //交易对最小下单数量
	MaxOrderAmt     float32 `json:"max-order-amt"`    //交易对最大下单数量
	MinOrderValue   float32 `json:"min-order-value"`  //交易对最小下单金额
	LeverageRatio   float32 `json:"leverage-ratio"`   //杠杆倍率
}

type Kline struct {
	ID     int     `json:"id"`     //kLine时间戳
	Amount float64 `json:"amount"` //基础币种计量的交易量
	Vol    float64 `json:"vol"`    //报价币种计量的交易量
	Count  int     `json:"count"`  //交易次数
	Open   float64 `json:"open"`   //开盘价
	Close  float64 `json:"close"`  //收盘价
	Low    float64 `json:"low"`    //最低价
	High   float64 `json:"high"`   //最高价
}

type ask struct {
	price       float64
	quoteVolume float64
}
type bid struct {
	price       float64
	quotoVolume float64
}

//
func GetAllTradePair() {
	resp, err := http.Get("https://api.huobi.pro/v1/common/symbols")
	readAll, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}
	err = json.Unmarshal(readAll, &tpr)
	if err != nil {
		log.Println(err)
	}
	if tpr.Status != "ok" {
		log.Println(tpr.Status)
		return
	}
	TradePairArr = tpr.Data
	TradePairCount = len(TradePairArr)
	log.Println("finished")
}
func GetAllSupportedCurrency() {
	resp, err := http.Get("https://api.huobi.pro/v1/common/currencys")
	if err != nil {
		log.Println(err)
		return
	}
	readAll, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(readAll, &cur)
	if err != nil {
		log.Println(err)
		return
	}
	if cur.Status != "ok" {
		log.Println(err)
		return
	}
	Currency = cur.Data
	currencysCount = len(Currency)
	log.Println("finished")
}
func GetTimeStamp() int {
	resp, err := http.Get("https://api.huobi.pro/v1/common/timestamp")
	if err != nil {
		log.Println(err)
		return 0
	}
	readAll, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return 0
	}
	err = json.Unmarshal(readAll, &tp)
	if err != nil {
		log.Println(err)
		return 0
	} else {
		return tp.Data
	}
}

func GetKline(symble string, period string, size int) []Kline {
	KL := new(Kl)
	resp, err := http.Get(fmt.Sprintf(
		"https://api.huobi.pro/market/history/kline?period=%v&size=%v&symbol=%v",
		period, size, symble))
	if err != nil {
		log.Println(err)
		return nil
	}
	readAll, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = json.Unmarshal(readAll, KL)
	if err != nil {
		log.Println(err)
		return nil
	}
	fmt.Println("status", KL.Status)
	return KL.Data
}

func GetTicker(symbol string) *ticker {
	ti := new(TickerJS)
	resp, err := http.Get(fmt.Sprintf("https://api.huobi.pro/market/detail/merged?symbol=%v", symbol))
	if err != nil {
		log.Println(err)
		return nil
	}
	readAll, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = json.Unmarshal(readAll, ti)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &ti.Tick
}

func GetTickers() []tickers {
	ti := new(Tickers)
	resp, err := http.Get("https://api.huobi.pro/market/tickers")
	if err != nil {
		log.Println(err)
		return nil
	}
	readAll, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = json.Unmarshal(readAll, ti)
	if err != nil {
		log.Println(err)
		return nil
	}
	return ti.Data
}

func GetDepth(symbol string, depth int, tp string) (depthTick, error) {
	dpt := new(Depth)
	resp, err := http.Get(fmt.Sprintf(
		"https://api.huobi.pro/market/depth?symbol=%v&depth=%v&type=%v",
		symbol, depth, tp))
	if err != nil {
		log.Print(err)
		return depthTick{}, err
	}
	readAll, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return depthTick{}, err
	}
	err = json.Unmarshal(readAll, dpt)
	if err != nil {
		log.Print(err)
		return depthTick{}, err
	}
	return dpt.Tick, nil
}

func GetCurTrade(symbol string) (curTrade, error) {
	ct := new(CurTrade)
	resp, err := http.Get(fmt.Sprintf("https://api.huobi.pro/market/trade?symbol=%v", symbol))
	if err != nil {
		log.Print(err)
		return curTrade{}, nil
	}
	readAll, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return curTrade{}, err
	}
	err = json.Unmarshal(readAll, ct)
	if err != nil {
		log.Print(err)
		return curTrade{}, err
	}
	return ct.Tick, nil
}

func GetHistoryTrade(symbol string, size int) ([]dataHistoryTrade, error) {
	hisTrade := new(HistoryTrade)
	resp, err := http.Get(fmt.Sprintf(
		"https://api.huobi.pro/market/history/trade?symbol=%v&size=%v", symbol, size))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	readAll, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = json.Unmarshal(readAll, hisTrade)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return hisTrade.Data, nil
}

func Get24Detail(symbol string) (detail24,error){
	detail := new(Detail24)
	resp, err := http.Get(fmt.Sprintf("https://api.huobi.pro/market/detail?symbol=%v", symbol))
	if err != nil {
		log.Println(err)
		return detail24{}, err
	}
	readAll, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return detail24{}, err
	}
	err = json.Unmarshal(readAll, detail)
	if err != nil {
		log.Println(err)
		return detail24{}, err
	}
	return detail.Tick, nil
}
