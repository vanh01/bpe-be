package evaluate

type gateway struct {
	node
	branchingProbabilities []float64
}

func (g *gateway) isSplitGateway() bool {
	return len(g.Next) > 1
}

func (g *gateway) isJoinGateway() bool {
	return len(g.Previous) > 1
}
