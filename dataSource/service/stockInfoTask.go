package service

import (
	"bytes"
	"dataSource/model"
	"dataSource/util"
	"fmt"
	"github.com/jinzhu/gorm"
	"math"
	"strconv"
	"time"
)

type Task struct {
	StartDate     string
	EndDate       string
	ExecStartTime int64
}

var requestsLimit = 300

func (t *Task) Run(){
	t.ExecStartTime = time.Now().Unix()
	db, err := model.GormOpenDB()
	if err != nil {
		return
	}
	defer db.Close()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(1000)

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	t.Worker(db,ticker)
}

func (t *Task) Worker(db *gorm.DB, ticker *time.Ticker) {
	stocks,err := new(model.Stock).FindAll()
	if err !=nil {
		return
	}

	stocksLen := len(stocks)
	//fmt.Println(stocks,stocksLen)
	taskNum := math.Ceil(float64(stocksLen) / float64(requestsLimit))
	taskIndex := 1
	for i := 0; i < stocksLen; i += requestsLimit {
		var stockSLice []model.Stock
		if taskIndex == int(taskNum) {
			stockSLice = stocks[i:]
		} else{
			stockSLice = stocks[i:taskIndex*requestsLimit]
		}

		sliceStockLen := len(stockSLice)
		getChan := make(chan []Result, sliceStockLen)
		storeChan := make(chan int, sliceStockLen)

		go t.StoreTask(db,getChan,storeChan)

		for _, value := range stockSLice {
			go t.GetData(value, getChan)
		}

	Loop:
		for {
			select {
			case <-ticker.C:
				chanLen := len(storeChan)
				if chanLen == sliceStockLen {
					fmt.Printf("Task %d completed,Total spend:%ds\n\n", taskIndex, time.Now().Unix()-t.ExecStartTime)
					break Loop
				}

				progressString := fmt.Sprintf("%.2f", float64(chanLen)/float64(sliceStockLen))
				progressFloat, _ := strconv.ParseFloat(progressString, 64)
				progress := int(progressFloat * 100)

				fmt.Printf("Handing requests %d,Current progress:%d%s\n", chanLen, progress, "%")
			}
		}

		taskIndex++

		close(getChan)
		close(storeChan)
	}


}

func (t *Task) GetData(value model.Stock ,getChan chan []Result ){

	requestCode := util.ConvertString(int(value.Type)) + value.Code
	results, _ := Result{}.Request(requestCode,t.StartDate,t.EndDate)
	fmt.Println(results)
	getChan <- results
	return
}

func (t Task) StoreTask(db *gorm.DB, getChan chan []Result, storeChan chan int) {
	for {
		select {
		case s, ok := <-getChan:
			if !ok {
				return
			}
			go t.StoreData(db, s, storeChan)
		}
	}
}

func (t *Task) StoreData(db *gorm.DB, results []Result ,storeChan chan int ) error {

	if len(results) == 0 {
		storeChan <- 1
		return fmt.Errorf("empty result")
	}

	var buffer bytes.Buffer
	sql := fmt.Sprintf("INSERT INTO %s (`code`, `name`, `open_price`, `close_price`, `high_price`, `low_price`, `last_price`, `quota`, `percent`, `rate`, `amount`, `money_amount`, `total_value`, `market_value`, `date`) VALUES", model.StockInfo{}.TableName())
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}

	//fmt.Println(len(results))

	for i, v := range results {
		if i == len(results)-1 {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, '%s') ON DUPLICATE KEY UPDATE `name`=values(`name`), `amount`=values(`amount`), `money_amount`=values(`money_amount`),`total_value`=values(`total_value`),`market_value`=values(`market_value`);", v.Code[1:], v.Name, v.OpenPrice, v.ClosePrice, v.HighPrice, v.LowPrice, v.LastPrice, v.Quota, v.Percent, v.Rate, v.Amount, v.MoneyAmount, v.TotalValue, v.MarketValue, v.Date))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s', '%s', %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f, '%s'),", v.Code[1:], v.Name, v.OpenPrice, v.ClosePrice, v.HighPrice, v.LowPrice, v.LastPrice, v.Quota, v.Percent, v.Rate, v.Amount, v.MoneyAmount, v.TotalValue, v.MarketValue, v.Date))
		}
	}

	//fmt.Println(buffer.String())

	if err := db.Exec(buffer.String()).Error; err != nil {
		return err
	}

	storeChan <- 1

	return nil
}

