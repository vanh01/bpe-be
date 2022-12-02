package t

import (
	"bpe/model"
	"bpe/repository"
	"fmt"
)

type tUsecase struct {
	tRepositoryInterface repository.TRepositoryInterface
}

func NewTUsecase(tRepositoryInterface repository.TRepositoryInterface) *tUsecase {
	return &tUsecase{tRepositoryInterface: tRepositoryInterface}
}

func (t *tUsecase) Add(tModel *model.T) error {
	fmt.Println("usecase ne")
	t.tRepositoryInterface.Add(tModel)
	return nil
}
