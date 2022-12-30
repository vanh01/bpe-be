package router

import (
	tRepository "bpe/repository/t"
	evaluateRouter "bpe/router/evaluate"
	tRouter "bpe/router/t"
	evaluateUsecase "bpe/usecase/evaluate"
	TUsecase "bpe/usecase/t"
	"net/http"

	"github.com/gorilla/mux"
)

func MapRouters(r *mux.Router) {
	tRep := tRepository.NewTRepository()
	tUse := TUsecase.NewTUsecase(tRep)
	tRou := tRouter.NewTRouter(tUse)
	r.HandleFunc("/api/t", tRou.TTestHandler).Methods(http.MethodGet)
	evaluateUse := evaluateUsecase.NewEvaluateUsecase()
	evaluateRou := evaluateRouter.NewEvaluateRouter(evaluateUse)
	r.HandleFunc("/api/evaluate/time", evaluateRou.EvaluateHandler).Methods(http.MethodPost)
}
