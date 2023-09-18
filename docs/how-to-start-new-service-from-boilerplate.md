# How to start new service from boilerplate

For example, you need create a new service which called "Article"

1. Use the Fork button of Gitlab to fork new repository from boilerplate repository and name it as `aritcle-svc`. Your repository URL will be:
```shell
https://gitlab.klik.doctor/common-services/backend/article-svc.git
```
2. Replace Go module name of your repository (`marketplace-svc`) with the repository URL on `go.mod`
```shell
// module marketplace-svc // remove origin boilerplate module name
module gitlab.klik.doctor/common-services/backend/article-svc
```
3. Remove sample files and snippet code
   - Remove file `app/api/endpoint/sample_product_endpoint.go`
   - Remove file `app/api/transport/sample_product_http.go`
   - Remove file `app/model/entity/sample_product.go`
   - Remove file `app/model/request/sample_product_request.go`
   - Remove file `app/model/response/sample_product_response.go`
   - Remove file `app/service/sample_product_service.go`
   - Remove sample product registry snippet code in `app/registry/service_registry.go`
```golang
// remove this sample product registry snippet code
func RegisterSampleProductService(db *gorm.DB, logger log.Logger) service.SampleProductService {
   return service.NewSampleProductService(
      logger,
      rp.NewSampleProductRepository(db), 
   )
}
```
4. Add files for your service, example: Article service
   - `app/api/endpoint/article_endpoint.go`
   - `app/api/transport/article_http.go`
   - `app/model/entity/article.go`
   - `app/model/request/article_request.go`
   - `app/model/response/article_response.go`
   - `app/service/article_service.go`
   - Update your article service registry in `app/registry/service_registry.go`
```golang
func RegisterArticleService(db *gorm.DB, logger log.Logger) service.ArticleService {
   return service.NewArticleService(
      logger,
      rp.NewArticleRepository(db),
   )
}
```

# Summary
1. Update Go module name in `go.mod`
2. Update Boilerplate layers:
   - Endpoint layer (similar to Controller in Laravel/PHP) `app/api/endpoint/article_endpoint.go`
   - Routes definition `app/api/transport/article_http.go`
   - Entity layer `app/model/article.go`
   - Request layer `app/model/request/article_request.go`
   - Response layer `app/model/response/article_response.go`
   - Service layer `app/service/article_service.go`
   - Registry services `app/registry/service_registry.go`