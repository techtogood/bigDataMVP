package model

type StockInfo struct {
	Id           uint    `gorm:"column:id;type:int(11);PRIMARY_KEY;AUTO_INCREMENT;"`
	Code         string  `gorm:"column:code;type:varchar(10);DEFAULT:0;NOT NULL;"`
	Name         string  `gorm:"column:name;type:varchar(20);NOT NULL;"`
	OpenPrice    float64 `gorm:"column:open_price;type:float;DEFAULT:0;NOT NULL;"`
	ClosePrice   float64 `gorm:"column:close_price;type:float;DEFAULT:0;NOT NULL;"`
	HighPrice    float64 `gorm:"column:high_price;type:float;DEFAULT:0;NOT NULL;"`
	LowPrice     float64 `gorm:"column:low_price;type:float;DEFAULT:0;NOT NULL;"`
	LastPrice    float64 `gorm:"column:last_price;type:float;NOT NULL;"`
	Quota        float64 `gorm:"column:quota;type:float;DEFAULT:0;NOT NULL;"`
	Percent      float64 `gorm:"column:percent;type:float;DEFAULT:0;NOT NULL;"`
	Rate         float64 `gorm:"column:rate;type:float;DEFAULT:0;NOT NULL;"`
	Amount       float64 `gorm:"column:amount;type:float;DEFAULT:0;NOT NULL;"`
	MoneyAmount  float64 `gorm:"column:money_amount;type:float;DEFAULT:0;NOT NULL;"`
	TotalValue   float64 `gorm:"column:total_value;type:float;DEFAULT:0;NOT NULL;"`
	MarketValue  float64 `gorm:"column:market_value;type:float;DEFAULT:0;NOT NULL;"`
	Date         string  `gorm:"column:date;type:varchar(20);NOT NULL;"`
	CreatedAt    int64   `gorm:"column:created_at;type:int(11);DEFAULT:0;NOT NULL;"`
	UpdatedAt    int64   `gorm:"column:updated_at;type:int(11);DEFAULT:0;NOT NULL;"`
}

func (s StockInfo) TableName() string {
	return "stock_info"
}
