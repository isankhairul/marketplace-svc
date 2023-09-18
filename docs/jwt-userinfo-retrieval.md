# Retrieve JWT User info

This is the implementation steps of retrieving username and their UID from JWT Authorization header:

## Usage

1. There is new helper file [helper/global/jwt.go](#helper-code) introduced. You'll need to include it in your project.

2. Add a `ServerBefore` listener in go-kit transport to extract HttpContext in options at your httpHandler  

   app/api/transport/sample_product_http.go
   ```golang
   func ProductHttpHandler(s service.ProductService, logger log.Logger) http.Handler {
      pr := mux.NewRouter()

      ep := endpoint.MakeProductEndpoints(s)
      options := []httptransport.ServerOption{
         httptransport.ServerErrorLogger(logger),
         httptransport.ServerErrorEncoder(encoder.EncodeError),

         // Listener for Extract JWT from HTTP to Context
         httptransport.ServerBefore(jwt.HTTPToContext()),
      }
   }
   ```

3. Add column to `/model/request/{struct}`

   app/model/request/product_request.go
   ```golang
   type SaveProductRequest struct {
      Name string `json:"name"`
      Sku string `json:"sku"`
      Uom string `json:"uom"`
      Weight int32 `json:"weight"`
      Uid string `json:"uid" binding:"omitempty"`

      // Extend Jwt Info
      global.JWTInfo
   }
   ```

4. Call jwt-helper in your `/api/endpoint/{endpoint}` to get the JWTInfo struct

   app/api/endpoint/sample_product_endpoint.go
   ```golang
   func makeSaveProduct(s service.ProductService) endpoint.Endpoint {
      return func(ctx context.Context, rqst interface{}) (resp interface{}, err error) {
         req := rqst.(request.SaveProductRequest)

         // Sample implement JWT Info
         // Retrieve JWT Info
         jwtInfo, msg := global.SetJWTInfoFromContext(ctx)
         if msg.Code != message.SuccessMsg.Code {
            return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
         }

         // Set Jwt Info to Request Params
         req.JWTInfo = jwtInfo

         result, msg := s.CreateProduct(req)
         if msg.Code == 4000 {
            return base.SetHttpResponse(msg.Code, msg.Message, nil, nil), nil
         }

         return base.SetHttpResponse(msg.Code, msg.Message, result, nil), nil
      }
   }
   ```

5. Set the value of the JWTInfo struct to your `/model/request/{struct}`

   app/api/endpoint/sample_product_endpoint.go
   ```golang
   // Set Jwt Info to Request Params
   req.JWTInfo = jwtInfo
   ```

6. Final step is to set the JWTInfo from your `/model/request/{struct}` to your `/model/entity/{struct}` to persist the data

   app/service/sample_product_service.go
   ```golang
   func (s *sampleProductServiceImpl) CreateSampleProduct(input request.SaveSampleProductRequest) (*entity.SampleProduct, message.Message) {
      logger := log.With(s.logger, "SampleProductService", "CreateSampleProduct")
      s.baseRepo.BeginTx()
      //Set request to entity
      product := entity.SampleProduct{
         Name:   input.Name,
         Sku:    input.Sku,
         Uom:    input.Uom,
         Weight: input.Weight,
         
         // Set Jwt Info to Entity
         BaseIDModel: base.BaseIDModel{
            CreatedByUid: input.ActorUID,
            CreatedBy:    input.ActorName,
         },
      }

      result, err := s.sampleProductRepo.Create(&product)
   }   
   ```


## Skip JWT validation	
Sometimes for local testing, you may want to skip the JWT validation. You can skip it with this option in `config-xxx.yml`

```yaml
security:
  jwt:
    skip-validation: false  ## set it to true to skip JWT validation
```
Note : skip-validation is a flag to check value of JWT in request.


## Helper Code

helper/global/jwt.go
```golang
package global

import (
   "context"
   "fmt"

   "github.com/go-kit/kit/auth/jwt"
   jwtgo "github.com/golang-jwt/jwt"
)

type JWTInfo struct {
   ActorName string `json:"actor_name"`
   ActorUID string `json:"actor_uid"`
}

func SetJWTInfoFromContext(ctx context.Context) JWTInfo {
   jwtInfo := JWTInfo{}
   skipJWTValidation := config.GetConfigBool(viper.GetString("security.jwt.skip-validation"))

   if !skipJWTValidation {
      token, _, err := new(jwtgo.Parser).ParseUnverified(fmt.Sprint(ctx.Value(jwt.JWTContextKey)), jwtgo.MapClaims{})
      if err != nil {
         return &jwtInfo, message.ErrNoAuth
      }

      if claims, ok := token.Claims.(jwtgo.MapClaims); ok {
         jwtInfo.ActorName = fmt.Sprintf("%v", claims["name"])
         jwtInfo.ActorUID = fmt.Sprintf("%v", claims["sub"])
         return &jwtInfo, message.SuccessMsg
      } else {
         return &jwtInfo, message.ErrNoAuth
      }

   } else {
      return &jwtInfo, message.SuccessMsg
   }
}
```
 if skip-validation is false user need to add token in request, if not will return 401 code. 

