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

type EsProductFlat struct {
	ID                       string  `json:"id"`
	Name                     string  `json:"name"`
	Sku                      string  `json:"sku"`
	Slug                     string  `json:"slug"`
	Barcode                  string  `json:"barcode"`
	BrandCode                string  `json:"brand_code"`
	MetaTitle                string  `json:"meta_title"`
	MetaTitleH1              string  `json:"meta_title_h1"`
	MetaDescription          string  `json:"meta_description"`
	MetaKeyword              string  `json:"meta_keyword"`
	ShortDescription         string  `json:"short_description"`
	Description              string  `json:"description"`
	CompletionTerms          string  `json:"completion_terms"`
	FullTextSearch           string  `json:"full_text_search"`
	Weight                   float64 `json:"weight"`
	BasePoint32              float64 `json:"base_point32"`
	BasePoint32Rupiah        int32   `json:"base_point32_rupiah"`
	RewardPoint32SellProduct int32   `json:"reward_point32_sell_product"`
	IsFamilyGift             int32   `json:"is_family_gift"`
	IsFreeProduct            int32   `json:"is_free_product"`
	IsLangganan              int32   `json:"is_langganan"`
	IsSpot                   int32   `json:"is_spot"`
	IsTicket                 int32   `json:"is_ticket"`
	IsKliknow                int32   `json:"is_kliknow"`
	IsActive                 int32   `json:"is_active"`
	TypeID                   string  `json:"type_id"`
	Status                   int32   `json:"status"`
	CreatedAt                string  `json:"created_at"`
	UpdatedAt                string  `json:"updated_at"`
	Images                   []struct {
		ID        int32  `json:"id"`
		Thumbnail string `json:"thumbnail"`
		Default   string `json:"default"`
		Original  string `json:"original"`
		IsDefault int32  `json:"is_default"`
	} `json:"images"`
	Breadcrumbs []struct {
		Level int32  `json:"level"`
		Name  string `json:"name"`
		Slug  string `json:"slug"`
	} `json:"breadcrumbs"`
	Categories []EsProductFlatCategory `json:"categories"`
	Merchants  EsProductFlatMerchant   `json:"merchants"`
}

type EsProductFlatCategory struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
	InHome int32  `json:"in_home"`
}

type EsProductFlatMerchant struct {
	ID                int32                       `json:"id"`
	UID               string                      `json:"uid"`
	Code              string                      `json:"code"`
	Name              string                      `json:"name"`
	Slug              string                      `json:"slug"`
	Status            int32                       `json:"status"`
	TypeID            int32                       `json:"type_id"`
	TypeSlug          string                      `json:"type_slug"`
	Stock             int32                       `json:"stock"`
	ReservedStock     int32                       `json:"reserved_stock"`
	StockOnHand       int32                       `json:"stock_on_hand"`
	MaxPurchaseQty    int32                       `json:"max_purchase_qty"`
	Sold              int32                       `json:"sold"`
	ProvinceID        int32                       `json:"province_id"`
	Province          string                      `json:"province"`
	CityID            int32                       `json:"city_id"`
	City              string                      `json:"city"`
	DistrictID        int32                       `json:"district_id"`
	District          string                      `json:"district"`
	SubdistrictID     int32                       `json:"subdistrict_id"`
	Subdistrict       string                      `json:"subdistrict"`
	PostalcodeID      int32                       `json:"postalcode_id"`
	Zipcode           string                      `json:"zipcode"`
	Location          EsProductFlatLocation       `json:"location"`
	Image             string                      `json:"image"`
	Categories        []string                    `json:"categories"`
	MerchantProductID int32                       `json:"merchant_product_id"`
	MerchantSku       string                      `json:"merchant_sku"`
	ProductSku        string                      `json:"product_sku"`
	ProductStatus     int32                       `json:"product_status"`
	Rating            float64                     `json:"rating"`
	Review            int32                       `json:"review"`
	SellingPrice      int32                       `json:"selling_price"`
	SpecialPrices     []EsProductFlatSpecialPrice `json:"special_prices"`
	HidePrice         bool                        `json:"hide_price"`
}

type EsProductFlatLocation struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type EsProductFlatSpecialPrice struct {
	Price           int32  `json:"price"`
	ToTime          string `json:"to_time"`
	FromTime        string `json:"from_time"`
	CustomerGroupID int32  `json:"customer_group_id"`
}
