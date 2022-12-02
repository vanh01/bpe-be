package usecase

type EvaluateUsecaseInterface interface {
	EvaluateCycleTime([]byte) (float64, error)
}
