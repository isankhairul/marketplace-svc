package entity

import (
	"time"
)

type ProductFlat struct {
	ID                      int64     `json:"id"`
	Name                    string    `json:"name"`
	Sku                     string    `json:"sku"`
	Slug                    string    `json:"slug"`
	Barcode                 string    `json:"barcode"`
	BrandCode               string    `json:"brand_code"`
	PrincipalCode           string    `json:"principal_code"`
	ProductGroupID          int       `json:"product_group_id"`
	MetaTitle               string    `json:"meta_title"`
	MetaDescription         string    `json:"meta_description"`
	MetaKeyword             string    `json:"meta_keyword"`
	Description             string    `json:"description"`
	ShortDescription        string    `json:"short_description"`
	Weight                  float64   `json:"weight"`
	BasePoint               int       `json:"base_point"`
	BasePointRupiah         float64   `json:"base_point_rupiah"`
	AttributeSetID          int       `json:"attribute_set_id"`
	TypeID                  string    `json:"type_id"`
	Status                  int       `json:"status"`
	Length                  float32   `json:"length"`
	Height                  float32   `json:"height"`
	Width                   float32   `json:"width"`
	MetaTitleH1             string    `json:"meta_title_h1"`
	CustomerGroupID         int       `json:"customer_group_id"`
	TaxClassID              string    `json:"tax_class_id"`
	BasePrice               float64   `json:"base_price"`
	FinalPrice              float64   `json:"final_price"`
	MinPrice                float64   `json:"min_price"`
	MaxPrice                float64   `json:"max_price"`
	TierPrice               float64   `json:"tier_price"`
	GroupPrice              float64   `json:"group_price"`
	SpecialPrice            float64   `json:"special_price"`
	SpecialFromDate         time.Time `json:"special_from_date"`
	SpecialToDate           time.Time `json:"special_to_date"`
	StoreID                 int       `json:"store_id"`
	CategoryIds             string    `json:"category_ids"`
	UpperLimitPrice         float32   `json:"upper_limit_price"`
	Images                  string    `json:"images"`
	MaximumPurchaseQuantity string    `json:"maximum_purchase_quantity"`
	IsFreeProduct           int       `json:"is_free_product"`
	Saleable                int       `json:"saleable"`
	IsKilling               int       `json:"is_killing"`
	ProductSlug             string    `json:"product_slug"`
	CbpPrice                float64   `json:"cbp_price"`
	ProductMetaDescription  string    `json:"product_meta_description"`
	Visibility              string    `json:"visibility"`
	ManagementFee           string    `json:"management_fee"`
	Preorder                int       `json:"preorder"`
	ProductKn               int       `json:"product_kn"`
	ProductKalbe            int       `json:"product_kalbe"`
	IsSpot                  int       `json:"is_spot"`
	IsActive                int       `json:"is_active"`
	RewardPointSellProduct  string    `json:"reward_point_sell_product"`
	EpmPrice                float64   `json:"epm_price"`
	ProductMetaTitle        string    `json:"product_meta_title"`
	IsFamilyGift            int       `json:"is_family_gift"`
	IsLangganan             int       `json:"is_langganan"`
	IsTicket                int       `json:"is_ticket"`
	EpmSku                  string    `json:"epm_sku"`
	ProductDescription      string    `json:"product_description"`
	IsKliknow               int       `json:"is_kliknow"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	IsPrescription          int       `json:"is_prescription"`
	Uom                     string    `json:"uom"`
	Nie                     string    `json:"nie"`
	UomName                 string    `json:"uom_name"`
	PrincipalName           string    `json:"principal_name"`
	Product                 *Product  `gorm:"foreignKey:sku;references:sku"`
}

func (pl ProductFlat) TableName() string {
	return "product_flat"
}
