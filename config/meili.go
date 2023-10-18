package config

type MeiliIndexConfig struct {
	Name                 string   `yaml:"name"`
	PrimaryKey           string   `yaml:"primary_key"`
	UpdateConfig         bool     `yaml:"update_config"`
	SearchableAttributes []string `yaml:"searchable_attributes"`
	FilterableAttributes []string `yaml:"filterable_attributes"`
	DisplayedAttributes  []string `yaml:"displayed_attributes"`
	SortableAttributes   []string `yaml:"sortable_attributes"`
}

type MeiliConfig struct {
	Host    string                      `yaml:"host"`
	Token   string                      `yaml:"token"`
	Indices map[string]MeiliIndexConfig `yaml:"indices"`
}
