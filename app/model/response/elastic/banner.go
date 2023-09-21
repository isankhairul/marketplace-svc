package responseelastic

type BannerResponse struct {
	ID                  int    `json:"id,omitempty"`
	Title               string `json:"title,omitempty"`
	Slug                string `json:"slug,omitempty"`
	Image               string `json:"image,omitempty"`
	Brandlogo           string `json:"brandlogo,omitempty"`
	Description         string `json:"description,omitempty"`
	Link                string `json:"link,omitempty"`
	SlugProductCategory string `json:"slug_product_category,omitempty"`
	VoucherCode         string `json:"voucher_code,omitempty"`
	ChannelID           int    `json:"channel_id,omitempty"`
	CategorySlug        string `json:"category_slug,omitempty"`
	StoreCode           string `json:"store_code,omitempty"`
	Sort                any    `json:"sort,omitempty"`
	Status              int    `json:"status,omitempty"`
	CreatedAt           string `json:"created_at,omitempty"`
}
