package evaluate

// tinh so luong gateway
func numberOfGatewayInNodes(nodes []interface{}) int {
	count := 0
	for _, n := range nodes {
		if isGateway(n) {
			count += 1
		}
	}
	return count
}

func isGateway(g interface{}) bool {
	switch g.(type) {
	case *gateway:
		return true
	}
	return false
}
