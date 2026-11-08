package model

type Config struct {
	Output struct {
		ErrorFile string `yaml:"error"`
		Result    string `yaml:"result"`
	}
}
