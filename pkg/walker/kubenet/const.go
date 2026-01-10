package kubenet

const (
	TEST_SUITE_IMAGE = "harbor.sail.se.sjtu.edu.cn/utils/test-suite:v0.2"
)

var TEST_SUITE_PORTS = []int32{5000, 5001, 5002}

const (
	FLOW_GEN_MODE_RANDOM           = "random"
	FLOW_GEN_MODE_WALKER_NEIGHBOUR = "walker-neighbour"
)
