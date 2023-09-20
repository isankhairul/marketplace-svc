package registry

import (
	"marketplace-svc/app"
	rp "marketplace-svc/app/repository"
	elasticservice "marketplace-svc/app/service/elastic"
	"marketplace-svc/helper/elastic"
)

func RegisterElasticProductService(app *app.Infra) elasticservice.ElasticProductService {
	ec, _ := elastic.NewElasticClient(&app.Config.Elastic, app.Log)
	return elasticservice.NewElasticProductService(
		app.Log,
		rp.NewBaseRepository(app.DB),
		rp.NewProductFlatRepository(app.DB),
		ec,
	)
}
