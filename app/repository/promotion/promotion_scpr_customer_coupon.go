package repository

import (
	entity "marketplace-svc/app/model/entity/promotion"
	base "marketplace-svc/app/repository"
)

type promotionScprCustomerCouponRepo struct {
	base base.BaseRepository
}

type PromotionScprCustomerCouponRepository interface {
	CheckClaimed(dbc *base.DBContext, ruleID int, customerID int, storeID int) bool
}

func NewPromotionScprCustomerCouponRepository(br base.BaseRepository) PromotionScprCustomerCouponRepository {
	return &promotionScprCustomerCouponRepo{br}
}

func (s promotionScprCustomerCouponRepo) CheckClaimed(dbc *base.DBContext, ruleID int, customerID int, storeID int) bool {
	var pscc entity.PromotionScprCustomerCoupon
	err := dbc.DB.Where("customer_id = ? and promotion_scpr_id = ? and store_id = ?", customerID, ruleID, storeID).
		First(&pscc).Error

	return err == nil
}
