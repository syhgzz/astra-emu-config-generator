package tle

func NewDefaultTLEParams() TLEParams {
	return TLEParams{
		Height:           500,
		ConstellationIdx: 1,
		Num:              300,
		Plane:            10,
		NeededPlane:      10,
		DeltaInit:        0,
		Phase:            0,
		Inclination:      50,
	}
}
