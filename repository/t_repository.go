package repository

import "bpe/model"

type TRepositoryInterface interface {
	Add(*model.T) error
}
