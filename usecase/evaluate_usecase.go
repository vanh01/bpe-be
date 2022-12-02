package usecase

import "bpe/usecase/evaluate"

type EvaluateUsecaseInterface interface {
	EvaluateCycleTime(mapItem map[string]evaluate.Element) (float64, error)
}
