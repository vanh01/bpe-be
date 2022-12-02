package evaluate

import (
	"bpe/usecase"
	"bpe/usecase/evaluate"
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

func (t *evaluateRouter) EvaluateCycleTimeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s %s %s", time.Now().In(time.FixedZone("UTC+7", +7*60*60)), r.Proto, r.RequestURI, r.Method)
	var items map[string]evaluate.Element
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &items)
	result, _ := t.iUsecase.EvaluateCycleTime(items)
	w.Write([]byte(fmt.Sprint(result)))
}
