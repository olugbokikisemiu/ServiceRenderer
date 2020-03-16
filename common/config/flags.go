package config

type Flags struct {
	ServiceFlags `konfig:"service_config" yaml:"service_config" bson:"service_config"`

	Test       bool    `konfig:"test" long:"test" yaml:"test" bson:",omitempty" json:",omitempty" description:"enable test env"`
	ConfigPath string  `konfig:"config" long:"config" default:"config.yaml" yaml:"config" bson:",omitempty" json:",omitempty" description:"path to config.yaml"`
	Addr       *string `long:"addr" yaml:"addr" bson:",omitempty" json:",omitempty" description:"address to listen for server requests"`
}

type ServiceFlags struct {
}
