package evaluate

type evaluateTransparency struct {
}

func (et *evaluateTransparency) visit(i interface{}, c *context) (float64, interface{}) {
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

func (et *evaluateTransparency) visitForEvent(e *event, c *context) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", e.Id, e.Name)
	return 0, nil
}

func (et *evaluateTransparency) visitForTask(tk *task, c *context) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", tk.Id, tk.Name)
	return 0, nil
}

func (et *evaluateTransparency) visitForGateway(g *gateway, c *context) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", g.Id, g.Name)
	return 0.0, nil
}
