package main

import (
	"github.com/cngobd/gohuobi/example"
	"github.com/cngobd/gohuobi/hbcontrol"
	"github.com/cngobd/gohuobi/hbws/Subscribe"
	"log"
)
func main() {
	subBBO, err := hbcontrol.RunSubWork(Subscribe.BBO{TradePare:Subscribe.TradePare{
		Coin:     "btc",
		BaseCoin: "usdt",
	}})
	if err != nil {
		log.Println(err)
	}
	go example.Listen(subBBO)
	select {

	}
}
