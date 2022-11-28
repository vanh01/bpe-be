package t

import (
	"bpe/model"
	"bpe/repository"
	"fmt"
)

type TUsecase struct {
	tRepositoryInterface repository.TRepositoryInterface
}

func NewTUsecase(tRepositoryInterface repository.TRepositoryInterface) *TUsecase {
	return &TUsecase{tRepositoryInterface: tRepositoryInterface}
}

func (t *TUsecase) Add(tModel *model.T) error {
	fmt.Println("usecase ne")
	t.tRepositoryInterface.Add(tModel)
	return nil
}
