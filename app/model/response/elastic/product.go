package responseelastic

import (
	"strconv"
)

// swagger:model ProductResponse
type ProductResponse struct {
	ProdCode         string `json:"prod_code"`
	ProdName         string `json:"prod_name"`
	UOM              string `json:"uom"`
	UOMName          string `json:"uom_name"`
	Price            string `json:"price"`
	MinPrice         string `json:"min_price"`
	Weight           string `json:"weight"`
	Image            string `json:"image,omitempty"`
	Proportional     string `json:"proportional"`
	PharmacyCode     string `json:"pharmacy_code"`
	PrincipalName    string `json:"principal_name,omitempty"`
	Description      string `json:"description,omitempty"`
	ShortDescription string `json:"short_description,omitempty"`
}

func NewProductResponse(val map[string]interface{}) *ProductResponse {
	return &ProductResponse{
		ProdCode:         getESStringValue(val["sku"]),
		ProdName:         getESStringValue(val["name"]),
		UOM:              getESStringValue(val["uom"]),
		UOMName:          getESStringValue(val["uom_name"]),
		Price:            getESFloat64Value(val["price"]),
		MinPrice:         getESFloat64Value(val["min_price"]),
		Weight:           getESFloat64Value(val["weight"]),
		PrincipalName:    getESStringValue(val["principal_name"]),
		Description:      getESStringValue(val["description"]),
		ShortDescription: getESStringValue(val["short_description"]),
		Image:            getESImagesThumbnailValues(val["images"]),
		Proportional:     "0",
		PharmacyCode:     getESStringValue(val["pharmacy_code"]),
	}
}

func getESStringValue(value interface{}) string {
	if value == nil {
		return ""
	}
	strValue, _ := value.(string)
	return strValue
}

func getESFloat64Value(value interface{}) string {
	if value == nil {
		return ""
	}

	float64Value, _ := value.(float64)
	return strconv.FormatFloat(float64Value, 'f', -1, 64)
}

func getESImagesThumbnailValues(value interface{}) string {
	images, _ := value.([]interface{})
	var thumbnailURL string
	if len(images) > 0 {
		imageInfo, _ := images[0].(map[string]interface{})
		thumbnailURL, _ = imageInfo["thumbnail"].(string)
	}

	return thumbnailURL
}
