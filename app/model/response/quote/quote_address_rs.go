package responsequote

import entity "marketplace-svc/app/model/entity/quote"

type QuoteAddressRs struct {
	ID                  uint64 `json:"id,omitempty"`
	Title               string `json:"title,omitempty"`
	ReceiverName        string `json:"receiver_name,omitempty"`
	Street              string `json:"street,omitempty"`
	Telephone           string `json:"telephone,omitempty"`
	Coordinate          string `json:"coordinate,omitempty"`
	Province            string `json:"province,omitempty"`
	City                string `json:"city,omitempty"`
	District            string `json:"district,omitempty"`
	Subdistrict         string `json:"subdistrict,omitempty"`
	Postcode            string `json:"postcode,omitempty"`
	CustomerNotes       string `json:"customer_notes,omitempty"`
	DiscountDescription string `json:"discount_description,omitempty"`
}

func (qr QuoteAddressRs) Transform(qa *entity.OrderQuoteAddress) *[]QuoteAddressRs {
	var response []QuoteAddressRs //nolint:prealloc
	if qa == nil {
		return nil
	}

	// set response
	response = append(response, QuoteAddressRs{
		Title:         qa.Title,
		Province:      qa.Province,
		City:          qa.City,
		District:      qa.District,
		Subdistrict:   qa.Subdistrict,
		Street:        qa.Street,
		Telephone:     qa.PhoneNumber,
		Coordinate:    qa.Coordinate,
		CustomerNotes: qa.CustomerNotes,
		Postcode:      qa.Zipcode,
		ReceiverName:  qa.ReceiverName,
	})

	return &response
}
