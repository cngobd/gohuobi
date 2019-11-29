package main

import (
	"github.com/cngobd/gohuobi/example"
	"github.com/cngobd/gohuobi/hbcontrol"
	"github.com/cngobd/gohuobi/hbws/Subscribe"
	"log"
)

func main() {
	subTradeDetail, err := hbcontrol.RunSubWork(Subscribe.TradeDetail{TradePare:Subscribe.TradePare{
		Coin:     "btc",
		BaseCoin: "usdt",
	}})
	if err != nil {
		log.Println(err)
	}
	go example.Listen(subTradeDetail)
	select {

	}
}
