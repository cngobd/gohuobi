package hbws

var (
	Subscribe = "sub"
	Request = "req"
)
const (
	Kline = "kline"
	Depth = "depth"
	Trade = "trade"
	Detail = "detail"
)
//
type WsReq struct {
	Detail string //detail of the market
	Action string //sub, unsub, req, stop
	TradePare string //btcusdt, ethusdt...
	Kind      string //kline, depth, trade, detail
	Parameter string //1.1min... 2.step0... 3.detail 4.nil
}

//this struct helps to divide the returned data into sub or req
type selectMsg1 struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Ch     string `json:"ch"`
	Rep    string `json:"rep"`
	Ping   int    `json:"ping"`
	Subbed string `json:"subbed"`
}

type subStatus struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Subbed string `json:"subbed"`
}
type reqStatus struct {
	Status string `json:"status"`
	Rep string `json:"rep"`
	Tick []byte `json:"tick"`
}
type unSubStatus struct {
	Id       string `json:"id"`
	Status   string `json:"status"`
	Unsubbed string `json:"unsubbed"`
}
type updateDataSelect struct {
	Ch string `json:"ch"`
	Rep string `json:"rep"`
	Tick interface{} `json:"tick"`
}

type UpdateKlinePre struct {
	Tick UpdateKline `json:"tick"`
}
type UpdateKline struct {
	TradePare string
	Period    string  //trade period  //1min, 5min, 15min, 30min, 60min, 1day, 1mon, 1week, 1year
	Id        int     `json:"id"`     //unix time
	Amount    float64 `json:"amount"` //the trade amount in this period
	Count     int     `json:"count"`  //the count of trade in this period
	Open      float64 `json:"open"`   //the open price
	Close     float64 `json:"close"`  //the close price
	Low       float64 `json:"low"`    //the lowest price in this period
	High      float64 `json:"high"`   //the highest price in this period
	Vol       float64 `json:"vol"`    //the trade sum
}
type UpdateDepthPre struct {
	Tick UpdateDepth `json:"tick"`
}
type UpdateDepth struct {
	TradePare string
	Bids      [][]float64 `json:"bids"`
	Asks       [][]float64 `json:"asks"`
}
type UpdateTradeDetailPre struct {
	Tick UpdateTradeDetail `json:"tick"`
}
type UpdateTradeDetail struct { //This topic sends the latest completed trade
	TradePare string
	Data      []TradeDetail
}
type TradeDetail struct {
	Amount    float64 `json:"amount"`
	Ts        int     `json:"ts"`
	Price     float64 `json:"price"`
	Direction string  `json:"direction"`
}

//latest 24h summary of market
type UpdateMarketDetailPre struct {
	Tick UpdateMarketDetail `json:"tick"`
}
type UpdateMarketDetail struct {
	TradePare string
	Ts     int     `json:"ts"`
	Amount float64 `json:"amount"`
	Count  int     `json:"count"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Vol    float64 `json:"vol"`
}

type UpdataBboPre struct {
	Tick UpdateBBO `json:"tick"`
}
type UpdateBBO struct {
	TradePare string
	QuoteTime int64 `json:"quoteTime"`
	Bid float32 `json:"bid"`
	BidSize float32 `json:"bidSize"`
	Ask float32 `json:"ask"`
	AskSize float32 `json:"askSize"`
}
type RequestOnce struct {

}



type RespTradeDetail struct {

}
/*
//sub status return
{
"id": "id1",
"status": "ok",
"subbed": "market.btccny.kline.1min",
"ts": 1489474081631
}

//update market kline
{
  "ch": "market.btccny.kline.1min",
  "ts": 1489474082831,
  "tick": {
    "id": 1489464480,
    "amount": 0.0,
    "count": 0,
    "open": 7962.62,
    "close": 7962.62,
    "low": 7962.62,
    "high": 7962.62,
    "vol": 0.0
  }
}

//update market depth
{
  "ch": "market.btcusdt.depth.step0",
  "ts": 1489474082831,
  "tick": {
    "bids": [
    [9999.3900,0.0098], // [price, amount]
    [9992.5947,0.0560]
    // more Market Depth data here
    ],
    "asks": [
    [10010.9800,0.0099],
    [10011.3900,2.0000]
    //more data here
    ]
  }
}

//update market trade detail
{
  "ch": "market.btcusdt.trade.detail",
  "ts": 1489474082831,
  "tick": {
        "id": 14650745135,
        "ts": 1533265950234,
        "data": [
            {
                "amount": 0.0099,
                "ts": 1533265950234,
                "id": 146507451359183894799,
                "price": 401.74,
                "direction": "buy"
            }
            // more Trade Detail data here
        ]
  }
}

//update market detail
  "tick": {
    "amount": 12224.2922,
    "open":   9790.52,
    "close":  10195.00,
    "high":   10300.00,
    "ts":     1494496390000,
    "id":     1494496390,
    "count":  15195,
    "low":    9657.00,
    "vol":    121906001.754751
  }


//request data once
{
  "status": "ok",
  "rep": "market.btccny.kline.1min",
  "tick": [
    {
      "amount": 1.6206,
      "count":  3,
      "id":     1494465840,
      "open":   9887.00,
      "close":  9885.00,
      "low":    9885.00,
      "high":   9887.00,
      "vol":    16021.632026
    },
    {
      "amount": 2.2124,
      "count":  6,
      "id":     1494465900,
      "open":   9885.00,
      "close":  9880.00,
      "low":    9880.00,
      "high":   9885.00,
      "vol":    21859.023500
    }
  ]
}


*/
