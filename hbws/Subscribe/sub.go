package Subscribe

type Kline struct {
	TradePare TradePare
	//x min
	Period int
}
type TradeDetail struct {
	TradePare TradePare
}

type MarketDetail struct {
	TradePare TradePare
}
type Depth struct {
	TradePare TradePare
	Step      int
}
type BBO struct {
	TradePare TradePare
}
type TradePare struct {
	Coin     string //btc
	BaseCoin string //usdt
}