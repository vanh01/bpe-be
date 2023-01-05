package evaluate

import "fmt"

type traverse struct {
}

func (et *traverse) visit(i interface{}, c *context, r *result) interface{} {
	switch i.(type) {
	case *event:
		return et.visitForEvent(i.(*event), c, r)
	case *task:
		return et.visitForTask(i.(*task), c, r)
	case *gateway:
		return et.visitForGateway(i.(*gateway), c, r)
	}
	r.CurrentCycleTime = 0
	return nil
}

func (et *traverse) visitForEvent(e *event, c *context, r *result) interface{} {
	fmt.Printf("Visit: %-50s| %-20s\n", e.Id, e.Name)
	var totalCycleTime = 0.0
	for _, n := range e.Next {
		nextNode := et.visit(n, c, r)
		nextResult := r.CurrentCycleTime
		totalCycleTime += nextResult
		et.calculateCyclyTimeNextNode(nextNode, c, r)
		totalCycleTime += r.CurrentCycleTime
	}
	r.CurrentCycleTime = totalCycleTime
	return nil
}

func (et *traverse) visitForTask(tk *task, c *context, r *result) interface{} {
	fmt.Printf("Visit: %-50s| %-20s\n", tk.Id, tk.Name)
	if c.inLoop == 0 && c.inBlock == 0 {
		block := blockCycleTime{Text: fmt.Sprintf("Step %d: calculate task %s", len(r.LogsCycleTime)+1, tk.Name), Blocks: nil}
		r.LogsCycleTime = append(r.LogsCycleTime, block)
	} else {
		block := blockCycleTime{Text: fmt.Sprintf("Calculate task %s", tk.Name), Blocks: nil}
		et.addBlockByLevel(r.LogsCycleTime, block, c.inBlock+c.inLoop)
	}
	nextNode := et.visit(tk.Next[0], c, r)
	r.NumberOfTotalTasks += 1
	if c.inXorBlock > 0 {
		r.NumberOfOptionalTasks += 1
	}
	nextResult := r.CurrentCycleTime
	totalCycleTime := tk.CycleTime + nextResult
	et.calculateCyclyTimeNextNode(nextNode, c, r)
	totalCycleTime += r.CurrentCycleTime
	r.CurrentCycleTime = totalCycleTime
	return nil
}

func (et *traverse) visitForGateway(g *gateway, c *context, r *result) interface{} {
	fmt.Printf("Visit: %-50s| %-20s\n", g.Id, g.Name)
	c.listGatewayTraveled[g.Id] = g

	if g.isJoinGateway() {
		return et.handleForJoinGateway(g, c, r)
	}
	if g.isSplitGateway() {
		return et.handleForSplitGateway(g, c, r)
	}
	r.CurrentCycleTime = 0
	return nil
}

func (et *traverse) handleForJoinGateway(g *gateway, c *context, r *result) interface{} {
	if _, check := c.listGateway[g.Id]; check {
		c.listGateway[g.Id] += 1
	} else {
		c.listGateway[g.Id] = 1 + numberOfGatewayInNodes(g.Previous)
	}
	// check so lan da duyet cua cong join
	if c.listGateway[g.Id] < len(g.Previous) {
		r.CurrentCycleTime = 0
		return nil
	}
	// kiem tra xem day la mot gateway bat dau khoi loop hay khong
	if check, previous := et.checkNodeTraveled(g.Previous, c); !check {
		fmt.Println("Start loop!")
		if c.inLoop == 0 && c.inBlock == 0 {
			block := blockCycleTime{Text: fmt.Sprintf("Step %d: calculate loop", len(r.LogsCycleTime)+1), Blocks: nil}
			r.LogsCycleTime = append(r.LogsCycleTime, block)
		} else {
			block := blockCycleTime{Text: "Calculate loop", Blocks: nil}
			et.addBlockByLevel(r.LogsCycleTime, block, c.inBlock+c.inLoop)
		}
		c.inLoop += 1
		c.stackEndLoop.Push(previous.(*gateway))
		return et.handleForLoop(g, previous, c, r)
	}
	fmt.Println("End gateway!")
	c.stackNextGateway.Push(g)
	r.CurrentCycleTime = 0
	return nil
}

func (et *traverse) handleForSplitGateway(g *gateway, c *context, r *result) interface{} {
	var totalCycleTime = 0.0
	var nextNode interface{}
	// xu li cho gateway dong loop
	if c.stackEndLoop.Size() > 0 && len(g.Next) == 2 && c.stackEndLoop.Top() == g {
		fmt.Println("End loop!")
		c.inLoop -= 1
		c.stackEndLoop.Pop()
		r.CurrentCycleTime = 0
		return nil
	}
	fmt.Println("Start gateway!")
	// xu li cho split gateway binh thuong
	if c.inLoop == 0 && c.inBlock == 0 {
		block := blockCycleTime{Text: fmt.Sprintf("Step %d: calculate block %s", len(r.LogsCycleTime)+1, g.Name), Blocks: nil}
		r.LogsCycleTime = append(r.LogsCycleTime, block)
	} else {
		block := blockCycleTime{Text: fmt.Sprintf("Calculate block %s", g.Name), Blocks: nil}
		et.addBlockByLevel(r.LogsCycleTime, block, c.inBlock+c.inLoop)
	}
	c.inBlock += 1
	if g.Name == "ExclusiveGateway" {
		c.inXorBlock += 1
	}
	for i, branch := range g.Next {
		nextN := et.visit(branch, c, r)
		branchCycleTime := r.CurrentCycleTime
		et.calculateCyclyTimeNextNode(nextN, c, r)
		branchCycleTime += r.CurrentCycleTime
		nextNode = nil
		switch g.Name {
		case "ParallelGateway":
			if totalCycleTime < branchCycleTime {
				totalCycleTime = branchCycleTime
			}
		case "InclusiveGateway":
			// TODO: Handle or gateway
			r.CurrentCycleTime = 0
			return nil
		case "ExclusiveGateway":
			totalCycleTime += g.branchingProbabilities[i] * branchCycleTime
		}
	}
	if g.Name == "ExclusiveGateway" {
		c.inXorBlock -= 1
	}
	c.inBlock -= 1
	if c.stackNextGateway.IsNotEmpty() {
		nextNode, _ = c.stackNextGateway.Pop()
		nextNode = nextNode.(*gateway).Next[0]
	}
	r.CurrentCycleTime = totalCycleTime
	return nextNode
}

// check xem co gateway da duoc duyet hay chua
func (et *traverse) checkNodeTraveled(nodes []interface{}, c *context) (bool, interface{}) {
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
func (et *traverse) calculateCyclyTimeNextNode(nextNode interface{}, c *context, r *result) {
	timeResult := 0.0
	for nextNode != nil {
		nextNextNode := et.visit(nextNode, c, r)
		nextNextResult := r.CurrentCycleTime
		nextNode = nextNextNode
		timeResult += nextNextResult
	}
	r.CurrentCycleTime = timeResult
}

// xu li tinh toan cho loop
func (et *traverse) handleForLoop(start interface{}, end interface{}, c *context, r *result) interface{} {
	startGateway := start.(*gateway)
	endGateway := end.(*gateway)
	nextN := et.visit(startGateway.Next[0], c, r)
	timeResult := r.CurrentCycleTime
	et.calculateCyclyTimeNextNode(nextN, c, r)
	timeResult += r.CurrentCycleTime
	var reloop float64
	var nextNode interface{}

	// len(endGateway.Next) = 2
	for i, n := range endGateway.Next {
		if isGateway(n) && start == n {
			reloop = endGateway.branchingProbabilities[i]
		} else {
			nextNode = endGateway.Next[i]
		}
	}
	r.CurrentCycleTime = timeResult / (1 - reloop)
	if c.inLoop == 0 {
		r.TotalCycleTimeAllLoops += r.CurrentCycleTime
		block := blockQuality{Text: "Loop", Start: startGateway.Id, End: endGateway.Id, CycleTime: timeResult, ReworkProbability: 1 - reloop}
		r.LogsQuality = append(r.LogsQuality, block)
	}
	return nextNode
}

// level must be greater than 0
// example:
// [
//
//	    {
//	        "Text": "Step 1: calculate task 1",
//	        "Blocks": null
//	    },
//	    {
//	        "Text": "Step 2: calculate block ExclusiveGateway",
//	        "Blocks": [
//	            {
//	                "Text": "Calculate task 2",
//	                "Blocks": null
//	            },
//	            {
//	                "Text": "Calculate block ParallelGateway",
//	                "Blocks": null
//	            },
//	        ]
//	    },
//	]
//
// add(logsCycletime, {"Text": "Calculate task 3", "Blocks": null}, 1) to parallel block, result:
// [
//
//	    {
//	        "Text": "Step 1: calculate task 1",
//	        "Blocks": null
//	    },
//	    {
//	        "Text": "Step 2: calculate block ExclusiveGateway",
//	        "Blocks": [
//	            {
//	                "Text": "Calculate task 2",
//	                "Blocks": null
//	            },
//	            {
//	                "Text": "Calculate block ParallelGateway",
//	                "Blocks": [
//						{
//							"Text": "Calculate task 3",
//							"Blocks": null,
//						},
//					]
//	            },
//	        ]
//	    },
//	]
func (et *traverse) addBlockByLevel(block []blockCycleTime, newBlock blockCycleTime, level int) {
	var temp []blockCycleTime
	var tempParent []blockCycleTime = block
	temp = block
	for ; level > 0; level-- {
		tempParent = temp
		temp = temp[len(temp)-1].Blocks
	}
	temp = append(temp, newBlock)
	tempParent[len(tempParent)-1].Blocks = temp
}
