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
	inLoop              int
	inBlock             int
}

// export result
type blockCycleTime struct {
	Text   string           `json:"text"`
	Blocks []blockCycleTime `json:"blocks"`
}

type blockQuality struct {
	Text              string  `json:"text"`
	Start             string  `json:"start"`
	End               string  `json:"end"`
	CycleTime         float64 `json:"cycleTime"`
	ReworkProbability float64 `json:"reworkProbability"`
}

type blockFlexibility struct {
	Text    string   `json:"text"`
	TaskIDs []string `json:"taskIDs"`
}

type result struct {
	CurrentCycleTime       float64            `json:"currentCycleTime"`
	NumberOfOptionalTasks  int                `json:"numberOfOptionalTasks"`
	NumberOfTotalTasks     int                `json:"numberOfTotalTasks"`
	TotalCycleTimeAllLoops float64            `json:"totalCycleTimeAllLoops"`
	LogsCycleTime          []blockCycleTime   `json:"logsCycleTime"`
	LogsQuality            []blockQuality     `json:"logsQuality"`
	LogsFlexibility        []blockFlexibility `json:"logsFlexibility"`
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

func (e *evaluateUsecase) Evaluate(body []byte) []byte {
	env := &env{}
	if err := json.Unmarshal(body, &env.mapElement); err != nil {
		return nil
	}
	env.mapNodeCreated = createNodeList(env.mapElement)
	env.startNode = buildGraph(env.mapElement, &env.mapNodeCreated)
	rlt := result{
		CurrentCycleTime:       0.0,
		NumberOfOptionalTasks:  0,
		NumberOfTotalTasks:     0,
		TotalCycleTimeAllLoops: 0.0,
		LogsCycleTime:          []blockCycleTime{},
		LogsQuality:            []blockQuality{},
		LogsFlexibility:        []blockFlexibility{},
	}
	evaluateTime := &traverse{}
	contextTime := context{listGateway: make(map[string]int), listGatewayTraveled: make(map[string]interface{}), inXorBlock: 0, inLoop: 0}
	evaluateTime.visit(env.startNode, &contextTime, &rlt)
	type export struct {
		result
		Flexibility float64 `json:"flexibility"`
		Quality     float64 `json:"quality"`
	}
	rs := export{
		result:      rlt,
		Flexibility: float64(rlt.NumberOfOptionalTasks) / float64(rlt.NumberOfTotalTasks),
		Quality:     rlt.TotalCycleTimeAllLoops / rlt.CurrentCycleTime,
	}
	result, err := json.Marshal(rs)
	if err != nil {
		return nil
	}
	return result
}
