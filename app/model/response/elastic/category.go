package responseelastic

type CategoryResponse struct {
	ID              int    `json:"id,omitempty"`
	ParentID        int    `json:"parent_id,omitempty"`
	StoreID         int    `json:"store_id,omitempty"`
	Name            string `json:"name,omitempty"`
	Slug            string `json:"slug,omitempty"`
	FullSlug        string `json:"full_slug,omitempty"`
	URLKey          string `json:"url_key,omitempty"`
	Icon            string `json:"icon,omitempty"`
	Image           string `json:"image,omitempty"`
	Description     string `json:"description,omitempty"`
	MetaTitle       string `json:"meta_title,omitempty"`
	MetaTitleH1     any    `json:"meta_title_h1,omitempty"`
	MetaDescription string `json:"meta_description,omitempty"`
	MetaKeyword     string `json:"meta_keyword,omitempty"`
	Position        int    `json:"position,omitempty"`
	Level           int    `json:"level,omitempty"`
	SortOrder       int    `json:"sort_order,omitempty"`
	ChildrenCount   int    `json:"children_count,omitempty"`
	LandingPage     string `json:"landing_page,omitempty"`
	InMenu          int    `json:"in_menu,omitempty"`
	InHome          int    `json:"in_home,omitempty"`
	InHomepage      int    `json:"in_homepage,omitempty"`
	Status          int    `json:"status,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
}

type CategoryTree struct {
	ID              int            `json:"id,omitempty"`
	Name            string         `json:"name,omitempty"`
	Slug            string         `json:"slug,omitempty"`
	URLKey          string         `json:"url_key,omitempty"`
	Level           int            `json:"level,omitempty"`
	ParentID        int            `json:"parent_id,omitempty"`
	Icon            string         `json:"icon,omitempty"`
	Image           string         `json:"image,omitempty"`
	MetaTitle       string         `json:"meta_title,omitempty"`
	MetaDescription string         `json:"meta_description,omitempty"`
	Sub             []CategoryTree `json:"sub,omitempty"`
}
