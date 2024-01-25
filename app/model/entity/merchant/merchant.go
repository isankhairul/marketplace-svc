package entity

import "time"

type Merchant struct {
	ID                       uint64             `json:"id"`
	MerchantCode             string             `json:"merchant_code"`
	CreatedAt                time.Time          `json:"created_at"`
	UpdatedAt                time.Time          `json:"updated_at"`
	Status                   int                `json:"status"`
	Address                  string             `json:"address"`
	MerchantName             string             `json:"merchant_name"`
	ProvinceID               int                `json:"province_id"`
	CityID                   int                `json:"city_id"`
	DistrictID               int                `json:"district_id"`
	SubdistrictID            int                `json:"subdistrict_id"`
	PostalcodeID             int                `json:"postalcode_id"`
	Email                    string             `json:"email"`
	MerchantTypeID           uint64             `json:"merchant_type_id"`
	Image                    string             `json:"image"`
	Slug                     string             `json:"slug"`
	PhoneNumber              string             `json:"phone_number"`
	Coordinate               string             `json:"coordinate"`
	NotifWhatsapp            int                `json:"notif_whatsapp"`
	SubdistCode              string             `json:"subdist_code"`
	Latitude                 float64            `json:"latitude"`
	Longitude                float64            `json:"longitude"`
	MerchantWarehouseDefault int                `json:"merchant_warehouse_default"`
	MetaTitle                string             `json:"meta_title"`
	MetaDescription          string             `json:"meta_description"`
	MerchantFlaggingID       int                `json:"merchant_flagging_id"`
	HidePrice                bool               `json:"hide_price"`
	FulfillmentID            int                `json:"fulfillment_id"`
	FulfillmentCode          string             `json:"fulfillment_code"`
	MerchantUID              string             `json:"merchant_uid"`
	CustomerContactIds       string             `json:"customer_contact_ids"`
	IsOfficial               int                `json:"is_official"`
	IsPharmacy               int                `json:"is_pharmacy"`
	UpPrice                  int                `json:"up_price"`
	MerchantInfo             MerchantInfo       `gorm:"foreignKey:merchant_id;references:id" json:"-"`
	MerchantType             MerchantType       `gorm:"foreignKey:merchant_type_id;references:id" json:"-"`
	MerchantShipping         []MerchantShipping `gorm:"foreignKey:merchant_id;references:id" json:"-"`
	MerchantStores           []MerchantStores   `gorm:"foreignKey:merchant_id;references:id;"`
	MerchantProduct          []MerchantProduct  `gorm:"foreignKey:merchant_id;references:id;"`
}

func (m Merchant) TableName() string {
	return "merchant"
}
