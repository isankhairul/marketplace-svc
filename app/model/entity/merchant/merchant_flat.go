package entity

import "time"

type MerchantFlat struct {
	ID                int       `json:"id"`
	Code              string    `json:"code"`
	Name              string    `json:"name"`
	Slug              string    `json:"slug"`
	Status            int       `json:"status"`
	TypeID            int       `json:"type_id"`
	TypeSlug          string    `json:"type_slug"`
	Stock             int       `json:"stock"`
	ReservedStock     int       `json:"reserved_stock"`
	StockOnHand       int       `json:"stock_on_hand"`
	MaxPurchaseQty    int       `json:"max_purchase_qty"`
	Sold              int       `json:"sold"`
	ProvinceID        int       `json:"province_id"`
	Province          string    `json:"province"`
	CityID            int       `json:"city_id"`
	City              string    `json:"city"`
	DistrictID        int       `json:"district_id"`
	District          string    `json:"district"`
	SubdistrictID     int       `json:"subdistrict_id"`
	Subdistrict       string    `json:"subdistrict"`
	PostalcodeID      int       `json:"postalcode_id"`
	Zipcode           string    `json:"zipcode"`
	Image             string    `json:"image"`
	Categories        any       `json:"categories"`
	MerchantProductID int       `json:"merchant_product_id"`
	MerchantSku       string    `json:"merchant_sku"`
	ProductSku        string    `json:"product_sku"`
	ProductStatus     int       `json:"product_status"`
	Rating            int       `json:"rating"`
	Review            int       `json:"review"`
	SellingPrice      int       `json:"selling_price"`
	SpecialPrices     string    `json:"special_prices"`
	HidePrice         bool      `json:"hide_price"`
	UpdatedAt         time.Time `json:"updated_at"`
	MerchantUID       string    `json:"merchant_uid"`
	Latitude          string    `json:"latitude"`
	Longitude         string    `json:"longitude"`
}

func (m MerchantFlat) TableName() string {
	return "merchant_flat_"
}
