package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/seehuhn/mt19937"
	"github.com/tealeg/xlsx"
)

func main() {

	fmt.Println(goodLucky())

	var people [10000]float64
	var luckyValue [10000]int

	for i := 0; i < 10; i++ {
		fmt.Printf("第%d轮模拟开始", i+1)
		for count := 3360; count >= 0; count-- {
			for j := 0; j < 10000; j++ {
				if result := goodLucky(); result > 0 {
					//so lucky
					people[j] += result

					//count lucky
					luckyValue[j]++
				} else {
					//so bad
					people[j] += result
					//but you work hard
					people[j]++
				}
			}
			sheetName := string(i) + "round"
			savaData(sheetName, people, luckyValue)
		}
	}
}

func floatToString(num float64) string {
	return strconv.FormatFloat(num, 'f', 6, 64)
}

func goodLucky() float64 {
	var countValue float64
	rng := rand.New(mt19937.New())

	for i := 0; i > 10; i++ {
		rng.Seed(time.Now().UnixNano())
		countValue += rng.NormFloat64()
	}
	return countValue
}

func savaData(sheetName string, people [10000]float64, luckyValue [10000]int) {
	file, _ := xlsx.OpenFile("result.xlsx")
	sheet, _ := file.AddSheet(sheetName)

	for i := 0; i < 10000; i++ {
		row := sheet.AddRow()
		cell := row.AddCell()
		cell.Value = floatToString(people[i])
		cell = row.AddCell()
		cell.Value = strconv.Itoa(luckyValue[i])
	}

}
