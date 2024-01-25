package entity

import "time"

type Product struct {
	ID               uint64          `json:"id"`
	Name             string          `json:"name"`
	Sku              string          `json:"sku"`
	Slug             string          `json:"slug"`
	Barcode          string          `json:"barcode"`
	BrandCode        string          `json:"brand_code"`
	ProductGroupID   int             `json:"product_group_id"`
	MetaTitle        string          `json:"meta_title"`
	MetaDescription  string          `json:"meta_description"`
	MetaKeyword      string          `json:"meta_keyword"`
	Description      string          `json:"description"`
	ShortDescription string          `json:"short_description"`
	Weight           float64         `json:"weight"`
	BasePoint        int             `json:"base_point"`
	BasePointRupiah  float64         `json:"base_point_rupiah"`
	AttributeSetID   int             `json:"attribute_set_id"`
	TypeID           string          `json:"type_id"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	Status           int             `json:"status"`
	Length           float64         `json:"length"`
	Height           float64         `json:"height"`
	Width            float64         `json:"width"`
	MetaTitleH1      string          `json:"meta_title_h1"`
	UomCode          string          `json:"uom_code"`
	BasePrice        int             `json:"base_price"`
	ProductImage     *[]ProductImage `gorm:"foreignKey:product_id;references:id"`
}

func (pi Product) TableName() string {
	return "product"
}
