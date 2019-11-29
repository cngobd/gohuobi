package main

import (
	"github.com/cngobd/gohuobi/example"
	"github.com/cngobd/gohuobi/hbcontrol"
	"github.com/cngobd/gohuobi/hbws/Request"
	"log"
)

func main() {
	reqDepth, err := hbcontrol.RunReqWork(Request.Depth{
		TradePare: Request.TradePare{"btc","usdt"},
		Step:      0,
	})
	if err != nil {
		log.Println(err)
	}
	go example.Listen(reqDepth)
	select {

	}
}
