package responseelastic

type BrandResponse struct {
	ID              int    `json:"id,omitempty"`
	Code            string `json:"code,omitempty"`
	Name            string `json:"name,omitempty"`
	Slug            string `json:"slug,omitempty"`
	PrincipalCode   string `json:"principal_code,omitempty"`
	Image           string `json:"image,omitempty"`
	MetaTitle       string `json:"meta_title,omitempty"`
	MetaDescription string `json:"meta_description,omitempty"`
	SortOrder       int    `json:"sort_order,omitempty"`
	Status          int    `json:"status,omitempty"`
	ShowOfficial    int    `json:"show_official,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
}
