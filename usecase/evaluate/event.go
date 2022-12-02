package evaluate

type event struct {
	node
}

func (e *event) accept(v visitor) (float64, interface{}) {
	return v.visitForEvent(e)
}
