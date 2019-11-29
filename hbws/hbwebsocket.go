package hbws

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
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
	Client    *websocket.Conn
	ReqCh     chan ReqMessage
	RespCh    chan interface{}
	StrChan chan interface{}
	subTopics []*subTopic  //the topic subscribed currently
}
type ReqMessage struct {
	Command string
	Topic string
}
//this is a subscribe work instance
type subTopic struct {
	subbed string
	tradePare string
	kind      string
	parameter string
}

//start a websocket client and listen to the chan of ReqCh
func StartAWebSocketClient(newWorker *WsWorker) {
	defer fmt.Printf("closed")
	//log.Println("ws client start")
	var err error
	newWorker.Client, err = websocket.Dial(url, "", origin)
	defer newWorker.Client.Close()
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Println("dial success, start working")

	go newWorker.messageHandle()
	for {
		select {
		case msg := <- newWorker.ReqCh:
			switch msg.Command {
			//subscription a topic
			case "sub":
				//subscribe a new topic
				newWorker.subNewTopic(msg.Topic)
			//unsubscription the topic
			case "unsub":
				newWorker.unsubTopic(msg.Topic)
			//request the index
			case "req":
				newWorker.requestDataOnce(msg.Topic)
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
		n, _ := w.Client.Read(msg)
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
			w.subStatusHandle(index)
			return
		case len(sel.Ch) > 0 || len(sel.Rep) > 0:
			w.updateData(index) //handle the returned data which just subscribed
		}
	}
}

func (w *WsWorker) subStatusHandle(index []byte) {
	status := new(subStatus)
	err := json.Unmarshal(index, status)
	if err != nil {
		panic(err)
		return
	} else {
		w.RespCh <- []string{status.Id, status.Status, status.Subbed}
		/*if status.Status != "ok" {
			return
		}
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
		}*/
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
	_, err = w.Client.Write(result)
	if err != nil {
		log.Println(err)
	}
}

//subscribe a new topic, topic parameter: []string(tradePare, kind, parameter)
func (w *WsWorker) subNewTopic(topic string) {
	//fmt.Printf("start sub:%v\n",topic)
	js := map[string]string{ //other subscriptions
		"sub": topic,
		"id":  topic,
	}
	marshal, err := json.Marshal(js)
	if err != nil {
		log.Println(err)
		return
	} else {
		w.Client.Write(marshal)
	}
	//[tradePare, kind, parameter]
	//kind : 1.kline 2.depth 3.trade 4.detail
	//parameter: 1.1min... 2.step0... 3.detail 4.nil


	/*if topic.Detail == "detail" { //sub the topic of the latest 24h market detail
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

	}*/
}

//
func (w *WsWorker) unsubTopic(topic string) {

}

//request data 10times per second automatically
func (w *WsWorker) requestDataOnce(topic string) {
	//[tradePare, kind, parameter]
	//kind : 1.kline 2.depth 3.trade 4.detail
	//parameter: 1.1min... 2.step0... 3.detail 4.nil
	//log.Println("topic", topic)

	midToSend := map[string]string{
		"req":topic,
		"id": topic,
		//"from": fmt.Sprintf("%v",time.Now().Unix() -100),
		//"to": fmt.Sprintf("%v",time.Now().Unix()-10),
	}
	//log.Println("mid:", midToSend)
	if jsToSend, err := json.Marshal(midToSend); err != nil {
		log.Println(err)
		panic(err)
		return
	} else {
		//1572920982
		//1501174800
		//log.Println("js to write:", string(jsToSend))
		_, err := w.Client.Write(jsToSend)
		if err != nil {
			panic(err)
		}
	}
	/*if topic.Detail == "detail" { //request the latest 24h market detail
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
		var msgToSend string
		if topic.Parameter == "" {
			msgToSend = fmt.Sprintf("market.%v.%v",topic.TradePare, topic.Kind)
		} else {
			msgToSend = fmt.Sprintf("market.%v.%v.%v",topic.TradePare, topic.Kind, topic.Parameter)
		}

		log.Println("msgToSend:", msgToSend)



	}*/



}

func (w *WsWorker) updateData(index []byte) {
	//fmt.Printf("start updateData\n")
	//fmt.Println("update date:", string(index))
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
			//upKPre.Tick.TradePare, upKPre.Tick.Period = result[1], result[3]
			if len(upKPre.Data) <= 0  {
				//log.Println("got tick")
				w.RespCh <- upKPre.Tick
			} else {
				//log.Println("got data")
				w.RespCh <- upKPre.Data
			}

		case "depth":
			upDepPre := new(UpdateDepthPre)
			err := json.Unmarshal(index,upDepPre)
			if err != nil {
				panic(err)
			}
			upDepPre.Data.TradePare = result[1]
			w.RespCh <- upDepPre.Data

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
