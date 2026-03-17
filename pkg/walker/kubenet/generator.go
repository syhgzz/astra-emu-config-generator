package kubenet

import "fmt"

type KubenetGenerator struct {
	Config       GeneratorConfig
	IntfParamSet *IntfParamSet
	TemplatePath string
	InputPath    string
	OutputPath   string
}

type Option func(*KubenetGenerator)

func WithTemplatePath(path string) Option {
	return func(g *KubenetGenerator) {
		g.TemplatePath = path
	}
}

func WithOutputPath(path string) Option {
	return func(g *KubenetGenerator) {
		g.OutputPath = path
	}
}

func WithInputPath(path string) Option {
	return func(g *KubenetGenerator) {
		g.InputPath = path
	}
}

func WithConfig(cfg GeneratorConfig) Option {
	return func(g *KubenetGenerator) {
		g.Config = cfg
	}
}

func NewKubenetGenerator(opts ...Option) *KubenetGenerator {
	g := &KubenetGenerator{}
	for _, opt := range opts {
		opt(g)
	}
	g.IntfParamSet = NewIntfParamSet(
		g.Config.GroundAccess, AccessParams{
			AccessIntfMaxLinkNum: g.Config.AccessIntfMaxLinkNum,
		},
	)
	return g
}

func (g *KubenetGenerator) Validate() error {
	cfg := g.Config

	if cfg.SatelliteNum <= 0 || cfg.PlaneNum <= 0 {
		return fmt.Errorf("satellite num and plane num must be positive")
	}
	if cfg.NeededPlane <= 0 || cfg.NeededPlane > cfg.PlaneNum {
		return fmt.Errorf("needed plane must be between 1 and plane")
	}
	if cfg.SatelliteNum%cfg.PlaneNum != 0 {
		return fmt.Errorf("satellite num must be divisible by plane")
	}
	if cfg.Phase < 0 || cfg.Phase >= float64(cfg.PlaneNum) {
		return fmt.Errorf("phase must be between 0 and plane")
	}
	return nil
}

func (g *KubenetGenerator) Generate() error {
	if err := g.Validate(); err != nil {
		return err
	}

	return g.GenerateKubenetConfig()
}
