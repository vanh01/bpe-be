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
	inXorBlock          int
}

type result struct {
	currentCycleTime      float64
	numberOfOptionalTasks int
	numberOfTotalTasks    int
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
	rlt := result{currentCycleTime: 0.0, numberOfOptionalTasks: 0, numberOfTotalTasks: 0}
	evaluateTime := &evaluateAll{}
	contextTime := context{listGateway: make(map[string]int), listGatewayTraveled: make(map[string]interface{})}
	evaluateTime.visit(env.startNode, &contextTime, &rlt)
	rs := make(map[string]interface{})
	rs["cycle_time"] = rlt.currentCycleTime
	rs["numberOfOptionalTasks"] = rlt.numberOfOptionalTasks
	rs["numberOfTotalTasks"] = rlt.numberOfTotalTasks
	rs["flexibility"] = float64(rlt.numberOfOptionalTasks) / float64(rlt.numberOfTotalTasks)
	return rs, nil
}
