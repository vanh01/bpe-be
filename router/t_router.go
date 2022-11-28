package router

import "net/http"

type TRouterInterface interface {
	TTestHandler(w http.ResponseWriter, r *http.Request)
}
