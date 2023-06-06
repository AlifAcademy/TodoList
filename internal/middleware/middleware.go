package middleware

import (
	"context"
	"net/http"
	"github.com/AlifAcademy/TodoList/pkg/types"
)


// Basic middleware
func Basic(auth func(login, password string) (int64, bool)) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			basicLogin, basicPassword, ok := request.BasicAuth()
			if !ok {
				http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			userID, isOk := auth(basicLogin, basicPassword)
			if  !isOk{
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(request.Context(), types.Key("key"), userID)
			handler.ServeHTTP(writer, request.WithContext(ctx))
		})
	}
}