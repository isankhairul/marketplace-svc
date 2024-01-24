package responsequote

import (
	"context"
	"fmt"
	"marketplace-svc/app"
	entitymerchant "marketplace-svc/app/model/entity/merchant"
	entity "marketplace-svc/app/model/entity/quote"
	"marketplace-svc/app/repository"
	repomerchant "marketplace-svc/app/repository/merchant"
	"marketplace-svc/pkg/util"
)

type QuoteItemRs struct {
	ID                    int64    `json:"id,omitempty"`
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
	Quantity              int      `json:"quantity,omitempty"`
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
	ProductID             int      `json:"product_id,omitempty"`
	Categories            []string `json:"categories,omitempty"`
	ProductMerchantUID    string   `json:"product_merchant_uid,omitempty"`
	PointEarned           float64  `json:"point_earned,omitempty"`
	PointSpent            float64  `json:"point_spent,omitempty"`
	StockStatus           int      `json:"stock_status,omitempty"`
}

func (qr QuoteItemRs) Transform(oqis *[]entity.OrderQuoteItem, merchant entitymerchant.Merchant, infra app.Infra) *[]QuoteItemRs {
	var response []QuoteItemRs
	if oqis == nil {
		return nil
	}
	baseRepo := repository.NewBaseRepository(infra.DB)
	mpRepo := repomerchant.NewMerchantProductRepository(baseRepo)
	dbc := repository.DBContext{DB: baseRepo.GetDB(), Context: context.Background()}

	for _, oqi := range *oqis {
		filterMp := map[string]interface{}{
			"merchant_id":  merchant.ID,
			"product_sku":  oqi.ProductSku,
			"merchant_sku": oqi.MerchantSku,
		}
		mp, err := mpRepo.FindFirstByParams(&dbc, filterMp, true)
		if err != nil {
			continue
		}
		firstProductImage := *oqi.Product.ProductImage

		intStockStatus := 0
		if oqi.ProductFlat.Status == 0 || oqi.ProductFlat.IsActive == 0 || mp.Status == 0 || mp.Stock < 1 {
			intStockStatus = 1
		}

		response = append(response, QuoteItemRs{
			ID:                    oqi.ID,
			ProductSku:            oqi.ProductSku,
			MerchantSku:           oqi.MerchantSku,
			Name:                  oqi.Name,
			Slug:                  oqi.ProductFlat.Slug,
			Status:                oqi.ProductFlat.Status,
			IsActive:              oqi.ProductFlat.IsActive,
			Image:                 infra.Config.URL.BaseImageURL + firstProductImage[0].ImageThumbnail,
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
		})
	}

	return &response
}
