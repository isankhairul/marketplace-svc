package entity

import "time"

type MerchantInfo struct {
	ID                    uint64    `json:"id"`
	MerchantID            uint64    `json:"merchant_id"`
	BankID                string    `json:"bank_id"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	TransferBankAccountID int       `json:"transfer_bank_account_id"`
	ApBankAccountID       int       `json:"ap_bank_account_id"`
	AccountNo             string    `json:"account_no"`
	AccountName           string    `json:"account_name"`
	BranchBankID          string    `json:"branch_bank_id"`
	JoinDate              time.Time `json:"join_date"`
	IsPaymentBuyingPrice  int       `json:"is_payment_buying_price"`
	EFulfillmentFee       float64   `json:"e_fulfillment_fee"`
	EDistributionFee      float64   `json:"e_distribution_fee"`
	MarfinRfpm            string    `json:"marfin_rfpm"`
	ManagementFee         float64   `json:"management_fee"`
	Npwp                  int16     `json:"npwp"`
	BusinessOwner         string    `json:"business_owner"`
	CustomerNumberEmos    string    `json:"customer_number_emos"`
	ShiptoNumberEmos      string    `json:"shipto_number_emos"`
}

func (m MerchantInfo) TableName() string {
	return "merchant_info"
}
