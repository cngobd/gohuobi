package main

import (
	"fmt"
	"github.com/cngobd/gohuobi/hbcontrol"
	"github.com/cngobd/hbws"
)
func printKline(resp hbws.UpdateKline) {
	fmt.Printf("tradePare:%v, amount:%v, open:%v, close:%v\n",
		resp.TradePare, resp.Amount, resp.Open, resp.Close)
}
func main() {
	//start a websocket subscribe work
	ch := hbcontrol.RunAWork(&hbws.WsReq{
		Detail:    "",
		Action:    hbws.Subscribe,
		TradePare: "btcusdt",
		Kind:      hbws.Kline,
		Parameter: "1min",
	})

	//start listen to the channel
	for {
		select {
		case x := <-ch :

			switch x.(type) {
			case []string :
				res := x.([]string)
				fmt.Printf("sub status  return:%v\n",res)
			case hbws.UpdateKline:
				printKline(x.(hbws.UpdateKline))

			case hbws.UpdateTradeDetail:

			case hbws.UpdateMarketDetail:

			case hbws.UpdateDepth:

			case hbws.UpdateBBO:
			}
		default:
		}
	}
}
