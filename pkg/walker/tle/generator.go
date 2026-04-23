package tle

import (
	"fmt"
)

type TLEGenerator struct {
	Params     TLEParams
	OutputPath string
}

type Option func(*TLEGenerator)

func WithOutput(path string) Option {
	return func(g *TLEGenerator) {
		g.OutputPath = path
	}
}

func WithParams(params TLEParams) Option {
	return func(g *TLEGenerator) {
		g.Params = params
	}
}

func NewTLEGenerator(opts ...Option) *TLEGenerator {
	g := &TLEGenerator{
		OutputPath: "test.tle",
	}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

func (g *TLEGenerator) Validate() error {
	p := g.Params

	if p.Height <= 0 || p.Num <= 0 || p.Plane <= 0 {
		return fmt.Errorf("height, num, plane must be positive")
	}
	if p.NeededPlane <= 0 || p.NeededPlane > p.Plane {
		return fmt.Errorf("needed plane must be between 1 and plane")
	}
	if p.Num%p.Plane != 0 {
		return fmt.Errorf("satellite num must be divisible by plane")
	}
	if p.Inclination < 0 || p.Inclination > 180 {
		return fmt.Errorf("inclination must be between 0 and 180")
	}
	if p.Phase < 0 || p.Phase >= float64(p.Plane) {
		return fmt.Errorf("phase must be between 0 and plane")
	}
	return nil
}

func (g *TLEGenerator) Generate() error {
	if err := g.Validate(); err != nil {
		return err
	}

	items := g.Params.NewTLE()
	err := items.WriteFile(g.OutputPath)
	if err != nil {
		return fmt.Errorf("Failed to write TLE file: %v\n", err)
	}
	return nil
}
