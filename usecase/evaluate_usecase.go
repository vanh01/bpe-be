package usecase

type EvaluateUsecaseInterface interface {
	Evaluate([]byte) (map[string]interface{}, error)
}
