package hbws

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
	"github.com/cngobd/hbMaps"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
)

var (
	wsWorkerArr []*WsWorker
	url         = "wss://api.huobi.pro/ws"
	origin      = "http://127.0.0.1:8080"
)

type WsWorker struct {
	wsName    string
	client    *websocket.Conn
	ReqCh     chan interface{}
	RespCh    chan interface{}
	subTopics []*subTopic  //the topic subscribed currently
}

//this is a subscribe work instance
type subTopic struct {
	subbed string
	tradePare string
	kind      string
	parameter string
}

//start a websocket client and listen to the chan of ReqCh
func StartAWebSocketClient(WSName string) {
	defer fmt.Printf("%v closed\n",WSName)

	fmt.Printf("%v start...\n",WSName)
	newWorker := new(WsWorker)
	newWorker.ReqCh = make(chan interface{})
	newWorker.RespCh = make(chan interface{})
	newWorker.wsName = WSName
	var err error
	newWorker.client, err = websocket.Dial(url, "", origin)
	fmt.Println("after dial")
	defer newWorker.client.Close()
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Println("start store")

	hbMaps.WSMap.Store(newWorker.wsName, newWorker)

	go newWorker.messageHandle()
	for {
		select {
		case msg := <-newWorker.ReqCh:
			switch msg.(*WsReq).Action {
			//subscription a topic
			case "sub":
				//subscribe a new topic
				newWorker.subNewTopic(msg.(*WsReq))

			//unsubscription the topic
			case "unsub":
				newWorker.unsubTopic(msg.(*WsReq))

			//request the index
			case "req":
				newWorker.requestDataOnce(msg.(*WsReq))
			case "stop":
				return
			}
		}
	}
}

//listen the websocket and get message
func (w *WsWorker) messageHandle() {
	gzipFeature := []byte{31, 139}

	var box []byte
	for {
		mark1:
		msg := make([]byte, 5120)
		n, _ := w.client.Read(msg)
		if n == 0 {
			log.Print(errors.New("socket was closed"))
			break
		}

		if reflect.DeepEqual(msg[:2],gzipFeature) { //a new start
			if len(box) <= 0 { //is an empty box
				box = msg[:n]
				goto mark1 //append the data into box and try to read next data
			} else { //the box is full of data
				w.handleGzip(box) //start decompress
				box = msg[:n] //make a new box and append the start date into box
			}
		} else { //the data read is the rest data of last one
			box = append(box, msg[:n]...)
			goto mark1 //append the data into box and try to read next data
		}
	}
	log.Print(errors.New("message handle stopped"))
}
func (w *WsWorker)handleGzip(box []byte) {
	reader, err := gzip.NewReader(bytes.NewReader(box))
	if err != nil {
		log.Printf("gzip panic, box data : %v, length:%v\n", box[:20],len(box))
		panic(err)
	} else {
		result, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Print(err)
		} else {
			w.unmarshalMsg(result)
		}
	}
}
//handle the message form hb server, unmarshal []byte to string of struct,
// and use different
func (w *WsWorker) unmarshalMsg(index []byte) {
	//fmt.Printf("string index in unmarshal:%v\n",string(index))
	sel := new(selectMsg1)
	err := json.Unmarshal(index, sel)
	if err != nil {
		panic(err)
	} else {
		switch {
		case sel.Ping != 0:
			w.pong(index)
			return
		case len(sel.Id) > 0 && len(sel.Subbed) > 0:
			fmt.Println("select:", sel)
			w.subStatusHandle(index)
			return
		case len(sel.Ch) > 0 || len(sel.Rep) > 0:
			w.updateData(index) //handle the returned data which just subscribed
		}
	}
}

func (w *WsWorker) subStatusHandle(index []byte) {
	status := new(subStatus)
	fmt.Println("decompressed index:",string(index))
	err := json.Unmarshal(index, status)
	if err != nil {
		panic(err)
		return
	} else {
		w.RespCh <- []string{status.Id, status.Status, status.Subbed}
		if status.Status != "ok" {
			return
		}
		fmt.Println("status:", status)
		result := strings.Split(status.Subbed,".")
		st := new(subTopic)
		if len(result) >3 {
			st.subbed = status.Subbed
			st.tradePare,st.kind,st.parameter = result[1],result[2],result[3]
			w.subTopics = append(w.subTopics, st)
		} else {
			st.subbed = status.Subbed
			fmt.Println("result", result)
			st.tradePare,st.kind = result[1],result[2]
			w.subTopics = append(w.subTopics, st)
		}
	}
}

//send pong info to the hb server
func (w *WsWorker) pong(index []byte) {
	msg := make(map[string][]byte)
	msg["pong"] = index
	result, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = w.client.Write(result)
	if err != nil {
		log.Println(err)
	}
}

//subscribe a new topic, topic parameter: []string(tradePare, kind, parameter)
func (w *WsWorker) subNewTopic(topic *WsReq) {
	//[tradePare, kind, parameter]
	//kind : 1.kline 2.depth 3.trade 4.detail
	//parameter: 1.1min... 2.step0... 3.detail 4.nil
	fmt.Printf("start sub:%v\n",topic)

	if topic.Detail == "detail" { //sub the topic of the latest 24h market detail
		content := fmt.Sprintf("market.%v.%v", topic.TradePare, topic.Detail)
		js := map[string]string{
			"sub": content,
			"id":  fmt.Sprintf("%v.%v.%v", topic.TradePare, topic.Detail),
		}
		if marshal, err := json.Marshal(js); err == nil {
			w.client.Write(marshal)
		} else {
			log.Println(err)
			return
		}
	} else {
		var sub string
		var id string
		if topic.Parameter != "" {
			sub = fmt.Sprintf("market.%v.%v.%v", topic.TradePare, topic.Kind, topic.Parameter)
			id = fmt.Sprintf("%v.%v.%v", topic.TradePare, topic.Kind, topic.Parameter)
		} else {
			sub = fmt.Sprintf("market.%v.%v", topic.TradePare, topic.Kind)
			id = fmt.Sprintf("%v.%v", topic.TradePare, topic.Kind)
		}
		js := map[string]string{ //other subscriptions

			"sub": sub,
			"id":  id,
		}
		marshal, err := json.Marshal(js)
		if err != nil {
			log.Println(err)
			return
		} else {
			w.client.Write(marshal)
		}
	}
}

//
func (w *WsWorker) unsubTopic(topic *WsReq) {

}

//request data
func (w *WsWorker) requestDataOnce(topic *WsReq) {
	//[tradePare, kind, parameter]
	//kind : 1.kline 2.depth 3.trade 4.detail
	//parameter: 1.1min... 2.step0... 3.detail 4.nil
	if topic.Detail == "detail" { //request the latest 24h market detail
		msgToSend := fmt.Sprintf("market.%v.%v", topic.TradePare, topic.Detail)
		midToSend := map[string]string{
			"req": msgToSend,
			"id" : fmt.Sprintf("%v.%v", topic.TradePare, topic.Detail),
		}
		if jsToSend, err := json.Marshal(midToSend); err == nil {
			_, err := w.client.Write(jsToSend)
			if err != nil {
				panic(err)
			}
		} else {
			log.Println("write err:",err)
			panic(err)
			return
		}
	} else { //other requests
		msgToSend := fmt.Sprintf("market.%v.%v.%v",topic.TradePare, topic.Kind, topic.Parameter)
		midToSend := map[string]string{
			"req":msgToSend,
			"id": fmt.Sprintf("%v.%v.%v",topic.TradePare, topic.Kind, topic.Parameter),
		}
		if jsToSend, err := json.Marshal(midToSend); err != nil {
			log.Println(err)
			panic(err)
			return
		} else {
			_, err := w.client.Write(jsToSend)
			if err != nil {
				panic(err)
			}
		}

	}



}

func (w *WsWorker) updateData(index []byte) {
	//fmt.Printf("start updateData\n")
	upstr := new(updateDataSelect)
	err := json.Unmarshal(index, upstr)
	if err != nil {
		fmt.Printf("unmarshal 1 err\n")
		panic(err)
		return
	} else {
		var result []string
		if upstr.Ch != "" {
			result = strings.Split(upstr.Ch,".")
		} else {
			result = strings.Split(upstr.Rep, ".")
		}
		switch result[2] {
		case "kline":
			upKPre := new(UpdateKlinePre)
			err := json.Unmarshal(index, upKPre)
			if err != nil {
				panic(err)
				return
			}
			upKPre.Tick.TradePare, upKPre.Tick.Period = result[1], result[3]
			w.RespCh <- upKPre.Tick
		case "depth":
			upDepPre := new(UpdateDepthPre)
			err := json.Unmarshal(index,upDepPre)
			if err != nil {
				panic(err)
			}
			upDepPre.Tick.TradePare = result[1]
			w.RespCh <- upDepPre.Tick

		case "trade":
			tradeStr := new(UpdateTradeDetailPre)
			if err = json.Unmarshal(index, &tradeStr.Tick); err != nil {
				panic(err)
				return
			}
			tradeStr.Tick.TradePare = result[1]
			w.RespCh <- tradeStr.Tick

		case "detail":
			MarketDetailPre := new(UpdateMarketDetailPre)
			err := json.Unmarshal(index, MarketDetailPre)
			if err != nil {
				panic(err)
			}
			MarketDetailPre.Tick.TradePare = result[1]
			w.RespCh <- MarketDetailPre.Tick

		case "bbo":
			bboPre := new(UpdataBboPre)
			err := json.Unmarshal(index, bboPre)
			if err != nil {
				panic(err)
			}
			bboPre.Tick.TradePare = result[1]
			w.RespCh <- bboPre.Tick
		}
	}
}
