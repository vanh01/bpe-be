package router

import "net/http"

type EvaluateRouterInterface interface {
	EvaluateCycleTimeHandler(w http.ResponseWriter, r *http.Request)
}
