package Request
type Kline struct {
	TradePare TradePare
	Period int
	From int64
	To int64
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
