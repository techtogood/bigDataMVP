package util

import (
	"time"
)

func IsWorkingDay() bool {
	switch time.Now().Weekday().String() {
	case "Monday","Tuesday","Wednesday","Thursday","Friday":
		return true
	case "Saturday","Sunday":
		return false
	}
	return false
}


func IsTradingTime() bool {
	h,m,s := time.Now().Clock()
	clock := h*10000+m*100+s
	if clock >= 93000 && clock <= 113000 ||  clock >= 130000 && clock <= 150000 {
		return  true
	}else{
		return false
	}
}
func IsUpdateDailyTime() bool{
	h,m,_ := time.Now().Clock()
	if h == 17 && (m == 0 || m == 1 || m == 2){
		return true
	}
	return false
}

func IsTruncateDailyTime() bool{
	h,m,_ := time.Now().Clock()
	if h == 23 && (m == 0 || m == 1){
		return true
	}
	return false
}

func GetNowTimeInt() int {
	h,m,s := time.Now().Clock()
	return h*10000+m*100+s
}

func GetWeekday() string{
	return time.Now().Weekday().String()
}

func GetNowDate() string{
	return time.Now().Format("20060102")
}