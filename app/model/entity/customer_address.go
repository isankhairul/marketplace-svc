package entity

import "time"

type CustomerAddress struct {
	ID                int64     `json:"id"`
	Title             string    `json:"title"`
	CustomerID        int       `json:"customer_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Status            int       `json:"status"`
	CityID            int       `json:"city_id"`
	Company           any       `json:"company"`
	CountryID         int       `json:"country_id"`
	Fax               any       `json:"fax"`
	PostcodeID        int       `json:"postcode_id"`
	Coordinate        string    `json:"coordinate"`
	Prefix            any       `json:"prefix"`
	DistrictID        int       `json:"district_id"`
	ProvinceID        int       `json:"province_id"`
	Street            string    `json:"street"`
	Suffix            any       `json:"suffix"`
	Telephone         any       `json:"telephone"`
	VatID             any       `json:"vat_id"`
	VatIsValid        any       `json:"vat_is_valid"`
	VatRequestDate    any       `json:"vat_request_date"`
	VatRequestID      any       `json:"vat_request_id"`
	VatRequestSuccess any       `json:"vat_request_success"`
	SubdistrictID     int       `json:"subdistrict_id"`
	IsDefault         int       `json:"is_default"`
	ReceiverName      string    `json:"receiver_name"`
	Province          string    `json:"province"`
	City              string    `json:"city"`
	District          string    `json:"district"`
	Subdistrict       string    `json:"subdistrict"`
	CitrixProvinceID  string    `json:"citrix_province_id"`
	CitrixCityID      string    `json:"citrix_city_id"`
	CitrixDistrictID  string    `json:"citrix_district_id"`
	CitrixAddressID   string    `json:"citrix_address_id"`
	Zipcode           string    `json:"zipcode"`
	GoogleAddress     string    `json:"google_address"`
	ReceiverPhone     string    `json:"receiver_phone"`
	Notes             string    `json:"notes"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	IsCompleted       int       `json:"is_completed"`
}
