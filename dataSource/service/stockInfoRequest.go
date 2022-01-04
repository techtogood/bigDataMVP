package service

import (
	"bytes"
	"dataSource/util"
	"encoding/csv"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"net/http"
)

type Result struct {
	Date        string
	Code        string
	Name        string
	ClosePrice  float64
	HighPrice   float64
	LowPrice    float64
	OpenPrice   float64
	LastPrice   float64
	Quota       float64
	Percent     float64
	Rate        float64
	Amount      float64
	MoneyAmount float64
	TotalValue  float64
	MarketValue float64
}

const CrawlUrl string = "http://quotes.money.163.com/service/chddata.html"

func (r Result) Request(code string, startDate string, endDate string) ([]Result, error) {
	body := &bytes.Buffer{}
	url := fmt.Sprintf(CrawlUrl+"?code=%s&start=%s&end=%s", code, startDate, endDate)
	//fmt.Println(url,body)
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, _ := client.Do(request)
	//fmt.Println(resp.Body)
	defer resp.Body.Close()

	utf8Reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	csvRead := csv.NewReader(utf8Reader)
	records, _ := csvRead.ReadAll()

	return r.ConvertResult(records), nil
}

func (r Result) ConvertResult(records [][]string) []Result {
	var results []Result
	for i, record := range records {
		if i == 0 {
			continue
		}

		trend := Result{
			Date:        record[0],
			Code:        record[1],
			Name:        record[2],
			ClosePrice:  util.ConvertFloat64(record[3]),
			HighPrice:   util.ConvertFloat64(record[4]),
			LowPrice:    util.ConvertFloat64(record[5]),
			OpenPrice:   util.ConvertFloat64(record[6]),
			LastPrice:   util.ConvertFloat64(record[7]),
			Quota:       util.ConvertFloat64(record[8]),
			Percent:     util.ConvertFloat64(record[9]),
			Rate:        util.ConvertFloat64(record[10]),
			Amount:      util.ConvertFloat64(record[11]),
			MoneyAmount: util.ConvertFloat64(record[12]),
			TotalValue:  util.ConvertFloat64(record[13]),
			MarketValue: util.ConvertFloat64(record[14]),
		}
		results = append(results, trend)
	}

	return results
}
