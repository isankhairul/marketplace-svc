package elasticregistry

import (
	"marketplace-svc/app"
	rp "marketplace-svc/app/repository"
	rppromo "marketplace-svc/app/repository/promotion"
	elasticservice "marketplace-svc/app/service/elastic"
	"marketplace-svc/helper/elastic"
)

func RegisterEsProductService(app *app.Infra) elasticservice.ElasticProductService {
	ec, _ := elastic.NewElasticClient(&app.Config.Elastic, app.Log)
	return elasticservice.NewElasticProductService(
		app.Log,
		rp.NewBaseRepository(app.DB),
		rp.NewProductFlatRepository(app.DB),
		ec,
	)
}

func RegisterEsBannerService(app *app.Infra) elasticservice.ElasticBannerService {
	ec, _ := elastic.NewElasticClient(&app.Config.Elastic, app.Log)
	return elasticservice.NewElasticBannerService(
		*app.Config,
		app.Log,
		rp.NewBaseRepository(app.DB),
		ec,
	)
}

func RegisterEsBrandService(app *app.Infra) elasticservice.ElasticBrandService {
	ec, _ := elastic.NewElasticClient(&app.Config.Elastic, app.Log)
	return elasticservice.NewElasticBrandService(
		*app.Config,
		app.Log,
		rp.NewBaseRepository(app.DB),
		ec,
	)
}

func RegisterEsCategoryService(app *app.Infra) elasticservice.ElasticCategoryService {
	ec, _ := elastic.NewElasticClient(&app.Config.Elastic, app.Log)
	return elasticservice.NewElasticCategoryService(
		*app.Config,
		app.Log,
		rp.NewBaseRepository(app.DB),
		ec,
	)
}

func RegisterEsVoucherService(app *app.Infra) elasticservice.ElasticVoucherService {
	ec, _ := elastic.NewElasticClient(&app.Config.Elastic, app.Log)
	return elasticservice.NewElasticVoucherService(
		*app.Config,
		app.Log,
		rp.NewBaseRepository(app.DB),
		rppromo.NewPromotionScprCustomerCouponRepository(rp.NewBaseRepository(app.DB)),
		ec,
	)
}

func RegisterEsMerchantService(app *app.Infra) elasticservice.ElasticMerchantService {
	ec, _ := elastic.NewElasticClient(&app.Config.Elastic, app.Log)
	return elasticservice.NewElasticMerchantService(
		*app.Config,
		app.Log,
		rp.NewBaseRepository(app.DB),
		ec,
	)
}

func RegisterEsOrderService(app *app.Infra) elasticservice.ElasticOrderService {
	ec, _ := elastic.NewElasticClient(&app.Config.Elastic, app.Log)
	return elasticservice.NewElasticOrderService(
		*app.Config,
		app.Log,
		rp.NewBaseRepository(app.DB),
		ec,
	)
}
