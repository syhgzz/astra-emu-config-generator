package simple

import "github.com/sirupsen/logrus"

type SimpleGenerator struct {
	config     *SimpleConfig
	outputPath string
}

func NewSimpleGenerator(config *SimpleConfig, outputPath string) *SimpleGenerator {
	return &SimpleGenerator{config: config, outputPath: outputPath}
}

func (g *SimpleGenerator) Generate() error {
	switch g.config.Topology {
	case "line":
		return g.GenerateLine()
	case "pair":
		return g.GeneratePair()
	default:
		logrus.Errorf("Unsupported topology type: %s", g.config.Topology)
		return nil
	}
}
