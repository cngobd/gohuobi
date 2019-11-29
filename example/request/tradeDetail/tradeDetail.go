package main

import (
	"github.com/cngobd/gohuobi/example"
	"github.com/cngobd/gohuobi/hbcontrol"
	"github.com/cngobd/gohuobi/hbws/Request"
	"log"
)

func main() {
	reqTradeDetail, err := hbcontrol.RunReqWork(Request.TradeDetail{TradePare: Request.TradePare{
			Coin:     "btc",
			BaseCoin: "usdt",
		}})
	if err != nil {
		log.Println(err)
	}
	go example.Listen(reqTradeDetail)
	select {

	}
}
