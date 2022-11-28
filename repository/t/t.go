package t

import (
	"bpe/model"
	"fmt"
)

type TRepository struct{}

func NewTRepository() *TRepository {
	return &TRepository{}
}

func (tRepository *TRepository) Add(t *model.T) error {
	fmt.Println("repository ne")
	return nil
}
