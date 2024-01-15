package modelelastic

type EsMerchantFlat struct {
	ID             string  `json:"id"`
	MerchantID     int     `json:"merchant_id"`
	MerchantUID    string  `json:"merchant_uid"`
	ProductID      int     `json:"product_id"`
	ProductSKU     string  `json:"product_sku"`
	MerchantSKU    string  `json:"merchant_sku"`
	ProductStatus  int     `json:"product_status"`
	Stock          int     `json:"stock"`
	StockOnHand    int     `json:"stock_on_hand"`
	Rating         float64 `json:"rating"`
	Review         int     `json:"review"`
	Categories     []int   `json:"categories"`
	MaxPurchaseQty int     `json:"max_purchase_qty"`
	SellingPrice   int     `json:"selling_price"`
	SpecialPrices  []struct {
		Price           int    `json:"price"`
		ToTime          string `json:"to_time"`
		FromTime        string `json:"from_time"`
		CustomerGroupID int    `json:"customer_group_id"`
	} `json:"special_prices"`
	Status   int `json:"status"`
	TypeID   int `json:"type_id"`
	Location struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	} `json:"location"`
	UpdatedAt string `json:"updated_at"`
}
