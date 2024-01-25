package responsequote

import (
	"encoding/json"
	"github.com/bytedance/sonic"
	"time"
)

type QuoteReceiptRs struct {
	UID   string `json:"uid"`
	Items []struct {
		Qty         int    `json:"qty"`
		Sku         string `json:"sku"`
		Uom         string `json:"uom"`
		Name        string `json:"name"`
		Note        string `json:"note"`
		Image       string `json:"image"`
		UomName     string `json:"uom_name"`
		PriceMax    int64  `json:"price_max"`
		PriceMin    int64  `json:"price_min"`
		IsOrdered   bool   `json:"is_ordered"`
		AturanPakai string `json:"aturan_pakai"`
	} `json:"items"`
	Status    string    `json:"status"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (oqr QuoteReceiptRs) Transform(dataReceipt *json.RawMessage) *QuoteReceiptRs {
	var response QuoteReceiptRs
	if dataReceipt == nil {
		return nil
	}
	err := sonic.Unmarshal(*dataReceipt, &response)
	if err != nil {
		return nil
	}

	return &response
}
