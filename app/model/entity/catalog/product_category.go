package entity

type ProductCategory struct {
	ProductID  int64       `json:"product_id"`
	CategoryID int64       `json:"category_id"`
	StoreID    int64       `json:"store_id"`
	Product    *[]Product  `gorm:"foreignKey:product_id;references:id"`
	Category   *[]Category `gorm:"foreignKey:category_id;references:id"`
}

func (pi ProductCategory) TableName() string {
	return "product_category"
}
