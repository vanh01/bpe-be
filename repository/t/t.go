package t

import (
	"bpe/model"
	"fmt"
)

type tRepository struct{}

func NewTRepository() *tRepository {
	return &tRepository{}
}

func (tRepository *tRepository) Add(t *model.T) error {
	fmt.Println("repository ne")
	return nil
}
