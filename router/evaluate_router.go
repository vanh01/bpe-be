package router

import "net/http"

type EvaluateRouterInterface interface {
	EvaluateHandler(w http.ResponseWriter, r *http.Request)
}
