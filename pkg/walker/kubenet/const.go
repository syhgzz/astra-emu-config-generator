package kubenet

const (
	FLOW_GEN_MODE_RANDOM           = "random"
	FLOW_GEN_MODE_WALKER_NEIGHBOUR = "walker-neighbour"
)

const (
	X_INTF_NAME  = "islx"
	XX_INTF_NAME = "islxx"
	Y_INTF_NAME  = "isly"
	YY_INTF_NAME = "islyy"
	Z_INTF_NAME  = "islz"
	ZZ_INTF_NAME = "islzz"

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
