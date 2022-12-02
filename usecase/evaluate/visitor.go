package evaluate

type visitor interface {
	visit(interface{}) (float64, interface{})
	visitForEvent(*event) (float64, interface{})
	visitForTask(*task) (float64, interface{})
	visitForGateway(*gateway) (float64, interface{})
}
