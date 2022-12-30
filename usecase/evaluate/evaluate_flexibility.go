package evaluate

type evaluateFlexibility struct {
}

func (ef *evaluateFlexibility) visit(i interface{}, c *context) interface{} {
	switch i.(type) {
	case *event:
		return ef.visitForEvent(i.(*event), c)
	case *task:
		return ef.visitForTask(i.(*task), c)
	case *gateway:
		return ef.visitForGateway(i.(*gateway), c)
	}
	return nil
}

func (ef *evaluateFlexibility) visitForEvent(e *event, c *context) interface{} {
	// fmt.Printf("Visit: %-50s| %-20s\n", e.Id, e.Name)
	return nil
}

func (ef *evaluateFlexibility) visitForTask(tk *task, c *context) interface{} {
	// fmt.Printf("Visit: %-50s| %-20s\n", tk.Id, tk.Name)
	return nil
}

func (ef *evaluateFlexibility) visitForGateway(g *gateway, c *context) interface{} {
	// fmt.Printf("Visit: %-50s| %-20s\n", g.Id, g.Name)
	return nil
}
