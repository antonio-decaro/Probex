package main

type Classificator struct {
}

func InitClassificator() *Classificator {
	return new(Classificator)
}

func (Classificator) ClassifyData(data SpaceProbeData) bool {
	return true
}
