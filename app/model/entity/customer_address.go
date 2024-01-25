package entity

import "time"

type CustomerAddress struct {
	ID            uint64    `json:"id"`
	Title         string    `json:"title"`
	CustomerID    uint64    `json:"customer_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Status        int       `json:"status"`
	CityID        int       `json:"city_id"`
	Company       string    `json:"company"`
	CountryID     int       `json:"country_id"`
	Fax           string    `json:"fax"`
	PostcodeID    int       `json:"postcode_id"`
	Coordinate    string    `json:"coordinate"`
	Prefix        string    `json:"prefix"`
	DistrictID    int       `json:"district_id"`
	ProvinceID    int       `json:"province_id"`
	Street        string    `json:"street"`
	Suffix        string    `json:"suffix"`
	Telephone     string    `json:"telephone"`
	SubdistrictID int       `json:"subdistrict_id"`
	IsDefault     int       `json:"is_default"`
	ReceiverName  string    `json:"receiver_name"`
	Province      string    `json:"province"`
	City          string    `json:"city"`
	District      string    `json:"district"`
	Subdistrict   string    `json:"subdistrict"`
	Zipcode       string    `json:"zipcode"`
	GoogleAddress string    `json:"google_address"`
	ReceiverPhone string    `json:"receiver_phone"`
	Notes         string    `json:"notes"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	IsCompleted   int       `json:"is_completed"`
}

func (m CustomerAddress) TableName() string {
	return "customer_address"
}
