package responseelastic

type OrderResponse struct {
	ID                     int     `json:"id,omitempty"`
	OrderParentNo          int     `json:"order_parent_no,omitempty"`
	OrderNo                string  `json:"order_no,omitempty"`
	StatusID               int     `json:"status_id,omitempty"`
	StatusLabel            string  `json:"status_label,omitempty"`
	Store                  string  `json:"store,omitempty"`
	TotalWeight            float64 `json:"total_weight,omitempty"`
	TotalPrice             float64 `json:"total_price,omitempty"`
	DiscountAmount         float64 `json:"discount_amount,omitempty"`
	StoreID                int     `json:"store_id,omitempty"`
	DiscountPercentage     float64 `json:"discount_percentage,omitempty"`
	TotalQuantity          float64 `json:"total_quantity,omitempty"`
	ShippingDiscountAmount float64 `json:"shipping_discount_amount,omitempty"`
	TotalPointBonus        float64 `json:"total_point_bonus,omitempty"`
	TotalPointDiscount     float64 `json:"total_point_discount,omitempty"`
	Redeem                 bool    `json:"redeem,omitempty"`
	Event                  bool    `json:"event,omitempty"`
	AppID                  string  `json:"app_id,omitempty"`
	MarketplaceStoreID     int     `json:"marketplace_store_id,omitempty"`
	DeviceID               int     `json:"device_id,omitempty"`
	AgentID                any     `json:"agent_id,omitempty"`
	CouponCode             string  `json:"coupon_code,omitempty"`
	Rules                  []any   `json:"rules,omitempty"`
	Email                  string  `json:"email,omitempty"`
	PhoneNumber            string  `json:"phone_number,omitempty"`
	Address                string  `json:"address,omitempty"`
	Street                 string  `json:"street,omitempty"`
	ShippingAmount         float64 `json:"shipping_amount,omitempty"`
	GrandTotal             float64 `json:"grand_total,omitempty"`
	Subtotal               float64 `json:"subtotal,omitempty"`
	TotalPointEarned       float64 `json:"total_point_earned,omitempty"`
	TotalPointSpent        float64 `json:"total_point_spent,omitempty"`
	OrderDate              string  `json:"order_date,omitempty"`
	PaymentInfo            struct {
		PaymentDate   string `json:"payment_date,omitempty"`
		PaymentStatus bool   `json:"payment_status,omitempty"`
		PaymentMethod string `json:"payment_method,omitempty"`
	} `json:"payment_info,omitempty"`
	OrderItems []struct {
		ProductName  string `json:"product_name,omitempty"`
		MerchantSku  string `json:"merchant_sku,omitempty"`
		ProductImage string `json:"product_image,omitempty"`
		PointSpent   int    `json:"point_spent,omitempty"`
		PointEarned  int    `json:"point_earned,omitempty"`
		Quantity     int    `json:"quantity,omitempty"`
		Event        bool   `json:"event,omitempty"`
		Redeem       bool   `json:"redeem,omitempty"`
		FreeProduct  bool   `json:"free_product,omitempty"`
		Price        int    `json:"price,omitempty"`
		Rating       any    `json:"rating,omitempty"`
		ItemNotes    string `json:"item_notes,omitempty"`
	} `json:"order_items,omitempty"`
	OrderAddress struct {
		Title            string `json:"title,omitempty"`
		ReceiverName     string `json:"receiver_name,omitempty"`
		Street           string `json:"street,omitempty"`
		DistrictID       string `json:"district_id,omitempty"`
		DistrictName     string `json:"district_name,omitempty"`
		SubdistrictID    string `json:"subdistrict_id,omitempty"`
		SubdistrictName  string `json:"subdistrict_name,omitempty"`
		CityID           int    `json:"city_id,omitempty"`
		CityName         string `json:"city_name,omitempty"`
		ProvinceID       int    `json:"province_id,omitempty"`
		ProvinceName     string `json:"province_name,omitempty"`
		PostcodeID       string `json:"postcode_id,omitempty"`
		ZipCode          string `json:"zip_code,omitempty"`
		Address          string `json:"address,omitempty"`
		Coordinate       string `json:"coordinate,omitempty"`
		PhoneNumber      string `json:"phone_number,omitempty"`
		TrackingNumber   string `json:"tracking_number,omitempty"`
		ShippingMethod   string `json:"shipping_method,omitempty"`
		ShippingDuration string `json:"shipping_duration,omitempty"`
		DeliveryDate     any    `json:"delivery_date,omitempty"`
	} `json:"order_address,omitempty"`
	Customer struct {
		ID        int    `json:"id,omitempty"`
		Phone     string `json:"phone,omitempty"`
		Name      string `json:"name,omitempty"`
		ContactID string `json:"contact_id,omitempty"`
	} `json:"customer,omitempty"`
	Merchant struct {
		ID       int    `json:"id,omitempty"`
		Name     string `json:"name,omitempty"`
		Province string `json:"province,omitempty"`
		City     string `json:"city,omitempty"`
		Slug     string `json:"slug,omitempty"`
	} `json:"merchant,omitempty"`
	PushCitrix bool `json:"push_citrix,omitempty"`
	IsReviewed int  `json:"is_reviewed,omitempty"`
	DataSource struct {
		ID   any `json:"id,omitempty"`
		Name any `json:"name,omitempty"`
	} `json:"data_source,omitempty"`
}
