package main

import (
	"github.com/cngobd/gohuobi/example"
	"github.com/cngobd/gohuobi/hbcontrol"
	"github.com/cngobd/gohuobi/hbws/Request"
	"log"
)

func main() {
	reqMarketDetail, err := hbcontrol.RunReqWork(Request.MarketDetail{TradePare:Request.TradePare{
		Coin:     "btc",
		BaseCoin: "usdt",
	}})
	if err != nil {
		log.Println(err)
	}
	go example.Listen(reqMarketDetail)
	select {
	}
}
