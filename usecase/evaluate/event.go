package evaluate

type event struct {
	node
}

func (e *event) accept(v visitor, c *context) (float64, interface{}) {
	return v.visitForEvent(e, c)
}
