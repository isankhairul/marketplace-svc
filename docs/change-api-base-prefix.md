# Change your Base-Prefix for API-Endpoints

You can now change your base-prefix dynamically from `environment-variable` or `config-xx.yml` without any hard-coding. This will also dynamically change the swagger `basePath` accordingly.  

You can do it with any of these 2 approaches:  
1. From environment variables, add 2 new env-vars:
   ```bash
   ## For your swagger-path, for example: http://localhost:5600/boilerplate-svc/docs
   KD_URL_BASEPATH=/boilerplate-svc/

   ## For your API-endpoints, for example: http://localhost:5600/boilerplate-svc/api/v1/products
   KD_URL_BASEPREFIX=/boilerplate-svc/api/v1/
   ```

2. From `config-xx.yml`, update these 2 properties:
   ```yaml
   url:
      basepath: /boilerplate-svc/
      baseprefix: /boilerplate-svc/api/v1/
   ```

Sample command to tests:
```bash
KD_URL_BASEPATH=/bp-svc/ KD_URL_BASEPREFIX=/bp-svc/api/v1/ go run main.go
```

## Implementation Changes for Old Boilerplate

This is the instruction for implementing dynamic base-prefix  in old boilerplate:

1. Set values on config.yml

    ```yaml
    url:
        basepath: /boilerplate-svc/
        baseprefix: /boilerplate-svc/api/v1/
    ```

2. Set value of config file in swaggerHttpHandler and append it to your URL paths

    ```golang
    basePath := config.GetConfigString(viper.GetString("url.basepath"))
	basePrefix := config.GetConfigString(viper.GetString("url.baseprefix"))
    ```

3. Read swagger.yaml and manipulate basePath value based on baseprefix

    ```golang
    // Handling & Manipulate swagger.yaml basePath with config-val
	pr.HandleFunc(basePath+"swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		fileBytes, err := ioutil.ReadFile("swagger.yaml")
		if err != nil {
			panic(err)
		}

		regex, _ := regexp.Compile(`^basePath\s*:\s+.*`)
		fileBytes = regex.ReplaceAll(fileBytes, []byte("basePath: "+basePrefix))

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/yaml")
		w.Write(fileBytes)
	})
    ```

4. Set route of swagger ui base on basePath

    ```golang
    opts := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml", BasePath: basePath}
	sh := middleware.SwaggerUI(opts, nil)
	pr.Handle(basePath+"docs", sh)

	//// documentation for share
	opts1 := middleware.RedocOpts{SpecURL: "swagger.yaml", BasePath: basePath, Path: "doc"}
	sh1 := middleware.Redoc(opts1, nil)
	pr.Handle(basePath+"doc", sh1)

	return pr
    ```

### Swagger Handler Code

api/http/swagger_http.go
```golang
func SwaggerHttpHandler(logger log.Logger) http.Handler {
    pr := mux.NewRouter()
    basePath := config.GetConfigString(viper.GetString("url.basepath"))
    basePrefix := config.GetConfigString(viper.GetString("url.baseprefix"))

    // Handling & Manipulate swagger.yaml basePath with config-val
    pr.HandleFunc(basePath+"swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
        fileBytes, err := ioutil.ReadFile("swagger.yaml")
        if err != nil {
            panic(err)
        }

        regex, _ := regexp.Compile(`^basePath\s*:\s+.*`)
        fileBytes = regex.ReplaceAll(fileBytes, []byte("basePath: "+basePrefix))

        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "text/yaml")
        w.Write(fileBytes)
    })
    opts := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml", BasePath: basePath}
    sh := middleware.SwaggerUI(opts, nil)
    pr.Handle(basePath+"docs", sh)

    //// documentation for share
    opts1 := middleware.RedocOpts{SpecURL: "swagger.yaml", BasePath: basePath, Path: "doc"}
    sh1 := middleware.Redoc(opts1, nil)
    pr.Handle(basePath+"doc", sh1)

    return pr
}   
```