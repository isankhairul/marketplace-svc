package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	entitycatalog "marketplace-svc/app/model/entity/catalog"
	base "marketplace-svc/app/repository"
)

type productCategoryRepository struct {
	base.BaseRepository
}

type ProductCategoryRepository interface {
	GetCategoryMenu(dbc *base.DBContext, productID uint64, storeID uint64) ([]string, error)
}

func NewProductCategoryRepository(br base.BaseRepository) ProductCategoryRepository {
	return &productCategoryRepository{br}
}

func (r *productCategoryRepository) GetCategoryMenu(dbc *base.DBContext, productID uint64, storeID uint64) ([]string, error) {
	var arrCategory []entitycatalog.Category
	var arrResponse []string
	query := dbc.DB.WithContext(dbc.Context)

	err := query.Raw(
		fmt.Sprintf(
			`
				 select c.name  
				 from product_category pc 
				 inner join category c on pc.category_id = c.id 
				 left join category parent on parent.id = c.parent_id  
				 where c.status=1 and c.level in (2,3) and c.in_menu=1 and pc.product_id=%s and pc.store_id=%s 
				 and parent.in_menu=1 and parent.status=1 and parent.store_id=1 
				 order by c.level asc, c.name asc
			`, fmt.Sprint(productID), fmt.Sprint(storeID))).
		Find(&arrCategory).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	if len(arrCategory) > 0 {
		for _, pc := range arrCategory {
			arrResponse = append(arrResponse, pc.Name)
		}
	}

	return arrResponse, nil
}
