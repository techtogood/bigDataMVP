package service

import (
	"bytes"
	"dataSource/model"
	"encoding/csv"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"net/http"
)


const getUrl string = "http://quotes.money.163.com/service/chddata.html"

func GetStockInfoImpl( code string, starDate string, endDate string) error{
	body := &bytes.Buffer{}
	url := fmt.Sprintf(getUrl+"?code=%s&start=%s&end=%s",code,starDate,endDate)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, _ := client.Do(request)
	//fmt.Println(resp.Body)
	defer resp.Body.Close()

	utf8Reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	csvRead := csv.NewReader(utf8Reader)
	records, _ := csvRead.ReadAll()
	//fmt.Printf("len=%d",len(records))
	//fmt.Println(records[1][1],records[1][2])
	//fmt.Println(realcode)

	if len(records) == 2 {
		realcode := records[1][1][1:]
		db, err := model.GormOpenDB()
		if err != nil {
			return err
		}
		defer db.Close()
		sql := fmt.Sprintf("INSERT INTO %s (`code`,`name`,`type`) VALUES ('%s','%s','%d') ON DUPLICATE KEY UPDATE `name`=values(`name`), `code`=values(`code`)",
			model.Stock{}.TableName(),realcode,records[1][2],1)
		fmt.Println(sql)
		if err := db.Exec(sql).Error; err != nil {
			fmt.Println("err .............")
			return err
		}
	}

	return nil

}