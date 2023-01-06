package evaluate

import (
	"bpe/usecase"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type evaluateRouter struct {
	iUsecase usecase.EvaluateUsecaseInterface
}

func NewEvaluateRouter(iUsecase usecase.EvaluateUsecaseInterface) *evaluateRouter {
	return &evaluateRouter{iUsecase: iUsecase}
}

func (t *evaluateRouter) EvaluateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s %s %s\n", time.Now().In(time.FixedZone("UTC+7", +7*60*60)), r.Proto, r.RequestURI, r.Method)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	result := t.iUsecase.Evaluate(body)
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
