package evaluate

type event struct {
	node
}

type linkEvent struct {
	event
	source []interface{}
	target interface{}
}
