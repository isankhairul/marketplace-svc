package entity

import "time"

type Province struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	CountryID   uint64    `json:"country_id"`
	Description string    `json:"description"`
	Code        string    `json:"code"`
	UpdatedAt   time.Time `json:"updated_at"`
}
