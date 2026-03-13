package kubenet

const (
	FLOW_GEN_MODE_RANDOM           = "random"
	FLOW_GEN_MODE_WALKER_NEIGHBOUR = "walker-neighbour"
)

const (
	X_INTF_NAME  = "isl-X"
	XX_INTF_NAME = "isl-XX"
	Y_INTF_NAME  = "isl-Y"
	YY_INTF_NAME = "isl-YY"
	Z_INTF_NAME  = "isl-Z"
	ZZ_INTF_NAME = "isl-ZZ"

	ACCESS_INTF_NAME = "acce"

	INTRA_ORBIT_ISL_TOPOLOGY_MODEL  = "intra-orbit-isl"
	INTER_ORBIT_ISL_TOPOLOGY_MODEL  = "inter-orbit-isl"
	INTER_SPHERE_ISL_TOPOLOGY_MODEL = "inter-sphere-isl"
	ACCE_TOPOLOGY_MODEL             = "cone1"
	EXTERNAL_TOPOLOGY_MODEL         = "external"

	ACCE_FADING_MODEL = "fading_model1"

	SAT_LINK_MODEL  = "satlaser"
	ACCE_LINK_MODEL = "Satellite"
)
