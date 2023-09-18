# Handle Panic with Recover

This guide will help you implement panic handling and recovery in Golang.

### Please follow these steps below

###  Middleware

create a file at `app/api/middleware/recovery.go`  and copy the code below then save it.

```golang
package middleware

import (
	"encoding/json"
	"marketplace-svc/app/model/base"
	"marketplace-svc/helper/message"
	"log"
	"net/http"
	"runtime/debug"
)

func ServeHTTP(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(
					base.SetHttpResponse(req.Context(), message.FailedMsg.Code, message.FailedMsg.Message, nil, nil),
				)
				log.Println(err)
				debug.PrintStack()
			}
		}()
		h.ServeHTTP(w, req)
	})
}
```

### Update the main.go file

from `main.go`

find the line of code `http.Handle("/", mux)`

Replace to `http.Handle("/", middleware.ServeHTTP(mux))`

