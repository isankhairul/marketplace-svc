package entity

import "time"

type City struct {
	ID         uint64    `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	ProvinceID uint64    `json:"province_id"`
	UpdatedAt  time.Time `json:"updated_at"`
}
