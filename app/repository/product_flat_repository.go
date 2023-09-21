package repository

import (
	"fmt"
	db "marketplace-svc/helper/database"
	"strings"
)

type productFlatRepository struct {
	BaseRepository
}

type ProductFlatRepository interface {
	GetQueryRawAllProduct(productIDs string, storeID int, typeProduct string, isCount bool, merchantID int) string
}

func NewProductFlatRepository(entityDb *db.Database) ProductFlatRepository {
	return &productFlatRepository{NewBaseRepository(entityDb)}
}

func (r *productFlatRepository) GetQueryRawAllProduct(productIDs string, storeID int, typeProduct string, isCount bool, merchantID int) string {
	// default value
	if storeID == 0 {
		storeID = 1
	}
	if typeProduct == "" {
		typeProduct = "sku"
	}
	var tbMerchantFlat string = "merchant_flat_" + fmt.Sprint(storeID)
	// end default value

	var where []string
	var sql string
	where = append(where, "pf.store_id="+fmt.Sprint(storeID))

	if productIDs != "" {
		if "sku" == strings.ToLower(typeProduct) {
			skus := "'" + strings.Join(strings.Split(productIDs, ","), "','") + "'"
			if merchantID != 0 {
				where = append(where, fmt.Sprintf("pf.sku in (%s) AND %s.id = %d", skus, tbMerchantFlat, merchantID))
			} else {
				where = append(where, fmt.Sprintf("pf.sku in (%s)", skus))
			}
		} else {
			ids := strings.Join(strings.Split(productIDs, ","), ",")
			if merchantID != 0 {
				where = append(where, fmt.Sprintf("pf.id in (%s) AND %s.id = %d", ids, tbMerchantFlat, merchantID))
			} else {
				where = append(where, fmt.Sprintf("pf.id in (%s)", ids))
			}
		}
	} else {
		if merchantID != 0 {
			where = append(where, fmt.Sprintf("%s.id = %d", tbMerchantFlat, merchantID))
		}
	}

	strWhere := strings.Join(where, " AND ")

	if isCount {
		sql = fmt.Sprintf(`SELECT COUNT(pf.id) AS total 
		FROM product_flat pf 
		INNER JOIN %s  ON %s.merchant_product_id = pf.id and pf.store_id = %d 
		WHERE %s`, tbMerchantFlat, tbMerchantFlat, storeID, strWhere)
	} else {
		sql = fmt.Sprintf(`SELECT  pf.id AS product_id, 
                            pf.name AS product_name, 
                            pf.sku AS product_sku, 
                            pf.slug AS product_slug, 
                            pf.barcode AS product_barcode, 
                            pf.brand_code AS product_brand_code,
                            pf.meta_title AS product_meta_title,
                            pf.meta_title_h1 AS product_meta_title_h1,
                            pf.meta_description AS product_meta_description,
                            pf.meta_keyword AS product_meta_keyword,
                            pf.description AS product_description,
                            pf.short_description AS product_short_description,
                            pf.images AS product_images,
                            pf.weight AS product_weight,
                            pf.base_point AS product_base_point,
                            pf.base_point_rupiah AS product_base_point_rupiah,
                            pf.reward_point_sell_product AS product_reward_point_sell_product,
                            pf.is_family_gift AS product_is_family_gift,
                            pf.is_free_product AS product_is_free_product,
                            pf.is_langganan AS product_is_langganan,
                            pf.is_spot AS product_is_spot,
                            pf.is_ticket AS product_is_ticket,
                            pf.is_kliknow AS product_is_kliknow,
                            pf.is_active AS product_is_active,
                            pf.type_id AS product_type_id,
                            pf.status AS product_master_status,
                            pf.created_at AS product_created_at,
                            pf.updated_at AS product_updated_at,
                            %s.*
			FROM product_flat pf
		INNER JOIN %s ON %s.merchant_product_id = pf.id and pf.store_id = %d
		WHERE %s`, tbMerchantFlat, tbMerchantFlat, tbMerchantFlat, storeID, strWhere)
	}

	return sql
}
