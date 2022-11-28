package usecase

import (
	"bpe/model"
)

type TUsecaseInterface interface {
	Add(*model.T) error
}
