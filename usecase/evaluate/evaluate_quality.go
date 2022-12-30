package evaluate

type evaluateQuality struct {
}

func (eq *evaluateQuality) visit(i interface{}, c *context) (float64, interface{}) {
	switch i.(type) {
	case *event:
		return eq.visitForEvent(i.(*event), c)
	case *task:
		return eq.visitForTask(i.(*task), c)
	case *gateway:
		return eq.visitForGateway(i.(*gateway), c)
	}
	return 0, nil
}

func (eq *evaluateQuality) visitForEvent(e *event, c *context) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", e.Id, e.Name)
	return 0, nil
}

func (eq *evaluateQuality) visitForTask(tk *task, c *context) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", tk.Id, tk.Name)
	return 0, nil
}

func (eq *evaluateQuality) visitForGateway(g *gateway, c *context) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", g.Id, g.Name)
	return 0.0, nil
}
