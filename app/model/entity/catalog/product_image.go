package entity

type ProductImage struct {
	ID             int64  `json:"id"`
	ProductID      int    `json:"product_id"`
	Image          string `json:"image"`
	ImageOriginal  string `json:"image_original"`
	ImageThumbnail string `json:"image_thumbnail"`
	Status         bool   `json:"status"`
	SortOrder      int    `json:"sort_order"`
	IsDefault      int    `json:"is_default"`
}

func (pi ProductImage) TableName() string {
	return "product_image"
}
