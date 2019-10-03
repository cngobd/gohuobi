package controller

import (
	"fmt"
	"hbProject/hb/hbMaps"
	"hbProject/hb/websocketServer"
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
	runAWork("ws1", &websocketServer.WsReqPRM{
		Detail:    "",
		Action:    "sub",
		TradePare: "btcusdt",
		Kind:      "depth",
		Parameter: "step0",
	})
}

//this is a test function, send a request to the ReqCh
//and start listen to the RespCh, print the index from
//this channel

func runWs(wsName string) {
	fmt.Printf("listen %v start...\n",wsName)
	for {
		load, ok := hbMaps.WSMap.Load(wsName)
		if !ok {
			log.Print("not find in the map")
			time.Sleep(time.Second)
			continue
		} else {
			result := load.(*websocketServer.WsWorker)
			//[tradePare, kind, parameter]
			//kind : 1.kline 2.depth 3.trade 4.detail
			//parameter: 1.1min... 2.step0... 3.detail 4.nil
			fmt.Println("start send to chan")
			reqwork := new(websocketServer.WsReqPRM)
			reqwork.Action = "sub"
			reqwork.TradePare = "btcusdt"
			reqwork.Kind = "kline"
			reqwork.Parameter = "1min"
			result.ReqCh <- reqwork

			go listen(result.RespCh)
			break
		}
	}
}
func runSubBBO(wsName string) {

}

func runAWork(workName string, req *websocketServer.WsReqPRM) {
	go websocketServer.StartAWebSocketClient(workName)
	switch req.Action {
	case "sub":
		runASub(workName, req)
	case "req":
		runReqOnce(workName, req)
	}

}
func runASub(wsName string, param *websocketServer.WsReqPRM) {
	fmt.Printf("start a new work :%v, :%v\n",wsName, param)
	for {
		load, ok := hbMaps.WSMap.Load(wsName)
		if !ok {
			log.Print("check the network")
			time.Sleep(time.Second)
			continue
		} else {
			result := load.(*websocketServer.WsWorker)
			result.ReqCh <- param

			go listen(result.RespCh)
			break
		}
	}
}
func runReqOnce(wsName string, param *websocketServer.WsReqPRM) {
	fmt.Printf("start a new work :%v, :%v\n",wsName, param)
	for {
		load, ok := hbMaps.WSMap.Load(wsName)
		if !ok {
			log.Print("check the network...")
			time.Sleep(time.Second)
			continue
		} else {
			result := load.(*websocketServer.WsWorker)
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
			case websocketServer.UpdateKline:
				get := x.(websocketServer.UpdateKline)
				print(get)
			case websocketServer.UpdateTradeDetail:
				get := x.(websocketServer.UpdateTradeDetail)
				printTradeDetail(get)
			case websocketServer.UpdateMarketDetail:

			case websocketServer.UpdateDepth:
				printDepth(x.(websocketServer.UpdateDepth))
			case websocketServer.UpdateBBO:
				printBBO(x.(websocketServer.UpdateBBO))
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

func printDepth(get websocketServer.UpdateDepth) {
	fmt.Printf("sub depth resp: \n bid:%v \n ask:%v\n", get.Bids[:3],get.Asks[:3])
	fmt.Println("-----------")
}
func printBBO(get websocketServer.UpdateBBO) {
	fmt.Println("sub bbo resp:", get)
}
//print the info from the RespCH
func print(get websocketServer.UpdateKline) {
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
func printTradeDetail(get websocketServer.UpdateTradeDetail) {
	fmt.Println("data length", len(get.Data))
	fmt.Printf("data: %v, %v, %v \n",get.Data[len(get.Data)-1].Amount,get.Data[len(get.Data)-1].Price,get.Data[len(get.Data)-1].Direction)
	fmt.Printf("data: %v, %v, %v \n",get.Data[len(get.Data)-2].Amount,get.Data[len(get.Data)-2].Price,get.Data[len(get.Data)-2].Direction)
	fmt.Printf("data: %v, %v, %v \n",get.Data[len(get.Data)-3].Amount,get.Data[len(get.Data)-3].Price,get.Data[len(get.Data)-3].Direction)
	fmt.Println("ts:",time.Now().Second())
	fmt.Println("------")
}