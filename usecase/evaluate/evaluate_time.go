package evaluate

type evaluateTime struct {
}

func (et *evaluateTime) visit(i interface{}, c *context) (float64, interface{}) {
	switch i.(type) {
	case *event:
		return et.visitForEvent(i.(*event), c)
	case *task:
		return et.visitForTask(i.(*task), c)
	case *gateway:
		return et.visitForGateway(i.(*gateway), c)
	}
	return 0, nil
}

func (et *evaluateTime) visitForEvent(e *event, c *context) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", e.Id, e.Name)
	var totalCycleTime = 0.0
	for _, n := range e.Next {
		nextResult, nextNode := et.visit(n, c)
		totalCycleTime += nextResult
		totalCycleTime += et.calculateCyclyTimeNextNode(nextNode, c)
	}
	return totalCycleTime, nil
}

func (et *evaluateTime) visitForTask(tk *task, c *context) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", tk.Id, tk.Name)
	nextResult, nextNode := et.visit(tk.Next[0], c)
	totalCycleTime := tk.CycleTime + nextResult
	totalCycleTime += et.calculateCyclyTimeNextNode(nextNode, c)
	return totalCycleTime, nil
}

func (et *evaluateTime) visitForGateway(g *gateway, c *context) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", g.Id, g.Name)
	c.listGatewayTraveled[g.Id] = g

	if g.isJoinGateway() {
		return et.handleForJoinGateway(g, c)
	}
	if g.isSplitGateway() {
		return et.handleForSplitGateway(g, c)
	}
	return 0.0, nil
}

func (et *evaluateTime) handleForJoinGateway(g *gateway, c *context) (float64, interface{}) {
	if _, check := c.listGateway[g.Id]; check {
		c.listGateway[g.Id] += 1
	} else {
		c.listGateway[g.Id] = 1 + numberOfGatewayInNodes(g.Previous)
	}
	// check so lan da duyet cua cong join
	if c.listGateway[g.Id] < len(g.Previous) {
		return 0, nil
	}
	// kiem tra xem day la mot gateway bat dau khoi loop hay khong
	if check, previous := et.checkNodeTravel(g.Previous, c); !check {
		// fmt.Println("Start loop!")
		c.stackEndLoop.Push(previous.(*gateway))
		return et.handleForLoop(g, previous, c)
	}
	// fmt.Println("End gateway!")
	c.stackNextGateway.Push(g)
	return 0, nil
}

func (et *evaluateTime) handleForSplitGateway(g *gateway, c *context) (float64, interface{}) {
	var totalCycleTime = 0.0
	var nextNode interface{}
	// xu li cho gateway dong loop
	if c.stackEndLoop.Size() > 0 && len(g.Next) == 2 && c.stackEndLoop.Top() == g {
		// fmt.Println("End loop!")
		c.stackEndLoop.Pop()
		return 0, nil
	}
	// fmt.Println("Start gateway!")
	// xu li cho split gateway binh thuong
	for i, branch := range g.Next {
		branchCycleTime, nextN := et.visit(branch, c)
		branchCycleTime += et.calculateCyclyTimeNextNode(nextN, c)
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
	if c.stackNextGateway.IsNotEmpty() {
		nextNode, _ = c.stackNextGateway.Pop()
		nextNode = nextNode.(*gateway).Next[0]
	}
	return totalCycleTime, nextNode
}

// check xem co gateway da duoc duyet hay chua
func (et *evaluateTime) checkNodeTravel(nodes []interface{}, c *context) (bool, interface{}) {
	for _, n := range nodes {
		if isGateway(n) {
			id := n.(*gateway).Id
			if _, check := c.listGatewayTraveled[id]; !check {
				return false, n
			}
		}
	}
	return true, nil
}

// tinh toan cycle time cho nhung next node tiep theo
func (et *evaluateTime) calculateCyclyTimeNextNode(nextNode interface{}, c *context) float64 {
	timeResult := 0.0
	for nextNode != nil {
		nextNextResult, nextNextNode := et.visit(nextNode, c)
		nextNode = nextNextNode
		timeResult += nextNextResult
	}
	return timeResult
}

// xu li tinh toan cho loop
func (et *evaluateTime) handleForLoop(start interface{}, end interface{}, c *context) (float64, interface{}) {
	startGateway := start.(*gateway)
	endGateway := end.(*gateway)
	timeResult, _ := et.visit(startGateway.Next[0], c)
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
