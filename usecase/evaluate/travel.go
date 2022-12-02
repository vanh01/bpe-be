package evaluate

type travel struct {
}

var listGateway map[string]int = make(map[string]int)
var listGatewayTraveled map[string]interface{} = make(map[string]interface{})
var stackNextGateway gateWayStack
var stackEndLoop gateWayStack

func (t *travel) visit(i interface{}) (float64, interface{}) {
	switch i.(type) {
	case *event:
		return t.visitForEvent(i.(*event))
	case *task:
		return t.visitForTask(i.(*task))
	case *gateway:
		return t.visitForGateway(i.(*gateway))
	}
	return 0, nil
}

func (t *travel) visitForEvent(e *event) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", e.Id, e.Name)
	var totalCycleTime = 0.0
	for _, n := range e.Next {
		nextResult, nextNode := t.visit(n)
		totalCycleTime += nextResult
		totalCycleTime += calculateCyclyTimeNextNode(t, nextNode)
	}
	return totalCycleTime, nil
}

func (t *travel) visitForTask(tk *task) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", tk.Id, tk.Name)
	nextResult, nextNode := t.visit(tk.Next[0])
	totalCycleTime := tk.CycleTime + nextResult
	totalCycleTime += calculateCyclyTimeNextNode(t, nextNode)
	return totalCycleTime, nil
}

func (t *travel) visitForGateway(g *gateway) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", g.Id, g.Name)
	listGatewayTraveled[g.Id] = g

	if g.isJoinGateway() {
		return handleForJoinGateway(t, g)
	}
	if g.isSplitGateway() {
		return handleForSplitGateway(t, g)
	}
	return 0.0, nil
}

func handleForJoinGateway(t *travel, g *gateway) (float64, interface{}) {
	if _, check := listGateway[g.Id]; check {
		listGateway[g.Id] += 1
	} else {
		listGateway[g.Id] = 1 + numberOfGatewayInNodes(g.Previous)
	}
	// check so lan da duyet cua cong join
	if listGateway[g.Id] < len(g.Previous) {
		return 0, nil
	}
	// kiem tra xem day la mot gateway bat dau khoi loop hay khong
	if check, previous := checkNodeTravel(g.Previous); !check {
		// fmt.Println("Start loop!")
		stackEndLoop.Push(previous.(*gateway))
		return handleForLoop(t, g, previous)
	}
	// fmt.Println("End gateway!")
	stackNextGateway.Push(g)
	return 0, nil
}

func handleForSplitGateway(t *travel, g *gateway) (float64, interface{}) {
	var totalCycleTime = 0.0
	var nextNode interface{}
	// xu li cho gateway dong loop
	if stackEndLoop.Size() > 0 && len(g.Next) == 2 && stackEndLoop.Top() == g {
		// fmt.Println("End loop!")
		stackEndLoop.Pop()
		return 0, nil
	}
	// fmt.Println("Start gateway!")
	// xu li cho split gateway binh thuong
	for i, branch := range g.Next {
		branchCycleTime, nextN := t.visit(branch)
		branchCycleTime += calculateCyclyTimeNextNode(t, nextN)
		nextNode = nil
		switch g.Name {
		case "ParallelGateway":
			if totalCycleTime < branchCycleTime {
				totalCycleTime = branchCycleTime
			}
		case "InclusiveGateway":
			// TODO: Handle or gateway
			return 0, nil
		case "ExclusiveGateway":
			totalCycleTime += g.branchingProbabilities[i] * branchCycleTime
		}
	}
	if stackNextGateway.IsNotEmpty() {
		nextNode, _ = stackNextGateway.Pop()
		nextNode = nextNode.(*gateway).Next[0]
	}
	return totalCycleTime, nextNode
}

// tinh so luong gateway
func numberOfGatewayInNodes(nodes []interface{}) int {
	count := 0
	for _, n := range nodes {
		if isGateway(n) {
			count += 1
		}
	}
	return count
}

// check xem co gateway da duoc duyet hay chua
func checkNodeTravel(nodes []interface{}) (bool, interface{}) {
	for _, n := range nodes {
		if isGateway(n) {
			id := n.(*gateway).Id
			if _, check := listGatewayTraveled[id]; !check {
				return false, n
			}
		}
	}
	return true, nil
}

func calculateCyclyTimeNextNode(t *travel, nextNode interface{}) float64 {
	timeResult := 0.0
	for nextNode != nil {
		nextNextResult, nextNextNode := t.visit(nextNode)
		nextNode = nextNextNode
		timeResult += nextNextResult
	}
	return timeResult
}

func isGateway(g interface{}) bool {
	switch g.(type) {
	case *gateway:
		return true
	}
	return false
}

// xu li tinh toan cho loop
func handleForLoop(t *travel, start interface{}, end interface{}) (float64, interface{}) {
	startGateway := start.(*gateway)
	endGateway := end.(*gateway)
	timeResult, _ := t.visit(startGateway.Next[0])
	var reloop float64
	var nextNode interface{}

	for i, n := range endGateway.Next {
		if isGateway(n) && start == n {
			reloop = endGateway.branchingProbabilities[i]
		}
	}
	if isGateway(endGateway.Next[0]) {
		nextNode = endGateway.Next[1]
	} else {
		nextNode = endGateway.Next[0]
	}
	return timeResult / (1 - reloop), nextNode
}
