package main

import (
	"fmt"
	_ "hbProject/hb/controller"
)
var (
	arr []float64
	cur15time []float64
	allarr []float64
	priceArr []float64
	start int = 1000
	longTrend int
	shortTrend int
)

func main() {
	fmt.Println("main start")
	defer fmt.Println("main stop")
	select {
	}
	/*var count int
    for {
		kline, err := hb.GetCurTrade("btcusdt")
		if err != nil {
			log.Print(err)
		} else {
			var ltArr []float64
			//var stArr []float64
			for _, x := range kline.Data {
				if len(priceArr) < 100 {
					priceArr = append(priceArr, x.Price)
				} else {
					priceArr = append(priceArr, x.Price)
					priceArr = priceArr[1:]
					a := (priceArr[0] + priceArr[1] + priceArr[2])/3
					ltArr = append(ltArr,a)
					a = (priceArr[19] + priceArr[20] + priceArr[21])/3
					ltArr = append(ltArr,a)
					a = (priceArr[39] + priceArr[40] + priceArr[41])/3
					ltArr = append(ltArr,a)
					a = (priceArr[59] + priceArr[60] + priceArr[61])/3
					ltArr = append(ltArr,a)
					a = (priceArr[89] + priceArr[90] + priceArr[91])/3
					ltArr = append(ltArr,a)
				}
				var u float64
				for i := 0; i < len(ltArr) - 1; i ++ {
					u += ltArr[i+1] -ltArr[i]
				}
				if u > 0 {
					longTrend = 1
				} else {
					longTrend = 0
				}
				allarr = append(allarr,x.Amount)
				if len(allarr) < 15 {
					var a float64
					for _, x := range allarr {
						a += x
					}

				} else {
					var a float64
					for _, x := range allarr {
						a += x
					}
					cur15time = append(cur15time, a)
					allarr = allarr[1:]
				}
				if len(cur15time) > 10 {
					cur15time = cur15time[1:]
				}

				if x.Amount > 100 {
					if len(arr) < 10 {
						arr = append(arr, x.Amount)
						fmt.Println("big trade happened")
						fmt.Println("arr:",arr)
					} else {
						arr = append(arr, x.Amount)
						fmt.Println("big trade happened")
						fmt.Println("arr:",arr[len(arr)-10:])
					}

				}

			}
		}
		count ++
		if count == 10 {
			fmt.Println("cur15time:",
				cur15time[len(cur15time)-1],
				"--",cur15time[len(cur15time)-2],
				"--",cur15time[len(cur15time)-3])
			if longTrend == 1 {
				fmt.Println("up")
			} else {
				fmt.Println("down")
			}
			count = 0
		}
		time.Sleep(time.Second/10)
	}*/


}
