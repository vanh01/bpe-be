package evaluate

type travel struct {
}

func (t *travel) visit(i interface{}, c *context) (float64, interface{}) {
	switch i.(type) {
	case *event:
		return t.visitForEvent(i.(*event), c)
	case *task:
		return t.visitForTask(i.(*task), c)
	case *gateway:
		return t.visitForGateway(i.(*gateway), c)
	}
	return 0, nil
}

func (t *travel) visitForEvent(e *event, c *context) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", e.Id, e.Name)
	var totalCycleTime = 0.0
	for _, n := range e.Next {
		nextResult, nextNode := t.visit(n, c)
		totalCycleTime += nextResult
		totalCycleTime += calculateCyclyTimeNextNode(t, nextNode, c)
	}
	return totalCycleTime, nil
}

func (t *travel) visitForTask(tk *task, c *context) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", tk.Id, tk.Name)
	nextResult, nextNode := t.visit(tk.Next[0], c)
	totalCycleTime := tk.CycleTime + nextResult
	totalCycleTime += calculateCyclyTimeNextNode(t, nextNode, c)
	return totalCycleTime, nil
}

func (t *travel) visitForGateway(g *gateway, c *context) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", g.Id, g.Name)
	c.listGatewayTraveled[g.Id] = g

	if g.isJoinGateway() {
		return handleForJoinGateway(t, g, c)
	}
	if g.isSplitGateway() {
		return handleForSplitGateway(t, g, c)
	}
	return 0.0, nil
}

func handleForJoinGateway(t *travel, g *gateway, c *context) (float64, interface{}) {
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
	if check, previous := checkNodeTravel(g.Previous, c); !check {
		// fmt.Println("Start loop!")
		c.stackEndLoop.Push(previous.(*gateway))
		return handleForLoop(t, g, previous, c)
	}
	// fmt.Println("End gateway!")
	c.stackNextGateway.Push(g)
	return 0, nil
}

func handleForSplitGateway(t *travel, g *gateway, c *context) (float64, interface{}) {
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
		branchCycleTime, nextN := t.visit(branch, c)
		branchCycleTime += calculateCyclyTimeNextNode(t, nextN, c)
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
func checkNodeTravel(nodes []interface{}, c *context) (bool, interface{}) {
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

func calculateCyclyTimeNextNode(t *travel, nextNode interface{}, c *context) float64 {
	timeResult := 0.0
	for nextNode != nil {
		nextNextResult, nextNextNode := t.visit(nextNode, c)
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
func handleForLoop(t *travel, start interface{}, end interface{}, c *context) (float64, interface{}) {
	startGateway := start.(*gateway)
	endGateway := end.(*gateway)
	timeResult, _ := t.visit(startGateway.Next[0], c)
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
