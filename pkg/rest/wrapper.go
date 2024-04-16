package rest

import (
	"encoding/json"
	"github.com/artshirshov/gastebin/pkg/gerror"
	"net/http"
)

func JsonWrapperHandler[T any](f func(http.ResponseWriter, *http.Request) (ResponseHolder[T], error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var response ResponseHolder[T]
		response, err = f(w, r)
		if err == nil {
			var jsonBytes []byte
			jsonBytes, err = json.MarshalIndent(response.Data, "", " ")
			if err == nil {
				w.WriteHeader(response.Status)
				_, _ = w.Write(jsonBytes)
				return
			}
		}
		if err != nil {
			gerror.Handle(w, err)
		}
	}
}
