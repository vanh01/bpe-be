package evaluate

type graphBuilder struct {
	inputElement element
	graph        *map[string]interface{}
	mapElement   map[string]element
}

func (g *graphBuilder) buildEventNode(startNode *event) {
	switch g.inputElement.Name {
	case "StartEvent":
		g.buildStartEvent(startNode)
	case "IntermediateCatchEvent":
		g.buildIntermediateCatchEvent()
	case "IntermediateThrowEvent":
		g.buildIntermediateThrowEvent()
	}

}

func (g *graphBuilder) buildEvent(ele *element) *event {
	var tmpNode *event
	tmpNode = (*g.graph)[ele.Id].(*event)
	tmpNode.Previous = getNodeOfIds(ele.Incoming, g.graph)
	tmpNode.Next = getNodeOfIds(ele.Outgoing, g.graph)
	(*g.graph)[ele.Id] = tmpNode
	return tmpNode
}

func (g *graphBuilder) buildStartEvent(startNode *event) {
	*startNode = *g.buildEvent(&g.inputElement)
}

func (g *graphBuilder) buildIntermediateCatchEvent() {
	afterTask := getNodeOfIds(g.inputElement.Outgoing, g.graph)[0].(*task)
	eventSourceNodes := getNodeOfIds(g.inputElement.Source, g.graph)
	for i, _ := range eventSourceNodes {
		beforeTask := getNodeOfIds((g.mapElement)[g.inputElement.Source[i]].Incoming, g.graph)[0].(*task)
		(*g.graph)[beforeTask.Id].(*task).Next = append(make([]interface{}, 0), afterTask)
	}
}

func (g *graphBuilder) buildIntermediateThrowEvent() {
	beforeTask := getNodeOfIds(g.inputElement.Incoming, g.graph)[0].(*task)
	afterTask := getNodeOfIds((g.mapElement)[g.inputElement.Target].Outgoing, g.graph)[0].(*task)
	(*g.graph)[afterTask.Id].(*task).Previous = append(make([]interface{}, 0), beforeTask)
}
