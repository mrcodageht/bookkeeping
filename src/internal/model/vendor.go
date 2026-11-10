package model

import "fmt"

type Vendor struct {
	Name     string   `yaml:"name"`
	Category string   `yaml:"category"`
	Patterns []string `yaml:"patterns"`
}

type VendorsYAML struct {
	Vendors []Vendor `yaml:"vendors"`
}

func UnknownVendor() *Vendor {
	return &Vendor{
		Name:     "Unknown",
		Category: "Unknown",
		Patterns: []string{},
	}
}

func (v *Vendor) ToString() string {

	return fmt.Sprintf("%15s \t|\t %20s \t| %s \t", v.Name, v.Category, v.Patterns)
}
