package main

import (
	"github.com/cngobd/gohuobi/example"
	"github.com/cngobd/gohuobi/hbcontrol"
	"github.com/cngobd/gohuobi/hbws/Subscribe"
	"log"
)

func main() {
	respch, err := hbcontrol.RunSubWork(Subscribe.Kline{
		TradePare: Subscribe.TradePare{"btc", "usdt"},
		Period:      1,
	})
	if err != nil {
		log.Println(err)
	}
	go example.Listen(respch)
	select {

	}
}
