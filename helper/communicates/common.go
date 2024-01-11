package communicates

type KalstoreCustomerInfo struct {
	ID         int                          `json:"id"`
	Email      string                       `json:"email"`
	Name       string                       `json:"name"`
	Phone      string                       `json:"phone"`
	DOB        string                       `json:"dob"`
	Gender     int                          `json:"gender"`
	Group      string                       `json:"group"`
	ContactID  string                       `json:"contact_id"`
	MemberID   string                       `json:"member_id"`
	StoreID    int                          `json:"store_id"`
	StoreCode  string                       `json:"store_code"`
	QuoteCode  string                       `json:"quote_code"`
	QuoteTotal int                          `json:"quote_total"`
	Address    *KalstoreCustomerInfoAddress `json:"address"`
}

type KalstoreCustomerInfoAddress struct {
	Street        string  `json:"street"`
	Telephone     string  `json:"telephone"`
	ProvinceID    int     `json:"province_id"`
	Province      string  `json:"province"`
	CityID        int     `json:"city_id"`
	City          string  `json:"city"`
	DistrictID    int     `json:"district_id"`
	District      string  `json:"district"`
	SubDistrict   string  `json:"subdistrict"`
	ZIPCode       string  `json:"zipcode"`
	Latitude      *string `json:"latitude"`
	Longitude     *string `json:"longitude"`
	GoogleAddress string  `json:"google_address"`
}

type KalstoreCustomerInfoDetail struct {
	Meta struct {
		CorrelationID string `json:"correlation_id"`
		Code          int    `json:"code"`
		Message       string `json:"message"`
		Pagination    struct {
			Records      int `json:"records"`
			TotalRecords int `json:"total_records"`
			Limit        int `json:"limit"`
			Page         int `json:"page"`
			Totalpage    int `json:"total_page"`
		} `json:"pagination"`
	} `json:"meta"`
	Data struct {
		Record KalstoreCustomerInfo `json:"record"`
	} `json:"data"`
	Message string `json:"message"`
}
