package router

import (
	tRepository "bpe/repository/t"
	tRouter "bpe/router/t"
	TUsecase "bpe/usecase/t"
	"net/http"

	"github.com/gorilla/mux"
)

func MapRouters(r *mux.Router) {
	tRep := tRepository.NewTRepository()
	tUse := TUsecase.NewTUsecase(tRep)
	tRou := tRouter.NewTRouter(tUse)
	r.HandleFunc("/api/t", tRou.TTestHandler).Methods(http.MethodGet)
}
