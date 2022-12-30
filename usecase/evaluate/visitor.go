package evaluate

type visitor interface {
	visit(interface{}, *context, *result) interface{}
	visitForEvent(*event, *context, *result) interface{}
	visitForTask(*task, *context, *result) interface{}
	visitForGateway(*gateway, *context, *result) interface{}
}
