package loader

// ConfigLoader loads  and unmarshal config files
type ConfigLoader interface {
	Load(config interface{}, file string) error      // fills the config with unmarshalled data from the file
	PlainLoad(config interface{}, file string) error // fills the config with unmarshalled data from the file, without extension checking
}

// ConfigDumper marshals config to a file
type ConfigDumper interface {
	Dump(config interface{}, file string) error // serializes config to file
}
