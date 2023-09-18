# marketplace-svc

Marketplace Service

This is boilerplate for klikdokter microservice, it's use go-kit as our base library
The example to use this boilerplate please see this code
[Klikdokter Boilerplate Example](https://gitlab.com/medkomtek/scratchpad/go-kit-boilerplate)

## How to run

- git config core.hooksPath .githooks
- Run consul agent (e.g consul agent -dev)
- Copy config file example `cp config.yml.example config-dev.yml`
- In a root folder type `make init && make precommit.rehooks`
- In a root folder type `make run-server`
- Swagger in path `/docs`
- Run unit test in root folder : `make test`
- Auto generate latest mock files : `make mocks`

## Configuration File

- All configurations are located in `config-xxx.yml` 
- There are 3 config files used, based on `KD_ENV`:
  - `config-dev.yml`: for local development 
  - `config-prd.yml`: for staging/production (KD_ENV=stg or KD_ENV=prd)
  - `config-test.yml`: for go test
- You can override your configuration by using environment variable with `KD_` prefix. Capitalize your YAML field and replace any . (dot character) with _ (underscore). Examples:
  - KD_SITE_PORT will override site.port 
  - KD_DATABASE_USERNAME will override database.username
  
## Directory Structure

In go naming convention, package/directory names are not plural. This is surprising to programmers who came from other languages and are retaining an old habit of pluralizing names. Don't name a package `httputils`, but `httputil`.

```tree
marketplace-svc-boilerplate
+-- app 
|   +-- api 
|   |   +-- endpoint
|   |   +-- initialization
|   |   +-- transport
|   +-- model
|   |   +-- base
|   |   +-- entity
|   |   +-- presenter 
|   |   +-- request 
|   |   +-- response
|   +-- registry
|   +-- repository
|   +-- service
+-- docker
+-- helper 
+-- pkg
```

### `/marketplace-svc-boilerplate`

This is your main project directory, you can change `marketplace-svc-boilerplate` with any project name you have.
When you change this directory, dont forget to rename it in `go.mod` file,
`module marketplace-svc-boilerplate` --> change this module name with your name,after you change this module name,
don't forget to replace the existing module name with a new module name when importing the package,
`import "marketplace-svc-boilerplate/pkg/logger"` --> replace `marketplace-svc-boilerplate` with your new name.

### `/app`

This is your application directory, put your model, database interaction, business rules, and presentation response in this directory.
In this directory, we have separate directories, each of which represents the layered architecture.

- `/api`
This is your API, middleware JSON schema files, protocol definition files reside.
In this directory, we must have 2 required directory as of go-kit architecture
  - `/endpoint` An endpoint is like an action/handler on a controller; it’s where safety and antifragile logic lives. If you implement two transports (HTTP and gRPC), you might have two methods of sending requests to the same endpoint.
  - `/transport` The transport layer is bound to concrete transports like HTTP or gRPC. If you want to have HTTP and gRPC protocol for the same endpoint, you don't need to create new endpoint, but you only have create a new transport protocol.
  - `/initialization` Use this file for your initialization for transport, service, routing and database. Please don't make changes on main.go file
- `/model` A model in Go is a set of data structures, instead, we separate pure database table models and custom models.
  - `/base` This directory is where you put your base models, for example, if you have a standard fields that will be used on many tables you can add a base model and you can use it in your models, this concept is like an inheritance.
  - `/entity` This is where your database table models reside.
  - `/request` This is where your request models reside, you can also add a validation.
  - `/response` Sometimes we need to present data that has custom fields or not the same as the database table, this is where your view models reside. For our HTTP response please use function `SetResponse`
- `/registry` This is where you registering your service, so we can use it.
- `/repository` This layer where you put your query or interaction with the database.
- `/service` This is where you put your business rules in action, It encapsulates and implements all of the use cases of the system,one `service` file can use many repositories.
  
### `/docker`

This is where your docker files & script that use on docker reside.

### `/helper`

If you need to create some function that can be used on the entire project, you can put your code here.
This helper function is bonded to your application.

### `/docs`

You can find documentations about some boilerplate features that have been introduced. For example:
- [Retrive JWT Info from Authorization Header](/docs/jwt-userinfo-retrieval.md)
- [Change your API base-prefix](change-api-base-prefix.md)

### `/pkg`

Library code that's ok to use by external applications. Other projects can import these libraries expecting them to work.
So make sure to put the library code here is not bonded to your application.

## 👨‍💻 List of Libraries

- [go-kit](https://github.com/go-kit/kit) - Microservice library
- [gorilla mux](https://github.com/gorilla/mux) - API NewRouter
- [gorrila schema](https://github.com/gorilla/schema) - Converts structs to and from form values.
- [jwt-go](https://github.com/dgrijalva/jwt-go) - JSON Web Tokens (JWT)
- [gorm](https://gorm.io/gorm) - ORM library
- [redigo](https://github.com/gomodule/redigo) - Redis client
- [ozzo-validation](https://github.com/itgelo/ozzo-validation/v4) - Data validation
- [elasticsearch](https://github.com/olivere/elastic/v7) - Elastic search lbrary
- [curl](https://moul.io/http2curl) - CURL generator based on API request
- [zap](https://github.com/uber-go/zap) - Logger
- [viper](https://github.com/spf13/viper) - Config solution
- [uuid](https://github.com/matoous/go-nanoid/v2) - UID for rest API

Another library can be added later.
