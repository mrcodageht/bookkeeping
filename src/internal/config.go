package internal

import (
	"os"

	"go.yaml.in/yaml/v4"
	"mrccommunity.com/bookkeping/src/internal/model"
)

func LoadVendor(vendorFile string) (*model.VendorsYAML, error) {
	vendors := &model.VendorsYAML{}
	err := load(vendorFile, vendors)
	return vendors, err
}

func LoadConfig(configFile string) (*model.Config, error) {
	configs := &model.Config{}
	err := load(configFile, configs)
	return configs, err
}

func load(filepath string, model any) error {
	var fileByte []byte
	var err error

	fileByte, err = os.ReadFile(filepath)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(fileByte, model)
	return err
}
