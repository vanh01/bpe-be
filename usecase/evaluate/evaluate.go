package evaluate

import (
	"encoding/json"
)

type evaluateUsecase struct{}

func NewEvaluateUsecase() *evaluateUsecase {
	return &evaluateUsecase{}
}

type Element struct {
	Id                     string    `json:"id"`
	Name                   string    `json:"name"`
	Incoming               []string  `json:"incoming"`
	Outgoing               []string  `json:"outgoing"`
	Type                   string    `json:"type"`
	CycleTime              float64   `json:"cycleTime"`
	BranchingProbabilities []float64 `json:"branchingProbabilities"`
}

type context struct {
	listGateway         map[string]int
	listGatewayTraveled map[string]interface{}
	stackNextGateway    gateWayStack
	stackEndLoop        gateWayStack
}

func createNodeList(mapElement map[string]Element) map[string]interface{} {
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

func getNodeOfIds(ids []string, mapElement map[string]Element, mapNodeCreated *map[string]interface{}) []interface{} {
	var result []interface{}
	for _, id := range ids {
		if _, check := mapElement[id]; check {
			result = append(result, (*mapNodeCreated)[id])
		}
	}
	return result
}

func buildGraph(mapElement map[string]Element, mapNodeCreated *map[string]interface{}) *event {
	var startNode *event
	for _, i := range mapElement {
		if i.Type == "event" {
			event := (*mapNodeCreated)[i.Id].(*event)
			event.Previous = getNodeOfIds(i.Incoming, mapElement, mapNodeCreated)
			event.Next = getNodeOfIds(i.Outgoing, mapElement, mapNodeCreated)
			if event.Name == "StartEvent" {
				startNode = event
			}
		} else if i.Type == "task" {
			task := (*mapNodeCreated)[i.Id].(*task)
			task.Previous = getNodeOfIds(i.Incoming, mapElement, mapNodeCreated)
			task.Next = getNodeOfIds(i.Outgoing, mapElement, mapNodeCreated)
		} else if i.Type == "gateway" {
			gateway := (*mapNodeCreated)[i.Id].(*gateway)
			gateway.Previous = getNodeOfIds(i.Incoming, mapElement, mapNodeCreated)
			gateway.Next = getNodeOfIds(i.Outgoing, mapElement, mapNodeCreated)
		}
	}
	return startNode
}

func (e *evaluateUsecase) EvaluateCycleTime(body []byte) (float64, error) {
	var mapElement map[string]Element
	json.Unmarshal(body, &mapElement)
	mapNodeCreated := createNodeList(mapElement)
	startNode := buildGraph(mapElement, &mapNodeCreated)
	travel := &travel{}

	result, _ := startNode.accept(travel, &context{listGateway: make(map[string]int), listGatewayTraveled: make(map[string]interface{})})
	return result, nil
}
