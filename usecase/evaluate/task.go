package evaluate

type task struct {
	node
	CycleTime float64
}

func (t *task) accept(v visitor, c *context) (float64, interface{}) {
	return v.visitForTask(t, c)
}
