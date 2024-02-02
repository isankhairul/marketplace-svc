package response_kalcare

type meta struct {
	CorrelationID string      `json:"correlation_id"`
	Code          int         `json:"code"`
	Message       string      `json:"message"`
	Pagination    *pagination `json:"pagination"`
}

type pagination struct {
	Records      int `json:"records"`
	TotalRecords int `json:"total_records"`
	Limit        int `json:"limit"`
	Page         int `json:"page"`
	TotalPage    int `json:"total_page"`
}
