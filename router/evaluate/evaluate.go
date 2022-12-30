package evaluate

import (
	"bpe/usecase"
	"encoding/json"
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
	body, _ := ioutil.ReadAll(r.Body)
	result, _ := t.iUsecase.Evaluate(body)
	jsonResponse, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
