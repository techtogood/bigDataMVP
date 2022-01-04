package service

import (
	"dataSource/util"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"strings"
)

type QuotationResult struct {
	Code        string
	Name        string
	Date        string
	CurPrice    float64
	OpenPrice   float64
	LastPrice   float64
	HighPrice   float64
	LowPrice    float64
	Quota       float64
	Percent     float64
	Rate        float64
	Amount      float64
	MoneyAmount      float64
}

const QuotationUrl string = "http://qt.gtimg.cn/q="

func (r QuotationResult) Request( code string ) (QuotationResult, error) {
	url := fmt.Sprintf(QuotationUrl + code)
	fmt.Println(url)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return r, err
	}

	client := &http.Client{}
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	utf8Reader := transform.NewReader(resp.Body, simplifiedchinese.GBK.NewDecoder())
	body, _ := ioutil.ReadAll(utf8Reader)
	realBody := r.getBody(body)
	if len(realBody) > 1{
		fmt.Println(realBody)
		records := strings.Split(realBody, "~")
		return r.ConvertResult(records), nil
	}

	return r,errors.New("code:{code} no quotation")
}

func (r QuotationResult) getBody(body []byte) string {
	startIndex,endIndex:= 0,0

	for i,v := range body {
		if v == '"'{
			if startIndex == 0 {
				startIndex = i
			}else{
				endIndex = i
			}
		}
	}
	fmt.Println(startIndex,endIndex)
	if endIndex - startIndex > 2{
		return string(body[startIndex+1: endIndex])
	} else{
		return ""
	}
}

func (r QuotationResult) ConvertResult(records []string) QuotationResult {


	results := QuotationResult{
		Date:        records[30],
		Code:        records[2],
		Name:        records[1],
		CurPrice:    util.ConvertFloat64(records[3]),
		LastPrice:   util.ConvertFloat64(records[4]),
		HighPrice:   util.ConvertFloat64(records[33]),
		LowPrice:    util.ConvertFloat64(records[34]),
		OpenPrice:   util.ConvertFloat64(records[5]),
		Quota:       util.ConvertFloat64(records[31]),
		Percent:     util.ConvertFloat64(records[32]),
		Rate:        util.ConvertFloat64(records[38]),
		Amount:      util.ConvertFloat64(records[36]),
		MoneyAmount: util.ConvertFloat64(records[37]),
	}

	return results
}
