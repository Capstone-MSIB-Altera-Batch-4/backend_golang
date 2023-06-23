package gen

import (
	"fmt"
	"point-of-sale/app/model"
	"point-of-sale/config"
	"time"
)

func GetOrderCount() int {
	today := time.Now().Format("02012006")
	var count int64
	config.Db.Model(&model.Order{}).Where("DATE_FORMAT(created_at, '%d%m%Y') = ?", today).Count(&count)
	return int(count)
}

func GenerateOrderCode(count int, date string) string {
	orderCount := count + 1
	orderCode := fmt.Sprintf("%02d-%s", orderCount, date)
	return orderCode
}
