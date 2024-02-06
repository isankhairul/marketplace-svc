package responsequote

import (
	"fmt"
	entitymerchant "marketplace-svc/app/model/entity/merchant"
	entity "marketplace-svc/app/model/entity/quote"
	"marketplace-svc/pkg/util"
)

type QuoteItemRs struct {
	ID                    uint64   `json:"id"`
	QuoteMerchantID       uint64   `json:"quote_merchant_id"`
	ProductSku            string   `json:"product_sku"`
	MerchantSku           string   `json:"merchant_sku"`
	Name                  string   `json:"name"`
	Slug                  string   `json:"slug"`
	Status                int      `json:"status"`
	IsActive              int      `json:"is_active"`
	MerchantProductStatus int      `json:"merchant_product_status"`
	Image                 string   `json:"image"`
	ItemNotes             string   `json:"item_notes"`
	Weight                float64  `json:"weight"`
	Quantity              int32    `json:"quantity"`
	Stock                 int      `json:"stock"`
	MaxQty                int      `json:"max_qty"`
	GlobalMaxQty          int      `json:"global_max_qty"`
	Price                 float64  `json:"price"`
	BrandName             string   `json:"brand_name"`
	BrandCode             string   `json:"brand_code"`
	DiscountPercentage    float64  `json:"discount_percentage"`
	DiscountAmount        float64  `json:"discount_amount"`
	RowWeight             float64  `json:"row_weight"`
	RowTotal              float64  `json:"row_total"`
	OriginalPrice         float64  `json:"original_price"`
	RowOriginalPrice      float64  `json:"row_original_price"`
	Redeem                int      `json:"redeem"`
	Selected              bool     `json:"selected"`
	ProductID             uint64   `json:"product_id"`
	Categories            []string `json:"categories"`
	ProductMerchantUID    string   `json:"product_merchant_uid"`
	StockStatus           int      `json:"stock_status"`
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
