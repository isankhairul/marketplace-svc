package responseelastic

type MerchantResponse struct {
	Code             string             `json:"code,omitempty"`
	City             string             `json:"city,omitempty"`
	MerchantStore    []MerchantStore    `json:"merchant_store,omitempty"`
	ChannelTypeID    int                `json:"channel_type_id,omitempty"`
	Rating           float64            `json:"rating,omitempty"`
	CreatedAt        string             `json:"created_at,omitempty"`
	MerchantCategory []MerchantCategory `json:"merchant_category,omitempty"`
	UID              string             `json:"uid,omitempty"`
	Province         string             `json:"province,omitempty"`
	UpdatedAt        string             `json:"updated_at,omitempty"`
	MerchantShipping []MerchantShipping `json:"merchant_shipping,omitempty"`
	ID               int                `json:"id,omitempty"`
	Slug             string             `json:"slug,omitempty"`
	Email            string             `json:"email,omitempty"`
	MerchantType     string             `json:"merchant_type,omitempty"`
	Image            string             `json:"image,omitempty"`
	Statistic        MerchantStatistic  `json:"statistic,omitempty"`
	Address          string             `json:"address,omitempty"`
	Coordinate       string             `json:"coordinate,omitempty"`
	MetaTitle        string             `json:"meta_title,omitempty"`
	TypeSlug         string             `json:"type_slug,omitempty"`
	TypeID           int                `json:"type_id,omitempty"`
	SuggestionTerms  string             `json:"suggestion_terms,omitempty"`
	PostalcodeID     int                `json:"postalcode_id,omitempty"`
	SubdistrictID    int                `json:"subdistrict_id,omitempty"`
	Zipcode          string             `json:"zipcode,omitempty"`
	MetaDescription  string             `json:"meta_description,omitempty"`
	ProvinceID       int                `json:"province_id,omitempty"`
	Subdistrict      string             `json:"subdistrict,omitempty"`
	District         string             `json:"district,omitempty"`
	Name             string             `json:"name,omitempty"`
	PhoneNumber      string             `json:"phone_number,omitempty"`
	Location         MerchantLocation   `json:"location,omitempty"`
	DistrictID       int                `json:"district_id,omitempty"`
	PicPhoneNumber   string             `json:"pic_phone_number,omitempty"`
	Status           int                `json:"status,omitempty"`
	CityID           int                `json:"city_id,omitempty"`
}

type MerchantStore struct {
	Code string `json:"code,omitempty"`
	Name string `json:"name,omitempty"`
	ID   int    `json:"id,omitempty"`
}

type MerchantCategory struct {
	Name string `json:"name,omitempty"`
	ID   int    `json:"id,omitempty"`
}

type MerchantShipping struct {
	Name string `json:"name,omitempty"`
	Logo string `json:"logo,omitempty"`
}

type MerchantStatistic struct {
	Sold          int                 `json:"sold,omitempty"`
	AverageRating float64             `json:"average_rating,omitempty"`
	OrderReview   MerchantOrderReview `json:"order_review,omitempty"`
	ProductRating struct {
		Num1 int `json:"1,omitempty"`
		Num2 int `json:"2,omitempty"`
		Num3 int `json:"3,omitempty"`
		Num4 int `json:"4,omitempty"`
		Num5 int `json:"5,omitempty"`
	} `json:"product_rating,omitempty"`
}

type MerchantLocation struct {
	Lon string `json:"lon,omitempty"`
	Lat string `json:"lat,omitempty"`
}

type MerchantOrderReview struct {
	Num1 int `json:"1,omitempty"`
	Num2 int `json:"2,omitempty"`
	Num3 int `json:"3,omitempty"`
}

// swagger:model MerchantProductResponse
type MerchantProductResponse struct {
	ID             float64                `json:"id"`
	UID            string                 `json:"uid"`
	Name           string                 `json:"name"`
	Type           string                 `json:"type"`
	Distance       float64                `json:"distance"`
	TotalPrice     float64                `json:"total_price"`
	Shippings      []interface{}          `json:"shippings"`
	AvailableItems int                    `json:"available_items"`
	TotalItems     int                    `json:"total_items"`
	Items          []MerchantProductItems `json:"items"`
}

// swagger:model MerchantProductItems
type MerchantProductItems struct {
	SKU          string  `json:"sku"`
	MerchantSKU  string  `json:"merchant_sku"`
	Name         string  `json:"name"`
	QTY          int     `json:"qty"`
	QTYAvailable float64 `json:"qty_available"`
	UOM          string  `json:"uom"`
	UOMName      string  `json:"uom_name"`
	SellingPrice float64 `json:"selling_price"`
	SpecialPrice float64 `json:"special_price"`
	TotalPrice   float64 `json:"total_price"`
	Image        string  `json:"image"`
	IsAvailable  bool    `json:"is_available"`
	Status       string  `json:"status,omitempty"`
}

type ProductsAvailable struct {
	SKU          string  `json:"sku"`
	MerchantSKU  string  `json:"merchant_sku"`
	Name         string  `json:"name"`
	QTY          float64 `json:"qty"`
	UOM          string  `json:"uom"`
	UOMName      string  `json:"uom_name"`
	SellingPrice float64 `json:"selling_price"`
	SpecialPrice float64 `json:"special_price"`
	Image        string  `json:"image"`
}

type ProductsOrdered struct {
	SKU          string  `json:"sku"`
	MerchantSKU  string  `json:"merchant_sku"`
	Name         string  `json:"name"`
	QTY          int     `json:"qty"`
	QTYAvailable float64 `json:"qty_available"`
	UOM          string  `json:"uom"`
	UOMName      string  `json:"uom_name"`
	SellingPrice float64 `json:"selling_price"`
	SpecialPrice float64 `json:"special_price"`
	Available    float64 `json:"available"`
	IsAvailable  bool    `json:"is_available"` // out of stock (false) or not available (true)
	Status       string  `json:"status"`       // out of stock or not available
	Image        string  `json:"image"`
}

type Coordinates struct {
	Lat float64
	Lon float64
}
