package entity

import (
	entity "marketplace-svc/app/model/entity/merchant"
	"time"
)

type ProductFlat struct {
	ID                       int32     `json:"id"`
	Name                     string    `json:"name"`
	SKU                      string    `json:"sku"`
	Slug                     string    `json:"slug"`
	Barcode                  string    `json:"barcode"`
	BrandCode                string    `json:"brand_code"`
	PrincipalCode            any       `json:"principal_code"`
	ProductGroupID           int32     `json:"product_group_id"`
	MetaTitle                string    `json:"meta_title"`
	MetaDescription          string    `json:"meta_description"`
	MetaKeyword              string    `json:"meta_keyword"`
	Description              string    `json:"description"`
	ShortDescription         string    `json:"short_description"`
	Weight                   float64   `json:"weight"`
	BasePoint32              int32     `json:"base_point32"`
	BasePoint32Rupiah        float64   `json:"base_point32_rupiah"`
	AttributeSetID           int32     `json:"attribute_set_id"`
	TypeID                   string    `json:"type_id"`
	Status                   int32     `json:"status"`
	DeletedAt                any       `json:"deleted_at"`
	Length                   float64   `json:"length"`
	Height                   float64   `json:"height"`
	Width                    float64   `json:"width"`
	MetaTitleH1              string    `json:"meta_title_h1"`
	CustomerGroupID          int32     `json:"customer_group_id"`
	TaxClassID               string    `json:"tax_class_id"`
	BasePrice                float64   `json:"base_price"`
	FinalPrice               any       `json:"final_price"`
	MinPrice                 float64   `json:"min_price"`
	MaxPrice                 float64   `json:"max_price"`
	TierPrice                any       `json:"tier_price"`
	GroupPrice               any       `json:"group_price"`
	SpecialPrice             any       `json:"special_price"`
	SpecialFromDate          any       `json:"special_from_date"`
	SpecialToDate            any       `json:"special_to_date"`
	StoreID                  int32     `json:"store_id"`
	CategoryIds              any       `json:"category_ids"`
	UpperLimitPrice          any       `json:"upper_limit_price"`
	Images                   string    `json:"images"`
	MaximumPurchaseQuantity  string    `json:"maximum_purchase_quantity"`
	ProductMetaDescription   string    `json:"product_meta_description"`
	Visibility               string    `json:"visibility"`
	ManagementFee            string    `json:"management_fee"`
	IsFreeProduct            int32     `json:"is_free_product"`
	Preorder                 int32     `json:"preorder"`
	ProductKn                int32     `json:"product_kn"`
	ProductKalbe             int32     `json:"product_kalbe"`
	IsSpot                   int32     `json:"is_spot"`
	IsActive                 int32     `json:"is_active"`
	Saleable                 int32     `json:"saleable"`
	RewardPoint32SellProduct string    `json:"reward_point32_sell_product"`
	EpmPrice                 float64   `json:"epm_price"`
	ProductMetaTitle         string    `json:"product_meta_title"`
	ProductSlug              string    `json:"product_slug"`
	CbpPrice                 float64   `json:"cbp_price"`
	IsFamilyGift             int32     `json:"is_family_gift"`
	IsLangganan              int32     `json:"is_langganan"`
	IsTicket                 int32     `json:"is_ticket"`
	EpmSku                   string    `json:"epm_sku"`
	NewsFromDate             any       `json:"news_from_date"`
	NewsToDate               any       `json:"news_to_date"`
	ProductDescription       string    `json:"product_description"`
	IsKliknow                int32     `json:"is_kliknow"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`

	// MerchantFlat
	// in: merchant_flat
	MerchantFlat *entity.MerchantFlat `gorm:"foreignKey:id;references:merchant_product_id" json:"merchant_flat"`

	// MerchantFlats
	// in: merchant_flats
	//MerchantFlats []entity.MerchantFlat `gorm:"foreignKey:id;references:merchant_product_id" json:"merchant_flats"`
}
