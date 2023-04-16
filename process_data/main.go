package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type JijinData struct {
	date     string
	value    float64
	buy      float64
	buy_sum  float64
	cost     float64
	left     float64
	sell     float64
	sell_sum float64
}

var output []JijinData

var sellThreshold float64 = 1.07
var buyPerTime float64 = 200

var fileName = "../fetch_data/005827.txt"
var outFileName = "005827.csv"

func loadData() {
	fileHandle, _ := os.OpenFile(fileName, os.O_RDONLY, 0777)
	defer fileHandle.Close()
	sc := bufio.NewScanner(fileHandle)

	for sc.Scan() {
		line := sc.Text()
		world := strings.Split(line, " ")
		if len(world) < 2 {
			continue
		}

		value, err := strconv.ParseFloat(world[1], 64)
		if err != nil {
			fmt.Println(err)
			return
		}

		data := JijinData{
			date:  world[0],
			value: value,
		}
		output = append([]JijinData{data}, output...)
	}

	//for _, v := range output {
	//	fmt.Println(v)
	//}
}

func buyJijin() {
	var cost float64 = 0.0
	for i, _ := range output {
		if i%5 == 0 {
			output[i].buy = buyPerTime
			cost += buyPerTime
		} else {
			output[i].buy = 0
		}
		output[i].cost = cost
	}
	//for _, v := range output {
	//	fmt.Println(v)
	//}
}

func sellJijin() {
	var sell_sum float64 = 0.0
	var left float64 = 0.0
	var buy_sum float64 = 0.0
	for i, v := range output {
		var sell float64 = 0.0

		if i != 0 {
			left = left * output[i].value / output[i-1].value
		}

		if left/buy_sum >= sellThreshold {
			sell = left
			sell_sum += sell
			left = 0
			buy_sum = 0
		}

		left += v.buy
		buy_sum += v.buy

		output[i].buy_sum = buy_sum
		output[i].left = left
		output[i].sell = sell
		output[i].sell_sum = sell_sum
	}
	//for _, v := range output {
	//	fmt.Println(v)
	//}
}

func writeCsv() {
	var s string = "日期,净值,购买,购买净值,总成本,剩余,出售,出售总值\n"
	add := func(v float64) {
		s += "," + fmt.Sprintf("%f", v)
	}
	for _, v := range output {
		s += v.date
		add(v.value)
		add(v.buy)
		add(v.buy_sum)
		add(v.cost)
		add(v.left)
		add(v.sell)
		add(v.sell_sum)
		s += "\n"
	}
	_ = os.WriteFile(outFileName, []byte(s), 0644)
}

func main() {
	loadData()
	buyJijin()
	sellJijin()
	writeCsv()
}
