package main

import (
	"github.com/cngobd/gohuobi/example"
	"github.com/cngobd/gohuobi/hbcontrol"
	"github.com/cngobd/gohuobi/hbws/Subscribe"
	"log"
)

func main() {
	subMarketDetail, err := hbcontrol.RunSubWork(Subscribe.MarketDetail{TradePare:Subscribe.TradePare{
		Coin:     "btc",
		BaseCoin: "usdt",
	}})
	if err != nil {
		log.Println(err)
	}
	go example.Listen(subMarketDetail)
	select {
	}

}
