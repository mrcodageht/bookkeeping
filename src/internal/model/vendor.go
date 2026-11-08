package model

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
