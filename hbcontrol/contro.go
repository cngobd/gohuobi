package hbcontrol

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cngobd/gohuobi/hbws"
	"github.com/cngobd/gohuobi/hbws/Subscribe"
	"github.com/cngobd/gohuobi/hbws/Request"

)


//this is a test function, send a request to the ReqCh
//and start listen to the RespCh, print the index from
//this channel

func RunSubWork(param interface{}) (chan interface{}, error){

	var workName string

	fmt.Printf("start a new sub work :%v\n",workName)
	newWorker := new(hbws.WsWorker)
	newWorker.ReqCh = make(chan hbws.ReqMessage,1000)
	newWorker.RespCh = make(chan interface{},1000)
	go hbws.StartAWebSocketClient(newWorker)
	command := "sub"
	topic := ""
	switch param.(type) {
	case Subscribe.Kline:
		tp := param.(Subscribe.Kline)
		topic = fmt.Sprintf("market.%v.kline.%vmin",getTradePair(tp.TradePare), tp.Period)


	case Subscribe.TradeDetail:
		tp := param.(Subscribe.TradeDetail)
		topic = fmt.Sprintf("market.%v.trade.detail",getTradePair(tp.TradePare))
	case Subscribe.Depth:
		tp := param.(Subscribe.Depth)
		topic = fmt.Sprintf("market.%v.depth.step%v",getTradePair(tp.TradePare),tp.Step)
	case Subscribe.BBO:
		tp := param.(Subscribe.BBO)
		topic = fmt.Sprintf("market.%v.bbo",getTradePair(tp.TradePare))
	case Subscribe.MarketDetail:
		tp := param.(Subscribe.MarketDetail)
		topic = fmt.Sprintf("market.%v.detail",getTradePair(tp.TradePare))
	default:
		return nil, errors.New("wrong type, start fail")
	}
	reqMsg := hbws.ReqMessage{
		Command: command,
		Topic:   topic,

	}
	newWorker.ReqCh <- reqMsg
	return newWorker.RespCh, nil

}
//request data 10times per second automatically
func RunReqWork(param interface{}) (chan interface{}, error){
	newWorker := new(hbws.WsWorker)
	newWorker.ReqCh = make(chan hbws.ReqMessage,1000)
	newWorker.RespCh = make(chan interface{},1000)
	go hbws.StartAWebSocketClient(newWorker)
	command := "req"
	topic := ""
	switch param.(type) {
	case Request.Kline:
		tp := param.(Request.Kline)
		topic = fmt.Sprintf("market.%v.kline.%vmin",getTradePair(tp.TradePare), tp.Period)

	case Request.TradeDetail:
		tp := param.(Request.TradeDetail)
		topic = fmt.Sprintf("market.%v.trade.detail",getTradePair(tp.TradePare))
	case Request.Depth:
		tp := param.(Request.Depth)
		topic = fmt.Sprintf("market.%v.depth.step%v",getTradePair(tp.TradePare),tp.Step)
	case Request.BBO:
		tp := param.(Request.BBO)
		topic = fmt.Sprintf("market.%v.bbo",getTradePair(tp.TradePare))
	case Request.MarketDetail:
		tp := param.(Request.MarketDetail)
		topic = fmt.Sprintf("market.%v.detail",getTradePair(tp.TradePare))
	default:
		return nil, errors.New("wrong type, start fail")
	}
	reqMsg := hbws.ReqMessage{
		Command: command,
		Topic:   topic,

	}
	go requestPolling(newWorker.ReqCh, reqMsg)
	return newWorker.RespCh, nil
}

func getTradePair(reqT interface{}) string {
	switch reqT.(type) {
	case Request.TradePare:
		req := reqT.(Request.TradePare)
		return fmt.Sprintf("%v%v",req.Coin, req.BaseCoin)
	case Subscribe.TradePare:
		req := reqT.(Subscribe.TradePare)
		return fmt.Sprintf("%v%v",req.Coin, req.BaseCoin)
	default:
		log.Println("wrong trade pare type")
		return ""
	}
}
func requestPolling(reqCH chan hbws.ReqMessage, reqMsg hbws.ReqMessage) {
	for {
		reqCH <- reqMsg
		time.Sleep(time.Second/10)
	}
}
