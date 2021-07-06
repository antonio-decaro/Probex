package main

type Classificator struct {
}

func InitClassificator() (*Classificator, error) {
	return new(Classificator), nil
}

func (Classificator) ClassifyData(data TelescopeData) bool { // THIS IS A MOCK CLASSIFICATOR
	const MIN_MASS = 0.0268
	const MAX_MASS = 2.0
	const MIN_RADIUS = 0.80
	const MAX_RADIUS = 3.9
	var minDistance, maxDistance float64

	switch data.StarType {
	case "OBA":
		minDistance = 2.56
		maxDistance = 3.65
	case "KGF":
		minDistance = 0.95
		maxDistance = 1.37
	default:
		return false
	}

	if minDistance > data.StarDistance || data.StarDistance > maxDistance ||
		data.Mass < MIN_MASS || data.Mass > MAX_MASS ||
		data.Radius < MIN_RADIUS || data.Radius > MAX_RADIUS {
		return false
	} else {
		return true
	}
}
