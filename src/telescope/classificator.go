package main

type Classificator struct {
}

func InitClassificator() (*Classificator, error) {
	return new(Classificator), nil
}

func (Classificator) ClassifyData(data TelescopeData) bool {
	return true
}
