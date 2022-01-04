package main

import (
	"dataSource/model"
	"dataSource/service"
	"dataSource/util"
	"strconv"
	"time"
	"fmt"
)

//获取
func InitData(){
	//清空数据
	db, err := model.GormOpenDB()
	if err != nil {
		return
	}

	sql := fmt.Sprintf("TRUNCATE TABLE %s", model.Stock{}.TableName())
	if err := db.Exec(sql).Error; err != nil {
		db.Close()
		return
	}

	sql = fmt.Sprintf("TRUNCATE TABLE %s", model.StockInfo{}.TableName())
	if err := db.Exec(sql).Error; err != nil {
		db.Close()
		return
	}
	db.Close()

	for code:=1300000;code<1300013;code++{
		service.GetStockInfoImpl(strconv.Itoa(code),"20211008","20211008")
	}

	task := service.Task{StartDate: "20211008", EndDate: util.GetNowDate()}
	task.Run()
}

func getStockInfo(){
	task := service.Task{StartDate: util.GetNowDate(), EndDate: util.GetNowDate()}
	task.Run()
}

func getStockQuotation(){
	task := service.QuotationTask{}
	task.Run()
}

func truncateStockQuotation(){
	task := service.QuotationTask{}
	task.TruncateData()
}

func main() {
	InitData()
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Println("main ticker heartbeat...")
			if util.IsWorkingDay(){
				if util.IsUpdateDailyTime(){
					fmt.Println("start to update daily stock data")
					getStockInfo()
				}
				if util.IsTruncateDailyTime(){
					fmt.Println("truncate stock quotation table")
					truncateStockQuotation()
				}
				if util.IsTradingTime() {
					fmt.Println("start to get trading data")
					getStockQuotation()
				}
			}
		}
	}

}

