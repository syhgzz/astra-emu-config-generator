package merge

import (
	"sync"

	"dslab.sjtu/kubenet/api/pkg/kubenet_config"
)

type Option func(*ConfigMerger)

// ConfigMerger handles merging of kubenet configuration files.
// It supports options for input/output paths, recursion, file filtering, and error handling.
type ConfigMerger struct {
	InputPaths  []string
	OutputPath  string
	Recursive   bool
	FilterExt   string
	IgnoreError bool

	networkExperiments []*kubenet_config.NetworkExperiment

	mu sync.Mutex
}

func WithInputPaths(paths []string) Option {
	return func(m *ConfigMerger) {
		m.InputPaths = paths
	}
}

func WithOutputPath(path string) Option {
	return func(m *ConfigMerger) {
		m.OutputPath = path
	}
}

func WithRecursive(recursive bool) Option {
	return func(m *ConfigMerger) {
		m.Recursive = recursive
	}
}

func WithFilterExt(ext string) Option {
	return func(m *ConfigMerger) {
		m.FilterExt = ext
	}
}

func WithIgnoreError(ignore bool) Option {
	return func(m *ConfigMerger) {
		m.IgnoreError = ignore
	}
}

// --- 构造函数 ---
func NewConfigMerger(opts ...Option) *ConfigMerger {
	m := &ConfigMerger{
		Recursive:   false,
		FilterExt:   ".yaml",
		IgnoreError: false,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}
