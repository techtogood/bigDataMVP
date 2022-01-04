package model

import "time"

type Stock struct {
	Id        uint32 `gorm:"column:id;type:int(11);PRIMARY_KEY;AUTO_INCREMENT;"`
	Code      string `gorm:"column:code;type:varchar(10);DEFAULT:0;NOT NULL;"`
	Name      string `gorm:"column:name;type:varchar(20);NOT NULL;"`
	Type      uint8  `gorm:"column:type;type:tinyint(1);DEFAULT:0;NOT NULL;"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (s Stock) TableName() string {
	return "stock"
}

func (s Stock) FindAll() ([]Stock, error) {
	db, err := GormOpenDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var stocks []Stock
	if err := db.Order("id ASC").Find(&stocks).Error; err != nil {
		return nil, err
	}

	return stocks, nil
}