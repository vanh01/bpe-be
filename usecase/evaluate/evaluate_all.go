package evaluate

type evaluateAll struct {
}

func (et *evaluateAll) visit(i interface{}, c *context, r *result) interface{} {
	switch i.(type) {
	case *event:
		return et.visitForEvent(i.(*event), c, r)
	case *task:
		return et.visitForTask(i.(*task), c, r)
	case *gateway:
		return et.visitForGateway(i.(*gateway), c, r)
	}
	r.currentCycleTime = 0
	return nil
}

func (et *evaluateAll) visitForEvent(e *event, c *context, r *result) interface{} {
	// fmt.Printf("Visit: %-50s| %-20s\n", e.Id, e.Name)
	var totalCycleTime = 0.0
	for _, n := range e.Next {
		nextNode := et.visit(n, c, r)
		nextResult := r.currentCycleTime
		totalCycleTime += nextResult
		et.calculateCyclyTimeNextNode(nextNode, c, r)
		totalCycleTime += r.currentCycleTime
	}
	r.currentCycleTime = totalCycleTime
	return nil
}

func (et *evaluateAll) visitForTask(tk *task, c *context, r *result) interface{} {
	// fmt.Printf("Visit: %-50s| %-20s\n", tk.Id, tk.Name)
	nextNode := et.visit(tk.Next[0], c, r)
	r.numberOfTotalTasks += 1
	if c.inXorBlock > 0 {
		r.numberOfOptionalTasks += 1
	}
	nextResult := r.currentCycleTime
	totalCycleTime := tk.CycleTime + nextResult
	et.calculateCyclyTimeNextNode(nextNode, c, r)
	totalCycleTime += r.currentCycleTime
	r.currentCycleTime = totalCycleTime
	return nil
}

func (et *evaluateAll) visitForGateway(g *gateway, c *context, r *result) interface{} {
	// fmt.Printf("Visit: %-50s| %-20s\n", g.Id, g.Name)
	c.listGatewayTraveled[g.Id] = g

	if g.isJoinGateway() {
		return et.handleForJoinGateway(g, c, r)
	}
	if g.isSplitGateway() {
		return et.handleForSplitGateway(g, c, r)
	}
	r.currentCycleTime = 0
	return nil
}

func (et *evaluateAll) handleForJoinGateway(g *gateway, c *context, r *result) interface{} {
	if _, check := c.listGateway[g.Id]; check {
		c.listGateway[g.Id] += 1
	} else {
		c.listGateway[g.Id] = 1 + numberOfGatewayInNodes(g.Previous)
	}
	// check so lan da duyet cua cong join
	if c.listGateway[g.Id] < len(g.Previous) {
		r.currentCycleTime = 0
		return nil
	}
	// kiem tra xem day la mot gateway bat dau khoi loop hay khong
	if check, previous := et.checkNodeTravel(g.Previous, c); !check {
		// fmt.Println("Start loop!")
		c.stackEndLoop.Push(previous.(*gateway))
		return et.handleForLoop(g, previous, c, r)
	}
	// fmt.Println("End gateway!")
	c.stackNextGateway.Push(g)
	r.currentCycleTime = 0
	return nil
}

func (et *evaluateAll) handleForSplitGateway(g *gateway, c *context, r *result) interface{} {
	var totalCycleTime = 0.0
	var nextNode interface{}
	// xu li cho gateway dong loop
	if c.stackEndLoop.Size() > 0 && len(g.Next) == 2 && c.stackEndLoop.Top() == g {
		// fmt.Println("End loop!")
		c.stackEndLoop.Pop()
		r.currentCycleTime = 0
		return nil
	}
	// fmt.Println("Start gateway!")
	// xu li cho split gateway binh thuong
	if g.Name == "ExclusiveGateway" {
		c.inXorBlock += 1
	}
	for i, branch := range g.Next {
		nextN := et.visit(branch, c, r)
		branchCycleTime := r.currentCycleTime
		et.calculateCyclyTimeNextNode(nextN, c, r)
		branchCycleTime += r.currentCycleTime
		nextNode = nil
		switch g.Name {
		case "ParallelGateway":
			if totalCycleTime < branchCycleTime {
				totalCycleTime = branchCycleTime
			}
		case "InclusiveGateway":
			// TODO: Handle or gateway
			r.currentCycleTime = 0
			return nil
		case "ExclusiveGateway":
			totalCycleTime += g.branchingProbabilities[i] * branchCycleTime
		}
	}
	if g.Name == "ExclusiveGateway" {
		c.inXorBlock -= 1
	}
	if c.stackNextGateway.IsNotEmpty() {
		nextNode, _ = c.stackNextGateway.Pop()
		nextNode = nextNode.(*gateway).Next[0]
	}
	r.currentCycleTime = totalCycleTime
	return nextNode
}

// check xem co gateway da duoc duyet hay chua
func (et *evaluateAll) checkNodeTravel(nodes []interface{}, c *context) (bool, interface{}) {
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
func (et *evaluateAll) calculateCyclyTimeNextNode(nextNode interface{}, c *context, r *result) {
	timeResult := 0.0
	for nextNode != nil {
		nextNextNode := et.visit(nextNode, c, r)
		nextNextResult := r.currentCycleTime
		nextNode = nextNextNode
		timeResult += nextNextResult
	}
	r.currentCycleTime = timeResult
}

// xu li tinh toan cho loop
func (et *evaluateAll) handleForLoop(start interface{}, end interface{}, c *context, r *result) interface{} {
	startGateway := start.(*gateway)
	endGateway := end.(*gateway)
	// timeResult, _ := et.visit(startGateway.Next[0], c, r)
	et.visit(startGateway.Next[0], c, r)
	timeResult := r.currentCycleTime
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
	r.currentCycleTime = timeResult / (1 - reloop)
	return nextNode
}
