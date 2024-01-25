package entity

import "time"

type PaymentMethod struct {
	ID                  uint64            `json:"id"`
	Code                string            `json:"code"`
	PaymentMethodTypeID uint64            `json:"payment_method_type_id"`
	Name                string            `json:"name"`
	Description         string            `json:"description"`
	Mdr                 int               `json:"mdr"`
	MdrTypeID           int               `json:"mdr_type_id"`
	BankAccountID       int               `json:"bank_account_id"`
	Status              int               `json:"status"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
	Logo                string            `json:"logo"`
	MinimumAmount       float64           `json:"minimum_amount"`
	SortOrder           int               `json:"sort_order"`
	MarfinDeduct        int               `json:"marfin_deduct"`
	Repay               int               `json:"repay"`
	UseSettlementData   int               `json:"use_settlement_data"`
	Extra               string            `json:"extra"`
	MdrMobile           int               `json:"mdr_mobile"`
	MdrTypeIDMobile     int               `json:"mdr_type_id_mobile"`
	AdminFee            float64           `json:"admin_fee"`
	AdminFeeType        int               `json:"admin_fee_type"`
	PaymentSvcUUID      string            `json:"payment_svc_uuid"`
	PaymentMethodType   PaymentMethodType `gorm:"foreignKey:payment_method_type_id;references:id" json:"-"`
}

func (m PaymentMethod) TableName() string {
	return "payment_method"
}
