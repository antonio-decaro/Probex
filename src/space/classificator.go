package main

type Classificator struct {
}

func InitClassificator() (*Classificator, error) {
	return new(Classificator), nil
}

func (Classificator) ClassifyData(data SpaceProbeData) bool {
	return true
}
