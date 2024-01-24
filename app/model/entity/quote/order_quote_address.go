package entity

import (
	"bytes"
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Point struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

func (p Point) Value() (driver.Value, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "(%f %f)", p.Latitude, p.Longitude)
	return buf.Bytes(), nil
}

func (p Point) String() string {
	return fmt.Sprintf("(%v %v)", p.Latitude, p.Longitude)
}

func (p Point) Scan(val interface{}) (err error) {
	if bb, ok := val.([]uint8); ok {
		tmp := bb[1 : len(bb)-1]
		coors := strings.Split(string(tmp[:]), ",")
		if p.Latitude, err = strconv.ParseFloat(coors[0], 64); err != nil {
			return err
		}
		if p.Longitude, err = strconv.ParseFloat(coors[1], 64); err != nil {
			return err
		}
	}
	return nil
}

type OrderQuoteAddress struct {
	ID                    int64     `json:"id,omitempty"`
	QuoteID               int64     `json:"quote_id,omitempty"`
	ContactID             string    `json:"contact_id,omitempty"`
	SaveInAddressBook     bool      `json:"save_in_address_book,omitempty"`
	CustomerAddressID     int       `json:"customer_address_id,omitempty"`
	AddressTypeID         int       `json:"address_type_id,omitempty"`
	Email                 string    `json:"email,omitempty"`
	ReceiverName          string    `json:"receiver_name,omitempty"`
	Street                string    `json:"street,omitempty"`
	CountryID             int       `json:"country_id,omitempty"`
	ProvinceID            int       `json:"province_id,omitempty"`
	CityID                int       `json:"city_id,omitempty"`
	DistrictID            int       `json:"district_id,omitempty"`
	SubdistrictID         int       `json:"subdistrict_id,omitempty"`
	Zipcode               string    `json:"zipcode,omitempty"`
	PhoneNumber           string    `json:"phone_number,omitempty"`
	CreatedAt             time.Time `json:"created_at,omitempty"`
	UpdatedAt             time.Time `json:"updated_at,omitempty"`
	CustomerNotes         string    `json:"customer_notes,omitempty"`
	Title                 string    `json:"title,omitempty"`
	PostcodeID            int       `json:"postcode_id,omitempty"`
	DiscountDescription   string    `json:"discount_description,omitempty"`
	BonusPointDescription string    `json:"bonus_point_description,omitempty"`
	Latitude              float64   `json:"latitude,omitempty"`
	Longitude             float64   `json:"longitude,omitempty"`
	Coordinate            string    `json:"coordinate,omitempty"`
	Province              string    `json:"province,omitempty"`
	City                  string    `json:"city,omitempty"`
	District              string    `json:"district,omitempty"`
	Subdistrict           string    `json:"subdistrict,omitempty"`
}

func (o OrderQuoteAddress) TableName() string {
	return "order_quote_address"
}
