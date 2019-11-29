package example

import (
	"fmt"
	"github.com/cngobd/gohuobi/hbws"
	"log"
)

func Listen(ch chan interface{}) {
	//var m hbws.UpdateTradeDetail
	for {
		select {
		case msg := <- ch:
			switch msg.(type) {
			case []string :
				m := msg.([]string)

				fmt.Println("listen resp status:",m)
			case hbws.UpdateKline:
				fmt.Println("listen resp kline:",msg.(hbws.UpdateKline))
			case []hbws.UpdateKline:
				log.Println("len kline:", len(msg.([]hbws.UpdateKline)))
				for _, v := range msg.([]hbws.UpdateKline) {
					fmt.Println("listen resp []kline:",v)
				}
				//continue
				//fmt.Println("listen resp kline", msg.(hbws.UpdateKline))
			case hbws.UpdateTradeDetail:
				fmt.Println("listen resp detail:",msg.(hbws.UpdateTradeDetail).Data[0])
			case hbws.UpdateBBO:
				fmt.Println("listen resp bbo",msg.(hbws.UpdateBBO))
			case hbws.UpdateDepth:
				fmt.Println("listen resp depth", msg.(hbws.UpdateDepth))
			case hbws.UpdateMarketDetail:
				fmt.Println("listen resp market detail:",msg.(hbws.UpdateMarketDetail))
			default:
				log.Println("resp unknown type")
			}
		}
	}
}
