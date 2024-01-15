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
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
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
	ID                    int64     `json:"id"`
	QuoteID               int64     `json:"quote_id"`
	ContactID             string    `json:"contact_id"`
	SaveInAddressBook     bool      `json:"save_in_address_book"`
	CustomerAddressID     int       `json:"customer_address_id"`
	AddressTypeID         int       `json:"address_type_id"`
	Email                 string    `json:"email"`
	ReceiverName          string    `json:"receiver_name"`
	Street                string    `json:"street"`
	CountryID             int       `json:"country_id"`
	ProvinceID            int       `json:"province_id"`
	CityID                int       `json:"city_id"`
	DistrictID            int       `json:"district_id"`
	SubdistrictID         int       `json:"subdistrict_id"`
	Zipcode               string    `json:"zipcode"`
	PhoneNumber           string    `json:"phone_number"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	CustomerNotes         string    `json:"customer_notes"`
	Title                 string    `json:"title"`
	PostcodeID            int       `json:"postcode_id"`
	DiscountDescription   string    `json:"discount_description"`
	BonusPointDescription string    `json:"bonus_point_description"`
	Latitude              float64   `json:"latitude"`
	Longitude             float64   `json:"longitude"`
	Coordinate            Point     `json:"coordinate"`
	Province              string    `json:"province"`
	City                  string    `json:"city"`
	District              string    `json:"district"`
	Subdistrict           string    `json:"subdistrict"`
}

func (o OrderQuoteAddress) TableName() string {
	return "order_quote_address"
}
