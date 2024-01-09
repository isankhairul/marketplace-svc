package requestelastic

import (
	"fmt"
	"marketplace-svc/app/model/base"
	"marketplace-svc/helper/global"

	validation "github.com/itgelo/ozzo-validation"
)

type MerchantRequest struct {
	Query   string `json:"q" schema:"q" binding:"omitempty"`
	Fields  string `json:"fields" schema:"fields" binding:"omitempty"`
	PID     string `json:"pid" schema:"pid"`
	CID     string `json:"cid" schema:"cid"`
	Rating  string `json:"rating" schema:"rating"`
	Type    string `json:"type" schema:"type"`
	Zipcode string `json:"zipcode" schema:"zipcode"`
	StoreID int    `json:"store_id" schema:"store_id" binding:"omitempty"`
	Store   string `json:"store" schema:"store" binding:"omitempty"`
	Sort    string `json:"sort" schema:"sort" binding:"omitempty"`
	Dir     string `json:"dir" schema:"dir" binding:"omitempty"`
	Page    int    `json:"page" schema:"page" binding:"omitempty"`
	Limit   int    `json:"limit" schema:"limit" binding:"omitempty"`
}

func (req MerchantRequest) ToString() string {
	return fmt.Sprintf("%s-%s-%d-%d-%s-%s-%s-%d-%d-%d", //nolint:govet
		req.Query, req.Fields, req.PID, req.CID, req.Rating, req.Type, req.Zipcode, req.StoreID, req.Page, req.Limit)
}

func (req MerchantRequest) DefaultPagination() MerchantRequest {
	if req.Limit == 0 {
		req.Limit = base.PAGINATION_MIN_LIMIT
	}
	if req.Limit > base.PAGINATION_MAX_LIMIT {
		req.Limit = base.PAGINATION_MAX_LIMIT
	}

	// Default page 1
	if req.Page == 0 {
		req.Page = 1
	}

	// default storeID
	if req.StoreID == 0 {
		req.StoreID = 1
	}

	return req
}

type MerchantDetailRequest struct {
	ID      int    `json:"id" schema:"id" binding:"omitempty"`
	Slug    string `json:"slug" schema:"slug" binding:"omitempty"`
	Fields  string `json:"fields" schema:"fields" binding:"omitempty"`
	StoreID int    `json:"store_id" schema:"store_id" binding:"omitempty"`
}

func (req MerchantDetailRequest) ToString() string {
	return fmt.Sprintf("%s-%d", req.Fields, req.StoreID)
}

type MerchantZipcodeRequest struct {
	Zipcode string `json:"zipcode" schema:"zipcode" binding:"omitempty"`
	Fields  string `json:"fields" schema:"fields" binding:"omitempty"`
	Page    int    `json:"page" schema:"page" binding:"omitempty"`
	Limit   int    `json:"limit" schema:"limit" binding:"omitempty"`
}

func (req MerchantZipcodeRequest) ToString() string {
	return fmt.Sprintf("%s-%s-%d-%d", req.Zipcode, req.Fields, req.Page, req.Limit)
}

func (req MerchantZipcodeRequest) DefaultPagination() MerchantZipcodeRequest {
	if req.Limit == 0 {
		req.Limit = base.PAGINATION_MIN_LIMIT
	}
	if req.Limit > base.PAGINATION_MAX_LIMIT {
		req.Limit = base.PAGINATION_MAX_LIMIT
	}

	// Default page 1
	if req.Page == 0 {
		req.Page = 1
	}

	return req
}

// swagger:parameters MerchantProductRequest
type MerchantProductRequest struct {
	// Additional Fields
	// Example: "description,short_description"
	Fields string `json:"fields" schema:"fields" binding:"omitempty"`
	// StoreID
	StoreID *int `json:"store_id" schema:"store_id" binding:"omitempty"`
	// Page number
	Page int `json:"page" schema:"page" binding:"omitempty"`
	// Maximum records per page
	Limit int `json:"limit" schema:"limit" binding:"omitempty"`
	// Field to be sorted
	Sort       string              `schema:"sort"`
	Body       BodyMerchantProduct `json:"body"`
	JwtPayload global.JWTPayload   `json:"-"`
	Token      string              `json:"-"`
}

type BodyMerchantProduct struct {
	Lat   float64             `json:"lat"`
	Lon   float64             `json:"lon"`
	Items []BodyReceiptsItems `json:"items"`
}

type BodyReceiptsItems struct {
	SKU string `json:"sku"`
	QTY int    `json:"qty"`
}

func (b MerchantProductRequest) ToString() string {
	return fmt.Sprintf("%s-%d-%d-%d", //nolint:govet
		b.Fields, b.StoreID, b.Page, b.Limit)
}

func (req MerchantProductRequest) DefaultPagination() MerchantProductRequest {
	if req.Limit == 0 {
		req.Limit = base.PAGINATION_MIN_LIMIT
	}
	if req.Limit > base.PAGINATION_MAX_LIMIT {
		req.Limit = base.PAGINATION_MAX_LIMIT
	}

	// Default page 1
	if req.Page == 0 {
		req.Page = 1
	}
	return req
}

func (req MerchantProductRequest) Validate() error {
	return validation.ValidateStruct(&req.Body,
		validation.Field(&req.Body.Lat, validation.Required),
		validation.Field(&req.Body.Lon, validation.Required),
	)
}
