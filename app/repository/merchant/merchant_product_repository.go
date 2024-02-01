package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	entity "marketplace-svc/app/model/entity/merchant"
	base "marketplace-svc/app/repository"
	"strings"
)

type merchantProductRepository struct {
	base.BaseRepository
}

type MerchantProductRepository interface {
	FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.MerchantProduct, error)
	FindFirstDetailMerchantProduct(dbc *base.DBContext, filter map[string]interface{}) (*entity.DetailMerchantProduct, error)
}

func NewMerchantProductRepository(br base.BaseRepository) MerchantProductRepository {
	return &merchantProductRepository{br}
}

func (r *merchantProductRepository) FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.MerchantProduct, error) {
	var merchantProduct entity.MerchantProduct
	query := dbc.DB.WithContext(dbc.Context).Table(merchantProduct.TableName())

	for key, v := range filter {
		if key == "merchant_id" && v != "" {
			query = query.Where("merchant_id = ?", v.(uint64))
		}
		if key == "product_sku" && v != "" {
			query = query.Where("product_sku = ?", v.(string))
		}
	}

	if isPreload {

	}

	err := query.Omit("reserved_stock,stock_on_hand,buffer_stockparent_reserved_stock,merchant_included_item,old_status,updated_by,created_at,updated_at,deleted_at").
		Preload("MerchantProductPrice", func(db *gorm.DB) *gorm.DB {
			return db.Select("id,merchant_product_id,selling_price,special_price,special_price_start_time,special_price_end_time,store_id,merchant_id")
		}).
		First(&merchantProduct).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &merchantProduct, nil
}

func (r *merchantProductRepository) FindFirstDetailMerchantProduct(dbc *base.DBContext, filter map[string]interface{}) (*entity.DetailMerchantProduct, error) {
	var response entity.DetailMerchantProduct
	query := dbc.DB.WithContext(dbc.Context)
	var arrWhere []string

	for key, v := range filter {
		if key == "merchant_id" && v != "" {
			arrWhere = append(arrWhere, fmt.Sprintf(" mp.merchant_id = %s ", fmt.Sprint(v.(uint64))))
		}
		if key == "merchant_sku" && v != "" {
			arrWhere = append(arrWhere, fmt.Sprintf(" mp.merchant_sku = '%s' ", fmt.Sprint(v.(string))))
		}
		if key == "store_id" && v != "" {
			arrWhere = append(arrWhere, fmt.Sprintf(" p.store_id = %s ", fmt.Sprint(v.(int))))
			arrWhere = append(arrWhere, fmt.Sprintf(" mpp.store_id = %s ", fmt.Sprint(v.(int))))
		}
	}
	querySql := `
			SELECT DISTINCT (p.id) AS id, p.sku AS sku, mp.merchant_sku AS merchant_sku,
			mp.merchant_included_item AS merchant_included_item, p.name AS name, p.slug AS slug, 
			p.weight AS weight, p.brand_code AS brand_code, p.product_kn AS product_kn, p.product_kalbe AS product_kalbe,
			p.base_point AS base_point, p.base_price AS base_price, p.type_id AS product_type, 
			p.attribute_set_id AS attribute_set_id, 0 AS reward_point_sell_product, 
			mpp.special_price_start_time AS special_price_start_time, mpp.special_price_end_time AS special_price_end_time, mpp.selling_price AS selling_price,
			mpp.special_price AS special_price 
			FROM product_flat AS p  
			INNER JOIN merchant_product AS mp ON mp.product_id = p.id 
			INNER JOIN merchant AS m ON m.id = mp.merchant_id 
			INNER JOIN merchant_product_price AS mpp ON mpp.merchant_product_id = mp.id AND mpp.merchant_id = mp.merchant_id 
			WHERE p.status = 1 AND p.is_family_gift = 0 AND m.status = 1 AND mp.status = 1
			`
	if len(arrWhere) > 0 {
		querySql += " AND " + strings.Join(arrWhere, " AND ")
	}
	err := query.Raw(querySql).First(&response).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &response, nil
}
