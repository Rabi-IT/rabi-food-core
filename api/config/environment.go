package config

type environment string

func (e environment) IsDevelopment() bool {
	return e == "development"
}

func (e environment) IsProduction() bool {
	return e == "production"
}

func (e environment) IsTest() bool {
	return e == "test"
}

func (e environment) String() string {
	return string(e)
}
