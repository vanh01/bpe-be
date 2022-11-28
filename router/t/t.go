package t

import (
	"bpe/model"
	"bpe/usecase"
	"net/http"
)

type TRouter struct {
	iUsecase usecase.TUsecaseInterface
}

func NewTRouter(iUsecase usecase.TUsecaseInterface) *TRouter {
	return &TRouter{iUsecase: iUsecase}
}

func (t *TRouter) TTestHandler(w http.ResponseWriter, r *http.Request) {
	t.iUsecase.Add(&model.T{T: "123"})
	w.Write([]byte("Gorilla!\n"))
}
