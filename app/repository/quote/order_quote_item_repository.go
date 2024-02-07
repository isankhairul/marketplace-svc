package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	modelbase "marketplace-svc/app/model/base"
	entitycatalog "marketplace-svc/app/model/entity/catalog"
	entitymerchant "marketplace-svc/app/model/entity/merchant"
	entity "marketplace-svc/app/model/entity/quote"
	base "marketplace-svc/app/repository"
	"marketplace-svc/pkg/util"
	"strings"
)

type orderQuoteItemRepository struct {
	base.BaseRepository
}

type OrderQuoteItemRepository interface {
	Create(dbc *base.DBContext, oqi *entity.OrderQuoteItem) (*entity.OrderQuoteItem, error)
	Save(dbc *base.DBContext, oqi *entity.OrderQuoteItem) (*entity.OrderQuoteItem, error)
	FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.OrderQuoteItem, error)
	FindByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool, limit int, page int) (*[]entity.OrderQuoteItem, *modelbase.Pagination, error)
	UpdateByID(dbc *base.DBContext, id uint64, data entity.OrderQuoteItem) error
	FindRawByParams(dbc *base.DBContext, filter map[string]interface{}) (*[]entity.OrderQuoteItem, error)
	DeleteByID(dbc *base.DBContext, id uint64) error
}

func NewOrderQuoteItemRepository(br base.BaseRepository) OrderQuoteItemRepository {
	return &orderQuoteItemRepository{br}
}

func (r *orderQuoteItemRepository) Create(dbc *base.DBContext, oqi *entity.OrderQuoteItem) (*entity.OrderQuoteItem, error) {
	err := dbc.DB.WithContext(dbc.Context).Create(oqi).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return oqi, nil
}

func (r *orderQuoteItemRepository) Save(dbc *base.DBContext, oqi *entity.OrderQuoteItem) (*entity.OrderQuoteItem, error) {
	err := dbc.DB.WithContext(dbc.Context).Save(oqi).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return oqi, nil
}

func (r *orderQuoteItemRepository) FindFirstByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool) (*entity.OrderQuoteItem, error) {
	var orderQuote entity.OrderQuoteItem
	query := dbc.DB.WithContext(dbc.Context).Table(orderQuote.TableName())

	for key, v := range filter {
		if key == "quote_id" && v != "" {
			query = query.Where("quote_id = ?", v.(uint64))
		}
		if key == "quote_merchant_id" && v != "" {
			query = query.Where("quote_merchant_id = ?", v.(uint64))
		}
		if key == "arr_quote_merchant_id" && v != "" {
			query = query.Where("quote_merchant_id in ?", v.([]string))
		}
		if key == "merchant_sku" && v != "" {
			query = query.Where("merchant_sku = ?", v.(string))
		}
	}
	if isPreload {
		query = query.
			Preload("ProductFlat", func(db *gorm.DB) *gorm.DB {
				return db.Select("sku", "is_active", "status", "slug", "name")
			}).
			Preload("Product", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "sku", "slug", "name")
			}).
			Preload("Product.ProductImage", func(db *gorm.DB) *gorm.DB {
				return db.Select("product_id", "image_thumbnail", "image").Where("is_default=1 and status=true")
			})
	}

	err := query.Omit("point_earned,point_spent,point_spent_conversion,created_at,updated_at,row_point_spent,redeem,product_kn,product_kalbe,event,event_online,free_product,free_product_commit,free_product_rule_id,non_changeable_item").
		Order("id DESC").
		Find(&orderQuote).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	return &orderQuote, nil
}

func (r *orderQuoteItemRepository) FindByParams(dbc *base.DBContext, filter map[string]interface{}, isPreload bool, limit int, page int) (*[]entity.OrderQuoteItem, *modelbase.Pagination, error) {
	var orderQuotes []entity.OrderQuoteItem
	var pagination modelbase.Pagination

	query := dbc.DB.WithContext(dbc.Context)
	pagination.Limit = limit
	pagination.Page = page

	for key, v := range filter {
		if key == "quote_id" && v != "" {
			query = query.Where("quote_id = ?", v.(uint64))
		}
		if key == "quote_merchant_id" && v != "" {
			query = query.Where("quote_merchant_id = ?", v.(uint64))
		}
		if key == "arr_quote_merchant_id" && v != "" {
			query = query.Where("quote_merchant_id in ?", v.([]uint64))
		}
	}
	if isPreload {
		query = query.
			Preload("ProductFlat", func(db *gorm.DB) *gorm.DB {
				return db.Select("sku", "is_active", "status", "slug", "name")
			}).
			Preload("Product", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "sku", "slug", "name")
			}).
			Preload("Product.ProductImage", func(db *gorm.DB) *gorm.DB {
				return db.Select("product_id", "image_thumbnail", "image").Where("is_default=1 and status=true")
			})
	}

	err := query.Omit("point_earned,point_spent,point_spent_conversion,created_at,updated_at,row_point_spent,redeem,product_kn,product_kalbe,event,event_online,free_product,free_product_commit,free_product_rule_id,non_changeable_item").
		Scopes(r.Paginate(orderQuotes, &pagination, query)).
		Order("id DESC").
		Find(&orderQuotes).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	pagination.Records = int64(len(orderQuotes))

	return &orderQuotes, &pagination, nil
}

func (r *orderQuoteItemRepository) FindRawByParams(dbc *base.DBContext, filter map[string]interface{}) (*[]entity.OrderQuoteItem, error) {
	var orderQuoteItems []entity.OrderQuoteItem
	var orderQuoteItemsJoin []entity.OrderQuoteItemJoin
	var arrWhere []string
	for key, v := range filter {
		if key == "arr_quote_merchant_id" && v != "" {
			arrWhere = append(arrWhere, fmt.Sprintf(" quote_merchant_id in (%s) ", util.IntToString(v.([]uint64))))
		}
		if key == "quote_merchant_id" && v != "" {
			arrWhere = append(arrWhere, fmt.Sprintf(" quote_merchant_id = %s ", fmt.Sprint(v.(uint64))))
		}
		if key == "merchant_id" && v != "" {
			arrWhere = append(arrWhere, fmt.Sprintf(" merchant_id = %s ", fmt.Sprint(v.(uint64))))
		}
		if key == "merchant_sku" && v != "" {
			arrWhere = append(arrWhere, fmt.Sprintf(" merchant_sku = '%s' ", fmt.Sprint(v.(string))))
		}
	}

	query := dbc.DB.WithContext(dbc.Context)
	err := query.Raw(
		fmt.Sprintf(
			`select oqi.id,oqi.quote_merchant_id,oqi.product_id,oqi.item_type_id,oqi.product_sku,oqi.merchant_sku,oqi.merchant_category_id,oqi.category_id,oqi.brand_id,oqi.name,oqi.item_notes,oqi.weight,oqi.quantity,oqi.price,oqi.discount_percentage,oqi.discount_amount,oqi.row_weight,oqi.row_total,oqi.original_price,oqi.row_original_price,oqi.base_price,oqi.base_discount_amount,oqi.base_row_total,oqi.promo_description,oqi.bonus_point,oqi.discount_point,oqi.discount_weight,oqi.min_price,oqi.max_price,oqi.brand_name,oqi.brand_code,oqi.selected,oqi.location,oqi.additional_info,oqi.parent_info,oqi.start_date,oqi.end_date,oqi.product_type,oqi.quote_id,oqi.attribute_set_id,oqi.merchant_special_price, oqi.selected,
			 pf.is_active as pf_is_active, pf.status as pf_status,
			p.id as p_id, p.sku as p_sku, p.name as p_name, p.slug as p_slug, p.status as p_status,
			oqm.merchant_id as oqm_merchant_id, oqm.selected as oqm_selected,
			m.id as m_id, m.merchant_uid as m_merchant_uid, m.merchant_name as m_merchant_name, m.slug as m_slug, m.merchant_code as m_merchant_code
			from order_quote_item oqi
			inner join order_quote_merchant oqm on oqi.quote_merchant_id = oqm.id
			inner join merchant m on oqm.merchant_id = m.id
			inner join product_flat pf on oqi.product_sku = pf.sku and store_id=1
			inner join product p on oqi.product_id = p.id
			where %s
			limit 50`, strings.Join(arrWhere, " and "))).
		Find(&orderQuoteItemsJoin).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	// set or assign to entity
	for _, oqi := range orderQuoteItemsJoin {
		orderQuoteItems = append(orderQuoteItems, entity.OrderQuoteItem{
			ID:                   oqi.ID,
			QuoteMerchantID:      oqi.QuoteMerchantID,
			ProductID:            oqi.ProductID,
			ProductSku:           oqi.ProductSku,
			MerchantSku:          oqi.MerchantSku,
			CategoryID:           oqi.CategoryID,
			BrandID:              oqi.BrandID,
			Name:                 oqi.Name,
			ItemNotes:            oqi.ItemNotes,
			Selected:             oqi.Selected,
			Weight:               oqi.Weight,
			Quantity:             oqi.Quantity,
			Price:                oqi.Price,
			DiscountPercentage:   oqi.DiscountPercentage,
			DiscountAmount:       oqi.DiscountAmount,
			RowWeight:            oqi.RowWeight,
			OriginalPrice:        oqi.OriginalPrice,
			RowOriginalPrice:     oqi.RowOriginalPrice,
			BrandCode:            oqi.BrandCode,
			BrandName:            oqi.BrandName,
			ProductType:          oqi.ProductType,
			QuoteID:              oqi.QuoteID,
			AttributeSetID:       oqi.AttributeSetID,
			MerchantSpecialPrice: oqi.MerchantSpecialPrice,
			Product:              &entitycatalog.Product{ID: oqi.PID, Sku: oqi.PSku, Slug: oqi.PSlug},
			ProductFlat:          &entitycatalog.ProductFlat{IsActive: oqi.PfIsActive, Status: oqi.PfStatus},
			Merchant:             &entitymerchant.Merchant{ID: oqi.MID, MerchantUID: oqi.MMerchantUID, MerchantCode: oqi.MMerchantCode, MerchantName: oqi.MMerchantName},
		})
	}

	return &orderQuoteItems, nil
}

func (r *orderQuoteItemRepository) UpdateByID(dbc *base.DBContext, id uint64, data entity.OrderQuoteItem) error {
	if id == 0 {
		return errors.New("id is required")
	}
	err := dbc.DB.WithContext(dbc.Context).
		Model(entity.OrderQuoteItem{}).
		Select("*").Omit("id", "created_at").
		Where("id = ?", id).
		Updates(data).
		Error
	return err
}

func (r *orderQuoteItemRepository) DeleteByID(dbc *base.DBContext, id uint64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	err := dbc.DB.WithContext(dbc.Context).
		Where("id = ?", id).
		Delete(entity.OrderQuoteItem{}).
		Error
	return err
}
