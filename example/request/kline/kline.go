package main

import (
	"github.com/cngobd/gohuobi/example"
	"github.com/cngobd/gohuobi/hbcontrol"
	"github.com/cngobd/gohuobi/hbws/Request"
	"log"
)

func main() {
	reqKline, err := hbcontrol.RunReqWork(Request.Kline{
		TradePare: Request.TradePare{"btc", "usdt"},
		Period:    1,
		//From:      0,
		//To:        0,
	})
	if err != nil {
		log.Println(err)
	}
	go example.Listen(reqKline)
	select {

	}
}
