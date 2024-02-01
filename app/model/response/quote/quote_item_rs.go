package responsequote

import (
	"fmt"
	entitymerchant "marketplace-svc/app/model/entity/merchant"
	entity "marketplace-svc/app/model/entity/quote"
	"marketplace-svc/pkg/util"
)

type QuoteItemRs struct {
	ID                    uint64   `json:"id,omitempty"`
	QuoteMerchantID       uint64   `json:"quote_merchant_id,omitempty"`
	ProductSku            string   `json:"product_sku,omitempty"`
	MerchantSku           string   `json:"merchant_sku,omitempty"`
	Name                  string   `json:"name,omitempty"`
	Slug                  string   `json:"slug,omitempty"`
	Status                int      `json:"status,omitempty"`
	IsActive              int      `json:"is_active,omitempty"`
	MerchantProductStatus int      `json:"merchant_product_status,omitempty"`
	Image                 string   `json:"image,omitempty"`
	ItemNotes             string   `json:"item_notes,omitempty"`
	Weight                float64  `json:"weight,omitempty"`
	Quantity              int32    `json:"quantity,omitempty"`
	Stock                 int      `json:"stock,omitempty"`
	MaxQty                int      `json:"max_qty,omitempty"`
	GlobalMaxQty          int      `json:"global_max_qty,omitempty"`
	Price                 float64  `json:"price,omitempty"`
	BrandName             string   `json:"brand_name,omitempty"`
	BrandCode             string   `json:"brand_code,omitempty"`
	DiscountPercentage    float64  `json:"discount_percentage,omitempty"`
	DiscountAmount        float64  `json:"discount_amount,omitempty"`
	RowWeight             float64  `json:"row_weight,omitempty"`
	RowTotal              float64  `json:"row_total,omitempty"`
	OriginalPrice         float64  `json:"original_price,omitempty"`
	RowOriginalPrice      float64  `json:"row_original_price,omitempty"`
	Redeem                int      `json:"redeem,omitempty"`
	Selected              bool     `json:"selected,omitempty"`
	MerchantIncludedItem  bool     `json:"merchant_included_item,omitempty"`
	IsKliknow             int      `json:"is_kliknow,omitempty"`
	ProductID             uint64   `json:"product_id,omitempty"`
	Categories            []string `json:"categories,omitempty"`
	ProductMerchantUID    string   `json:"product_merchant_uid,omitempty"`
	PointEarned           float64  `json:"point_earned,omitempty"`
	PointSpent            float64  `json:"point_spent,omitempty"`
	StockStatus           int      `json:"stock_status,omitempty"`
}

func (qr QuoteItemRs) Transform(oqi *entity.OrderQuoteItem, merchant entitymerchant.Merchant, mp entitymerchant.MerchantProduct, arrCategory []string, image string) *QuoteItemRs {
	var response QuoteItemRs
	if oqi == nil {
		return nil
	}

	intStockStatus := 0
	if oqi.ProductFlat.Status == 0 || oqi.ProductFlat.IsActive == 0 || mp.Status == 0 || mp.Stock < 1 {
		intStockStatus = 1
	}

	response = QuoteItemRs{
		ID:                    oqi.ID,
		QuoteMerchantID:       oqi.QuoteMerchantID,
		ProductSku:            oqi.ProductSku,
		MerchantSku:           oqi.MerchantSku,
		Name:                  oqi.Name,
		Slug:                  oqi.ProductFlat.Slug,
		Status:                oqi.ProductFlat.Status,
		IsActive:              oqi.ProductFlat.IsActive,
		Image:                 image,
		GlobalMaxQty:          util.StringToInt(oqi.ProductFlat.MaximumPurchaseQuantity),
		ItemNotes:             oqi.ItemNotes,
		Weight:                oqi.Weight,
		Quantity:              oqi.Quantity,
		MerchantProductStatus: mp.Status,
		Stock:                 mp.Stock,
		MaxQty:                mp.MaxPurchaseQty,
		Price:                 oqi.Price,
		BrandName:             oqi.BrandName,
		BrandCode:             oqi.BrandCode,
		DiscountPercentage:    oqi.DiscountPercentage,
		RowWeight:             oqi.RowWeight,
		RowTotal:              oqi.RowTotal,
		OriginalPrice:         oqi.OriginalPrice,
		RowOriginalPrice:      oqi.RowOriginalPrice,
		Selected:              oqi.Selected,
		ProductMerchantUID:    fmt.Sprint(oqi.ProductID) + "-" + merchant.MerchantUID,
		StockStatus:           intStockStatus,
		Categories:            arrCategory,
	}

	return &response
}
