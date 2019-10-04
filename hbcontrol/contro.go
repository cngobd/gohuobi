package hbcontrol

import (
	"fmt"
	"github.com/cngobd/hbMaps"
	"github.com/cngobd/hbws"
	"log"
	"time"
)

func init() {
	/*go websocketServer.StartAWebSocketClient("ws1")
	time.Sleep(time.Second)
	go runWs("ws1")*/

	/*go websocketServer.StartAWebSocketClient("ws2")
	time.Sleep(time.Second)
	go runReqOnce("ws2")*/

	/*go websocketServer.StartAWebSocketClient("ws3")
	time.Sleep(time.Second)
	go runSubBBO("ws3")*/
	RunAWork(&hbws.WsReq{
		Detail:    "",
		Action:    hbws.Subscribe,
		TradePare: "btcusdt",
		Kind:      hbws.Kline,
		Parameter: "step0",
	})
}

//this is a test function, send a request to the ReqCh
//and start listen to the RespCh, print the index from
//this channel


func RunAWork(req *hbws.WsReq) chan interface{}{
	workName := fmt.Sprintf("%v%v", req.Action, req.TradePare)
	go hbws.StartAWebSocketClient(workName)
	switch req.Action {
	case "sub":
		return runASub(workName, req)
	case "req":
		return runReqOnce(workName, req)
	}
	return nil
}
func runASub(wsName string, param *hbws.WsReq) chan interface{}{
	fmt.Printf("start a new work :%v, :%v\n",wsName, param)
	for {
		load, ok := hbMaps.WSMap.Load(wsName)
		if !ok {
			log.Print("check the network")
			time.Sleep(time.Second)
			continue
		} else {
			result := load.(*hbws.WsWorker)
			result.ReqCh <- param

			//go listen(result.RespCh)
			return result.RespCh
		}
	}
}
func runReqOnce(wsName string, param *hbws.WsReq) chan interface{}{
	fmt.Printf("start a new work :%v, :%v\n",wsName, param)
	for {
		load, ok := hbMaps.WSMap.Load(wsName)
		if !ok {
			log.Print("check the network...")
			time.Sleep(time.Second)
			continue
		} else {
			result := load.(*hbws.WsWorker)
			go listen(result.RespCh)
		mark1:
			//[tradePare, kind, parameter]
			//kind : 1.kline 2.depth 3.trade 4.detail
			//parameter: 1.1min... 2.step0... 3.detail 4.nil
			//fmt.Println("start send to chan")

			result.ReqCh <- param
			time.Sleep(time.Second/10)
			goto mark1
			//result.ReqCh <- []string{"sub","btcusdt","kline","1min"}
		}
	}
}
var count int
var lastTime = time.Now()
//listen to the RespCH
func listen(ch chan interface{}){
	for {
		select {
		case x := <-ch :
			count ++
			switch x.(type) {
			case []string :
				res := x.([]string)
				fmt.Printf("server status index return:%v\n",res)
			case hbws.UpdateKline:
				get := x.(hbws.UpdateKline)
				print(get)
			case hbws.UpdateTradeDetail:
				get := x.(hbws.UpdateTradeDetail)
				printTradeDetail(get)
			case hbws.UpdateMarketDetail:

			case hbws.UpdateDepth:
				printDepth(x.(hbws.UpdateDepth))
			case hbws.UpdateBBO:
				printBBO(x.(hbws.UpdateBBO))
			}

		default:
			if time.Since(lastTime) > 1*time.Second {
				log.Println("count:",count)
				lastTime = time.Now()
				count = 0
			}
		}
	}
}

func printDepth(get hbws.UpdateDepth) {
	fmt.Printf("sub depth resp: \n bid:%v \n ask:%v\n", get.Bids[:3],get.Asks[:3])
	fmt.Println("-----------")
}
func printBBO(get hbws.UpdateBBO) {
	fmt.Println("sub bbo resp:", get)
}
//print the info from the RespCH
func print(get hbws.UpdateKline) {
	fmt.Printf("tradepare:%v\n",get.TradePare)
	fmt.Printf("id:%v\n", get.Id)
	fmt.Printf("amount:%v\n",get.Amount)
	fmt.Printf("count:%v\n",get.Count)
	fmt.Printf("open:%v\n",get.Open)
	fmt.Printf("close:%v\n",get.Close)
	fmt.Printf("high:%v\n",get.High)
	fmt.Printf("low:%v\n",get.Low)
	fmt.Printf("volum:%v\n",get.Vol)
	fmt.Printf("%v\n","---------")
}
func printTradeDetail(get hbws.UpdateTradeDetail) {
	fmt.Println("data length", len(get.Data))
	fmt.Printf("data: %v, %v, %v \n",get.Data[len(get.Data)-1].Amount,get.Data[len(get.Data)-1].Price,get.Data[len(get.Data)-1].Direction)
	fmt.Printf("data: %v, %v, %v \n",get.Data[len(get.Data)-2].Amount,get.Data[len(get.Data)-2].Price,get.Data[len(get.Data)-2].Direction)
	fmt.Printf("data: %v, %v, %v \n",get.Data[len(get.Data)-3].Amount,get.Data[len(get.Data)-3].Price,get.Data[len(get.Data)-3].Direction)
	fmt.Println("ts:",time.Now().Second())
	fmt.Println("------")
}