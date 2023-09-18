package base

// swagger:model PaginationResponse
type Pagination struct {
	Records      int64 `json:"records"`
	TotalRecords int64 `json:"total_records"`
	Limit        int   `json:"limit"`
	Page         int   `json:"page"`
	TotalPage    int   `json:"total_page"`
}

type Paging struct {
	// Maximum records per page
	// in: query
	Limit int `json:"limit" schema:"limit" binding:"omitempty,numeric,min=1,max=100"`

	// Page No
	// in: query
	Page int `json:"page" schema:"page" binding:"omitempty,numeric,min=1"`

	// Sort fields
	// in: query
	Sort string `json:"sort" schema:"sort" binding:"omitempty"`

	Dir string `json:"dir" schema:"dir" binding:"omitempty"`

	//Global Search
	// in: query
	Query string `json:"q" schema:"q" binding:"omitempty"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}
