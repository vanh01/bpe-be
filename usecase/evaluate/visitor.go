package evaluate

type visitorTime interface {
	visit(interface{}, *context) (float64, interface{})
	visitForEvent(*event, *context) (float64, interface{})
	visitForTask(*task, *context) (float64, interface{})
	visitForGateway(*gateway, *context) (float64, interface{})
}

type visitorFlexibility interface {
	visit(interface{}, *context) interface{}
	visitForEvent(*event, *context) interface{}
	visitForTask(*task, *context) interface{}
	visitForGateway(*gateway, *context) interface{}
}
