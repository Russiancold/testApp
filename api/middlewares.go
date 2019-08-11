package api

import (
	"errors"
	"log"
	"net/http"
)

func RecoverMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Unknown error")
				}
				log.Println(err)
				writeInternal(w, err)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
