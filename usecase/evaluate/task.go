package evaluate

type task struct {
	node
	CycleTime    float64
	StartGateway string
	EndGateway   string
}

func (t *task) accept(v visitor) (float64, interface{}) {
	return v.visitForTask(t)
}
