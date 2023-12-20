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
	MaxPrice         string `json:"max_price"`
	Weight           string `json:"weight"`
	PrincipalName    string `json:"principal_name"`
	Description      string `json:"description"`
	ShortDescription string `json:"short_description"`
	Image            string `json:"image"`
	Proportional     string `json:"proportional"`
	PharmacyCode     string `json:"pharmacy_code"`
}

func NewProductResponse(val map[string]interface{}) *ProductResponse {
	return &ProductResponse{
		ProdCode:         getESStringValue(val["sku"]),
		ProdName:         getESStringValue(val["name"]),
		UOM:              getESStringValue(val["uom"]),
		UOMName:          getESStringValue(val["uom_name"]),
		Price:            getESStringValue(val["price"]),
		MinPrice:         getESStringValue(val["min_price"]),
		MaxPrice:         getESStringValue(val["max_price"]),
		Weight:           getESFloat64Value(val["weight"]),
		PrincipalName:    getESStringValue(val["principal_name"]),
		Description:      getESStringValue(val["description"]),
		ShortDescription: getESStringValue(val["short_description"]),
		Image:            getESImagesThumbnailValues(val["images"]),
		Proportional:     getESStringValue(val["proportional"]),
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
