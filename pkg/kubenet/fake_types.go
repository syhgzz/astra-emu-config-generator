package kubenet

type MetaData struct {
	Name string `yaml:"name"`

	Namespace string `yaml:"namespace"`
}

type NetworkExperiment struct {
	APIVersion string `yaml:"apiVersion"`

	MetaData MetaData `yaml:"metadata"`

	Kind string `yaml:"kind"`

	Spec NetworkExperimentSpec `yaml:"spec,omitempty"`
}

// NetworkExperimentSpec defines the desired state of NetworkExperiment
type NetworkExperimentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	NetworkNodes      []NetworkNodeSpec `yaml:"network_nodes,omitempty"`
	FixedLinks        []FixedLinkSpec   `yaml:"fixed_links,omitempty"`
	NetworkParameters NetworkParameters `yaml:"network_parameters,omitempty"`
	LargeScaleModel   string            `yaml:"large_scale_model,omitempty"`
	SmallScaleModel   string            `yaml:"small_scale_model,omitempty"`
	Protocol          string            `yaml:"protocol"`
}

type NetworkParameters struct {
	Frequency   float64 `yaml:"frequency,omitempty"`
	BitRate     float64 `yaml:"bitrate,omitempty"`
	Modulation  string  `yaml:"modulation,omitempty"`
	BandWidth   float64 `yaml:"bandwidth,omitempty"`
	M           float64 `yaml:"m,omitempty"`
	PowerInDbm  float64 `yaml:"powerindbm,omitempty"`
	MaxDistance float64 `yaml:"maxdistance,omitempty"`
}

type Location struct {
	Longitude float64 `yaml:"longitude"`
	Latitude  float64 `yaml:"latitude"`
	Altitude  float64 `yaml:"altitude"`
	Time      int64   `yaml:"time"`
}

type NetworkNodeSpec struct {
	Name       string     `yaml:"name"`
	IntfNum    int        `yaml:"intf_num"`
	Model      int        `yaml:"model,omitempty"`
	ModelType  int        `yaml:"type,omitempty"`
	NodeType   string     `yaml:"node_type,omitempty"`
	DstIP      string     `yaml:"dst_ip,omitempty"`
	Trajectory []Location `yaml:"trajectory,omitempty"`
}

type FixedLinkSpec struct {
	SrcName string `yaml:"src_name"`
	DstName string `yaml:"dst_name"`
}
