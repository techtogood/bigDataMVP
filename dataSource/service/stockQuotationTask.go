package service

import (
	"bytes"
	"dataSource/model"
	"fmt"
	"github.com/jinzhu/gorm"
	"math"
	"strconv"
	"time"
)

type QuotationTask struct {
	ExecStartTime int64
}

var quotationRequestsLimit = 200

func (t *QuotationTask) Run(){
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

func (t *QuotationTask) Worker(db *gorm.DB, ticker *time.Ticker) {
	stocks,err := new(model.Stock).FindAll()
	if err !=nil {
		return
	}

	stocksLen := len(stocks)
	fmt.Println(stocks,stocksLen)
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
		getChan := make(chan QuotationResult, sliceStockLen)
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

func (t *QuotationTask) GetData(value model.Stock ,getChan chan QuotationResult ){

	var requestCode = ""
	if value.Type == 0 {
		requestCode = "sh"+ value.Code
	}else if value.Type == 1 {
		requestCode = "sz"+ value.Code
	}

	result, _ := QuotationResult{}.Request(requestCode)
	fmt.Println(result)
	getChan <- result
	return
}

func (t QuotationTask) StoreTask(db *gorm.DB, getChan chan QuotationResult, storeChan chan int) {
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

func (t *QuotationTask) StoreData(db *gorm.DB, result QuotationResult ,storeChan chan int ) error {


	var buffer bytes.Buffer
	sql := fmt.Sprintf("INSERT INTO %s (`code`, `name`, `datetime`, `cur_price`, `open_price`, `last_price`, `high_price`, `low_price`, `quota`, `percent`, `rate`, `amount`, `money_amount`) VALUES", model.StockQuotation{}.TableName())
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}

	buffer.WriteString(fmt.Sprintf("('%s', '%s', %s, %f, %f, %f, %f, %f, %f, %f, %f, %f, %f) ON DUPLICATE KEY UPDATE `name`=values(`name`);", result.Code, result.Name, result.Date, result.CurPrice, result.OpenPrice, result.LastPrice, result.HighPrice, result.LowPrice, result.Quota, result.Percent, result.Rate, result.Amount, result.MoneyAmount))

	//fmt.Println(buffer.String())

	if err := db.Exec(buffer.String()).Error; err != nil {
		storeChan <- 1
		return err
	}

	storeChan <- 1

	return nil
}

func (t *QuotationTask) TruncateData() error {
	db, err := model.GormOpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	sql := fmt.Sprintf("TRUNCATE TABLE %s", model.StockQuotation{}.TableName())
	if err := db.Exec(sql).Error; err != nil {
		return err
	}
	return nil

}
