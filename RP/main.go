package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/seehuhn/mt19937"
	"github.com/tealeg/xlsx"
)

func main() {
	ch := make(chan bool)
	for i := 0; i < 9; i++ {
		if i != 8 {
			go goSimulation(i, false, nil)
		} else {
			go goSimulation(i, true, ch)
		}
	}
	<-ch
	fmt.Println("模拟结束")
	fmt.Scan()
}

func goSimulation(i int, checkEndChannel bool, ch chan bool) {
	var people [10000]float64
	var luckyValue [10000]int

	fmt.Printf("线程%d启动模拟\n", i+1)
	for count := 3360; count >= 0; count-- {
		for j := 0; j < 10000; j++ {
			if result := goodLucky(i); result > 0 {
				//fmt.Println(result)
				//so lucky
				people[j] += result

				//count lucky
				luckyValue[j]++
			} else {
				//fmt.Println(result)
				//so bad
				people[j] += result
				//but you work hard
				people[j]++
			}
		}
	}
	fmt.Printf("线程%d模拟结束,开始导出数据\n", i+1)
	sheetName := strconv.Itoa(i) + "round"

	for {
		if savaData(sheetName, people, luckyValue) == false {
			//如果文件被占用,等待下一个周期
			time.Sleep(time.Duration(3) * time.Second)
		} else {
			break
		}
	}

	fmt.Printf("线程%d导出数据完成\n", i+1)
	if checkEndChannel == true {
		ch <- true
	}
}

func floatToString(num float64) string {
	return strconv.FormatFloat(num, 'f', 6, 64)
}

func goodLucky(salt int) float64 {
	var countValue float64
	rng := rand.New(mt19937.New())

	for i := 0; i < 10; i++ {
		rng.Seed(time.Now().UnixNano() + int64(salt))
		countValue += rng.NormFloat64()
	}
	return countValue
}

func savaData(sheetName string, people [10000]float64, luckyValue [10000]int) bool {
	file, err := xlsx.OpenFile("result.xlsx")

	if err != nil {
		//可能是文件被占用
		fmt.Println(err)
		return false
	}

	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10000; i++ {
		row := sheet.AddRow()
		cell := row.AddCell()
		cell.Value = floatToString(people[i])
		cell = row.AddCell()
		cell.Value = strconv.Itoa(luckyValue[i])
	}
	file.Save("result.xlsx")
	return true
}
