package entity

type ProductCategory struct {
	ProductID  uint64      `json:"product_id"`
	CategoryID uint64      `json:"category_id"`
	StoreID    uint64      `json:"store_id"`
	Product    *[]Product  `gorm:"foreignKey:id;references:product_id"`
	Category   *[]Category `gorm:"foreignKey:id;references:category_id"`
}

func (pi ProductCategory) TableName() string {
	return "product_category"
}
