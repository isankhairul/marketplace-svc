package entity

import "time"

type Banner struct {
	ID                  int       `json:"id"`
	BannerCategoryID    int       `json:"banner_category_id"`
	ChannelID           int       `json:"channel_id"`
	Title               string    `json:"title"`
	Slug                string    `json:"slug"`
	Description         any       `json:"description"`
	Link                string    `json:"link"`
	Brandlogo           any       `json:"brandlogo"`
	Image               string    `json:"image"`
	VoucherCode         any       `json:"voucher_code"`
	SlugProductCategory any       `json:"slug_product_category"`
	Status              int       `json:"status"`
	MappingBannerID     any       `json:"mapping_banner_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	Sort                any       `json:"sort"`
}
