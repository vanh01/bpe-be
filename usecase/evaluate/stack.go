package evaluate

type gateWayStack []*gateway

// IsEmpty : check if stack is empty
func (g *gateWayStack) IsEmpty() bool {
	return len(*g) == 0
}

// IsNotEmpty : check if stack is empty
func (g *gateWayStack) IsNotEmpty() bool {
	return len(*g) != 0
}

// Push a new value onto the stack
func (g *gateWayStack) Push(gw *gateway) {
	*g = append(*g, gw) // Simply append the new value to the end of the stack
}

// Pop Remove and return top element of stack. Return false if stack is empty.
func (g *gateWayStack) Pop() (*gateway, bool) {
	if g.IsEmpty() {
		return nil, false
	} else {
		index := len(*g) - 1   // Get the index of the top most element.
		element := (*g)[index] // Index into the slice and obtain the element.
		*g = (*g)[:index]      // Remove it from the stack by slicing it off.
		return element, true
	}
}

func (g *gateWayStack) Top() *gateway {
	return (*g)[len(*g)-1]
}

func (g *gateWayStack) Size() int {
	return len(*g)
}

func (g *gateWayStack) Get(index int) *gateway {
	return (*g)[index]
}
