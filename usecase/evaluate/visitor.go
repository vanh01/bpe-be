package evaluate

type visitor interface {
	visit(interface{}, *context) (float64, interface{})
	visitForEvent(*event, *context) (float64, interface{})
	visitForTask(*task, *context) (float64, interface{})
	visitForGateway(*gateway, *context) (float64, interface{})
}
