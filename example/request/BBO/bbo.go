package main

import (
	"github.com/cngobd/gohuobi/example"
	"github.com/cngobd/gohuobi/hbcontrol"
	"github.com/cngobd/gohuobi/hbws/Request"
	"log"
)
func main() {
	reqBBO, err := hbcontrol.RunReqWork(Request.BBO{
		TradePare:Request.TradePare{
			Coin:     "btc",
			BaseCoin: "usdt",
		},
	})
	if err != nil {
		log.Println(err)
	}
	go example.Listen(reqBBO)
	select {

	}
}
