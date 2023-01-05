package usecase

type EvaluateUsecaseInterface interface {
	Evaluate([]byte) []byte
}
