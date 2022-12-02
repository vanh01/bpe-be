package t

import (
	"bpe/model"
	"bpe/usecase"
	"net/http"
)

type tRouter struct {
	iUsecase usecase.TUsecaseInterface
}

func NewTRouter(iUsecase usecase.TUsecaseInterface) *tRouter {
	return &tRouter{iUsecase: iUsecase}
}

func (t *tRouter) TTestHandler(w http.ResponseWriter, r *http.Request) {
	t.iUsecase.Add(&model.T{T: "123"})
	w.Write([]byte("Gorilla!\n"))
}
