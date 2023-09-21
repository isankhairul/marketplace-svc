package elasticregistry

import (
	"marketplace-svc/app"
	rp "marketplace-svc/app/repository"
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
		app.Log,
		rp.NewBaseRepository(app.DB),
		ec,
	)
}
