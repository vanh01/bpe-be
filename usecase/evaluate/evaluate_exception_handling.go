package evaluate

type evaluateExceptionHandling struct {
}

func (eeh *evaluateExceptionHandling) visit(i interface{}, c *context) (float64, interface{}) {
	switch i.(type) {
	case *event:
		return eeh.visitForEvent(i.(*event), c)
	case *task:
		return eeh.visitForTask(i.(*task), c)
	case *gateway:
		return eeh.visitForGateway(i.(*gateway), c)
	}
	return 0, nil
}

func (eeh *evaluateExceptionHandling) visitForEvent(e *event, c *context) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", e.Id, e.Name)
	return 0, nil
}

func (eeh *evaluateExceptionHandling) visitForTask(tk *task, c *context) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", tk.Id, tk.Name)
	return 0, nil
}

func (eeh *evaluateExceptionHandling) visitForGateway(g *gateway, c *context) (float64, interface{}) {
	// fmt.Printf("Visit: %-50s| %-20s\n", g.Id, g.Name)
	return 0.0, nil
}
