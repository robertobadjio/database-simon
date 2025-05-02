package config

// Engine ...
type Engine struct {
	Typ              string `yaml:"type"`
	PartitionsNumber int    `yaml:"partitions_number"`
}
