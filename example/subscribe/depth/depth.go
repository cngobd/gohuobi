package main

import (
	"github.com/cngobd/gohuobi/example"
	"github.com/cngobd/gohuobi/hbcontrol"
	"github.com/cngobd/gohuobi/hbws/Subscribe"
	"log"
)

func main() {
	subDepth, err := hbcontrol.RunSubWork(Subscribe.Depth{
		TradePare: Subscribe.TradePare{"btc","usdt"},
		Step:      0,
	})
	if err != nil {
		log.Println(err)
	}
	example.Listen(subDepth)
	select {

	}

}
