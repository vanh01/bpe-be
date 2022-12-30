package evaluate

import (
	"encoding/json"
)

type evaluateUsecase struct{}

func NewEvaluateUsecase() *evaluateUsecase {
	return &evaluateUsecase{}
}

type element struct {
	Id                     string    `json:"id"`
	Name                   string    `json:"name"`
	Incoming               []string  `json:"incoming"`
	Outgoing               []string  `json:"outgoing"`
	Type                   string    `json:"type"`
	CycleTime              float64   `json:"cycleTime"`
	BranchingProbabilities []float64 `json:"branchingProbabilities"`
}

type env struct {
	mapElement     map[string]element
	mapNodeCreated map[string]interface{}
	startNode      *event
}

type context struct {
	listGateway         map[string]int
	listGatewayTraveled map[string]interface{}
	stackNextGateway    gateWayStack
	stackEndLoop        gateWayStack
}

// tao mot map chua node tu cai json ban dau
func createNodeList(mapElement map[string]element) map[string]interface{} {
	mapNodeCreated := make(map[string]interface{})
	for _, i := range mapElement {
		if i.Type == "event" {
			event := event{node: node{Id: i.Id, Name: i.Name}}
			mapNodeCreated[i.Id] = &event
		} else if i.Type == "task" {
			task := task{node: node{Id: i.Id, Name: i.Name}, CycleTime: i.CycleTime}
			mapNodeCreated[i.Id] = &task
		} else if i.Type == "gateway" {
			gateway := gateway{node: node{Id: i.Id, Name: i.Name}, branchingProbabilities: i.BranchingProbabilities}
			mapNodeCreated[i.Id] = &gateway
		}
	}
	return mapNodeCreated
}

// get node from list id
func getNodeOfIds(ids []string, mapNodeCreated *map[string]interface{}) []interface{} {
	var result []interface{}
	for _, id := range ids {
		result = append(result, (*mapNodeCreated)[id])
	}
	return result
}

// lien ket cac node lai voi nhau
func buildGraph(mapElement map[string]element, mapNodeCreated *map[string]interface{}) *event {
	var startNode *event
	for _, i := range mapElement {
		if i.Type == "event" {
			event := (*mapNodeCreated)[i.Id].(*event)
			event.Previous = getNodeOfIds(i.Incoming, mapNodeCreated)
			event.Next = getNodeOfIds(i.Outgoing, mapNodeCreated)
			if event.Name == "StartEvent" {
				startNode = event
			}
		} else if i.Type == "task" {
			task := (*mapNodeCreated)[i.Id].(*task)
			task.Previous = getNodeOfIds(i.Incoming, mapNodeCreated)
			task.Next = getNodeOfIds(i.Outgoing, mapNodeCreated)
		} else if i.Type == "gateway" {
			gateway := (*mapNodeCreated)[i.Id].(*gateway)
			gateway.Previous = getNodeOfIds(i.Incoming, mapNodeCreated)
			gateway.Next = getNodeOfIds(i.Outgoing, mapNodeCreated)
		}
	}
	return startNode
}

func (e *evaluateUsecase) Evaluate(body []byte) (map[string]interface{}, error) {
	env := &env{}
	if err := json.Unmarshal(body, &env.mapElement); err != nil {
		return nil, err
	}
	env.mapNodeCreated = createNodeList(env.mapElement)
	env.startNode = buildGraph(env.mapElement, &env.mapNodeCreated)
	result := make(map[string]interface{})
	result["time"], _ = e.EvaluateCycleTime(env)
	result["quality"], _ = e.EvaluateQuality(env)
	result["flexibility"], _ = e.EvaluateQuality(env)
	result["transparency"], _ = e.EvaluateQuality(env)
	result["exceptionHandling"], _ = e.EvaluateQuality(env)
	return result, nil
}

func (e *evaluateUsecase) EvaluateCycleTime(en *env) (float64, error) {
	evaluateTime := &evaluateTime{}
	contextTime := context{listGateway: make(map[string]int), listGatewayTraveled: make(map[string]interface{})}
	result, _ := evaluateTime.visit(en.startNode, &contextTime)
	return result, nil
}

func (e *evaluateUsecase) EvaluateQuality(en *env) (float64, error) {
	// evaluateQuality := &evaluateQuality{}
	// contextQuality := context{listGateway: make(map[string]int), listGatewayTraveled: make(map[string]interface{})}
	// result, _ := evaluateQuality.visit(en.startNode, &contextQuality)
	return 0.0, nil
}

func (e *evaluateUsecase) EvaluateFlexibility(en *env) (float64, error) {
	// evaluateFlexibility := &evaluateFlexibility{}
	// contextFlexibility := context{listGateway: make(map[string]int), listGatewayTraveled: make(map[string]interface{})}
	// evaluateFlexibility.visit(en.startNode, &contextFlexibility)
	return 0.0, nil
}

func (e *evaluateUsecase) EvaluateTransparency(en *env) (float64, error) {
	// evaluateTransparency := &evaluateTransparency{}
	// contextTransparency := context{listGateway: make(map[string]int), listGatewayTraveled: make(map[string]interface{})}
	// result, _ := evaluateTransparency.visit(en.startNode, &contextTransparency)
	return 0.0, nil
}

func (e *evaluateUsecase) EvaluateExceptionHandling(en *env) (float64, error) {
	// evaluateExceptionHandling := &evaluateExceptionHandling{}
	// contextExceptionHandling := context{listGateway: make(map[string]int), listGatewayTraveled: make(map[string]interface{})}
	// result, _ := evaluateExceptionHandling.visit(en.startNode, &contextExceptionHandling)
	return 0.0, nil
}
